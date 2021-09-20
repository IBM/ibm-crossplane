/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package claim

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/claim"
)

// Error strings.
const (
	errUpdateClaim           = "cannot update composite resource claim"
	errUpdateComposite       = "cannot update composite resource"
	errBindClaimConflict     = "cannot bind claim that references a different composite resource"
	errBindCompositeConflict = "cannot bind composite resource that references a different claim"
	errGetSecret             = "cannot get composite resource's connection secret"
	errSecretConflict        = "cannot establish control of existing connection secret"
	errCreateOrUpdateSecret  = "cannot create or update connection secret"
	errCreateClient			 = "cannot create go client"
)

// An APIBinder binds claims to composites by updating them in a Kubernetes API
// server.
type APIBinder struct {
	client client.Client
	typer  runtime.ObjectTyper
}

// NewAPIBinder returns a new APIBinder.
func NewAPIBinder(c client.Client, t runtime.ObjectTyper) *APIBinder {
	return &APIBinder{client: c, typer: t}
}

// Bind the supplied claim to the supplied composite.
func (a *APIBinder) Bind(ctx context.Context, cm resource.CompositeClaim, cp resource.Composite) error {
	// IBM Patch: Move resourceRef to status
	existing := GetResourceReference(cm)
	proposed := meta.ReferenceTo(cp, resource.MustGetKind(cp, a.typer))
	if existing != nil && !cmp.Equal(existing, proposed, cmpopts.IgnoreFields(corev1.ObjectReference{}, "UID")) {
		return errors.New(errBindClaimConflict)
	}

	// Propagate the actual external name back from the composite to the claim.
	meta.SetExternalName(cm, meta.GetExternalName(cp))

	// We set the claim's resource reference first in order to reduce the chance
	// of leaking newly created composite resources. We want as few calls that
	// could fail and trigger a requeue between composite creation and reference
	// persistence as possible.
	// IBM Patch: Move resourceRef to status
	if err := SetResourceRef(ctx, a.client, cm, proposed); err != nil {
		return err
	}
	if err := a.client.Update(ctx, cm); err != nil {
		return errors.Wrap(err, errUpdateClaim)
	}

	existing = cp.GetClaimReference()
	proposed = meta.ReferenceTo(cm, resource.MustGetKind(cm, a.typer))
	if existing != nil && !cmp.Equal(existing, proposed, cmpopts.IgnoreFields(corev1.ObjectReference{}, "UID")) {
		return errors.New(errBindCompositeConflict)
	}

	cp.SetClaimReference(proposed)
	return errors.Wrap(a.client.Update(ctx, cp), errUpdateComposite)
}

// IBM Patch: Move resourceRef to status

// GetResourceReference Patch method which read data from status instead of spec
func GetResourceReference(cm resource.CompositeClaim) *corev1.ObjectReference {
	out := &corev1.ObjectReference{}
	data, ok := cm.(*claim.Unstructured)
	if data == nil || !ok {
		// back to standard inside one test where we can not change mock because of different repo
		return cm.GetResourceReference()
	}
	if err := fieldpath.Pave(data.Object).GetValueInto("status.resourceRef", out); err != nil {
		// back to standard inside one test where we can not change mock because of different repo
		if err := fieldpath.Pave(data.Object).GetValueInto("spec.resourceRef", out); err != nil {
			return nil
		}
		return out
	}
	if out.Name == "" {
		return nil
	}
	return out
}

// SetResourceRef - you can set resourceRef or you can remove it , if you set nil
func SetResourceRef(ctx context.Context, c client.Client, cm resource.CompositeClaim, resourceRef interface{}) error {
	data, ok := cm.(*claim.Unstructured)
	if !ok {
		return nil
	}

	// initailize status if we need to set it and it is not exists
	if data.Object["status"] == nil && resourceRef != nil {
		data.Object["status"] = map[string]interface{}{
			"resourceRef": resourceRef,
		}
	}

	iStatus, err := fieldpath.Pave(data.Object).GetValue("status")
	if err != nil {
		return errors.New(errUnsupportedClaimSpec)
	}

	status, ok := iStatus.(map[string]interface{})
	if !ok {
		return errors.New(errUnsupportedClaimSpec)
	}

	if resourceRef != nil {
		status["resourceRef"] = resourceRef
	} else {
		delete(status, "resourceRef")
	}
	if err := c.Status().Update(ctx, cm); err != nil {
		return errors.Wrap(err, errUpdateClaim)
	}
	return nil
}

// IBM Patch End: Move resourceRef to status

// An APIConnectionPropagator propagates connection details by reading
// them from and writing them to a Kubernetes API server.
type APIConnectionPropagator struct {
	client resource.ClientApplicator
	typer  runtime.ObjectTyper
}

// NewAPIConnectionPropagator returns a new APIConnectionPropagator.
func NewAPIConnectionPropagator(c client.Client, t runtime.ObjectTyper) *APIConnectionPropagator {
	return &APIConnectionPropagator{
		client: resource.ClientApplicator{Client: c, Applicator: resource.NewAPIUpdatingApplicator(c)},
		typer:  t,
	}
}

// PropagateConnection details from the supplied resource.
func (a *APIConnectionPropagator) PropagateConnection(ctx context.Context, to resource.LocalConnectionSecretOwner, from resource.ConnectionSecretOwner) (bool, error) {
	// Either from does not expose a connection secret, or to does not want one.
	if from.GetWriteConnectionSecretToReference() == nil || to.GetWriteConnectionSecretToReference() == nil {
		return false, nil
	}

	n := types.NamespacedName{
		Namespace: from.GetWriteConnectionSecretToReference().Namespace,
		Name:      from.GetWriteConnectionSecretToReference().Name,
	}
	fs := &corev1.Secret{}

	// IBM Patch: Remove cluster permission for Secrets
	// - replace a.client.Get() with newly created client
	//   to avoid using informers underneath (they had cluster scope)
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return false, errors.Wrap(err, errCreateClient)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	fs, err = clientset.CoreV1().Secrets(n.Namespace).Get(context.TODO(), n.Name, metav1.GetOptions{})
	// IBM Patch end

	if err != nil {
		return false, errors.Wrap(err, errGetSecret)
	}

	// Make sure 'from' is the controller of the connection secret it references
	// before we propagate it. This ensures a resource cannot use Crossplane to
	// circumvent RBAC by propagating a secret it does not own.
	if c := metav1.GetControllerOf(fs); c == nil || c.UID != from.GetUID() {
		return false, errors.New(errSecretConflict)
	}

	ts := resource.LocalConnectionSecretFor(to, resource.MustGetKind(to, a.typer))
	ts.Data = fs.Data

	// IBM Patch: Remove cluster permission for Secrets
	// - replace a.client.Apply() with .Get() and .Create() for newly created client
	if _, err := clientset.CoreV1().Secrets(ts.Namespace).Get(context.TODO(), ts.Name, metav1.GetOptions{}); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return false, errors.Wrap(err, errGetSecret)
		}
		_, err := clientset.CoreV1().Secrets(ts.ObjectMeta.Namespace).Create(context.TODO(), ts, metav1.CreateOptions{})
		if err != nil {
			return false, errors.Wrap(err, errCreateOrUpdateSecret)
		}
		return true, nil
	}
	return false, nil

	// IBM Patch end
}
