/*
Copyright 2023 The Crossplane Authors.

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

package revision

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

	pkgmetav1 "github.com/crossplane/crossplane/apis/pkg/meta/v1"
	v1 "github.com/crossplane/crossplane/apis/pkg/v1"
	"github.com/crossplane/crossplane/apis/pkg/v1alpha1"
)

type deploymentModifier func(*appsv1.Deployment)

func withPodTemplateLabels(labels map[string]string) deploymentModifier {
	return func(d *appsv1.Deployment) {
		d.Spec.Template.Labels = labels
	}
}

const (
	namespace = "ns"
)

// IBM Patch: reduce cluster permission
// added watchNamespace to the function because Crossplane propagates its
// WATCH_NAMESPACE value to Provider's deployment
func deployment(provider *pkgmetav1.Provider, revision string, watchNamespace string, modifiers ...deploymentModifier) *appsv1.Deployment {
	// IBM Patch end: reduce cluster permission
	var (
		replicas = int32(1)
	)

	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      revision,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"pkg.crossplane.io/revision": revision},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      provider.GetName(),
					Namespace: namespace,
					// IBM Patch: Add label for NSS operator
					Labels: map[string]string{"pkg.crossplane.io/revision": revision, "intent": "projected"},
					// IBM Patch end: Add label for NSS operator
				},
				Spec: corev1.PodSpec{
					// IBM Patch: rbac for Provider
					ServiceAccountName: provider.GetName(),
					// IBM Patch end: rbac for Provider
					Containers: []corev1.Container{
						{
							Name:  provider.GetName(),
							Image: provider.Spec.Controller.Image,
							// IBM Patch: reduce cluster permission
							Env: []corev1.EnvVar{{Name: "WATCH_NAMESPACE", Value: watchNamespace}},
							// IBM Patch end: reduce cluster permission
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}

	for _, modifier := range modifiers {
		modifier(d)
	}

	return d
}

func TestBuildProviderDeployment(t *testing.T) {
	type fields struct {
		provider *pkgmetav1.Provider
		revision *v1.ProviderRevision
		cc       *v1alpha1.ControllerConfig
	}

	provider := &pkgmetav1.Provider{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pkg",
		},
		Spec: pkgmetav1.ProviderSpec{
			Controller: pkgmetav1.ControllerSpec{
				Image: "img:tag",
			},
		},
	}

	revisionWithoutCC := &v1.ProviderRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rev-123",
		},
		Spec: v1.PackageRevisionSpec{
			ControllerConfigReference: nil,
			Package:                   "package",
			Revision:                  3,
		},
	}

	revisionWithCC := &v1.ProviderRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rev-123",
		},
		Spec: v1.PackageRevisionSpec{
			ControllerConfigReference: &xpv1.Reference{Name: "cc"},
			Package:                   "package",
			Revision:                  3,
		},
	}

	cc := &v1alpha1.ControllerConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: revisionWithCC.Name,
		},
		Spec: v1alpha1.ControllerConfigSpec{
			Metadata: &v1alpha1.PodObjectMeta{
				Labels: map[string]string{
					"k": "v",
				},
			},
		},
	}

	cases := map[string]struct {
		fields fields
		want   *appsv1.Deployment
	}{
		"MissingCC": {
			fields: fields{
				provider: provider,
				revision: revisionWithoutCC,
				cc:       nil,
			},
			// IBM Patch: reduce cluster permission
			want: deployment(provider, revisionWithCC.GetName(), "ns"),
			// IBM Patch end: reduce cluster permission
		},
		"CC": {
			fields: fields{
				provider: provider,
				revision: revisionWithCC,
				cc:       cc,
			},
			// IBM Patch: reduce cluster permission
			want: deployment(provider, revisionWithCC.GetName(), "ns", withPodTemplateLabels(map[string]string{
				// IBM Patch end: reduce cluster permission
				"pkg.crossplane.io/revision": revisionWithCC.GetName(),
				"k":                          "v",
				// IBM Patch: Add label for NSS operator
				"intent": "projected",
				// IBM Patch end: Add label for NSS operator
			})),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, deployment := buildProviderDeployment(tc.fields.provider, tc.fields.revision, tc.fields.cc, namespace)

			if diff := cmp.Diff(tc.want, deployment, cmpopts.IgnoreTypes(&corev1.SecurityContext{}, &corev1.PodSecurityContext{}, []metav1.OwnerReference{})); diff != "" {
				t.Errorf("-want, +got:\n%s\n", diff)
			}
		})
	}

}
