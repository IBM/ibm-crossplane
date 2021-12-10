/*
Copyright 2020 The Crossplane Authors.

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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/crossplane-runtime/pkg/meta"
	pkgmetav1 "github.com/crossplane/crossplane/apis/pkg/meta/v1"
	v1 "github.com/crossplane/crossplane/apis/pkg/v1"
	"github.com/crossplane/crossplane/apis/pkg/v1alpha1"
)

var (
	replicas = int32(1)
	// IBM Patch: Migration to use Provider.
	// hardcoded runAsUser fails on Openshift clusters
	// runAsUser                = int64(2000)
	// runAsGroup               = int64(2000)
	// IBM Patch end: Migration to use Provider.
	allowPrivilegeEscalation = false
	privileged               = false
	runAsNonRoot             = true
)

func buildProviderDeployment(provider *pkgmetav1.Provider, revision v1.PackageRevision, cc *v1alpha1.ControllerConfig, namespace string) (*corev1.ServiceAccount, *appsv1.Deployment) { // nolint:interfacer,gocyclo
	s := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			// IBM Patch: rbac for Provider
			// do not use revision name because its name is generated dynamically,
			// instead use known, constant provider name
			Name: provider.GetName(),
			// IBM Patch end: rbac for Provider
			Namespace:       namespace,
			OwnerReferences: []metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(revision, v1.ProviderRevisionGroupVersionKind))},
		},
	}
	pullPolicy := corev1.PullIfNotPresent
	if revision.GetPackagePullPolicy() != nil {
		pullPolicy = *revision.GetPackagePullPolicy()
	}
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            revision.GetName(),
			Namespace:       namespace,
			OwnerReferences: []metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(revision, v1.ProviderRevisionGroupVersionKind))},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"pkg.crossplane.io/revision": revision.GetName()},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      provider.GetName(),
					Namespace: namespace,
				},
				Spec: corev1.PodSpec{
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: &runAsNonRoot,
					},
					ServiceAccountName: s.GetName(),
					ImagePullSecrets:   revision.GetPackagePullSecrets(),
					Containers: []corev1.Container{
						{
							Name:  provider.GetName(),
							Image: provider.Spec.Controller.Image,
							Args: []string{
								"--debug",
							},
							// IBM Patch: reduce cluster permission
							// this env variable is needed in provider
							// to read NamespaceScope resource and restrict cache
							Env: []corev1.EnvVar{
								{
									Name:  "WATCH_NAMESPACE",
									Value: namespace,
								},
							},
							// IBM Patch end: reduce cluster permission
							ImagePullPolicy: pullPolicy,
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: &allowPrivilegeEscalation,
								Privileged:               &privileged,
								RunAsNonRoot:             &runAsNonRoot,
							},
						},
					},
				},
			},
		},
	}
	templateLabels := make(map[string]string)
	if cc != nil {
		s.Labels = cc.Labels
		s.Annotations = cc.Annotations
		d.Labels = cc.Labels
		d.Annotations = cc.Annotations
		if cc.Spec.Metadata != nil {
			d.Spec.Template.Annotations = cc.Spec.Metadata.Annotations
		}

		if cc.Spec.Metadata != nil {
			for k, v := range cc.Spec.Metadata.Labels {
				templateLabels[k] = v
			}
		}

		if cc.Spec.Replicas != nil {
			d.Spec.Replicas = cc.Spec.Replicas
		}
		if cc.Spec.Image != nil {
			d.Spec.Template.Spec.Containers[0].Image = *cc.Spec.Image
		}
		if cc.Spec.ImagePullPolicy != nil {
			d.Spec.Template.Spec.Containers[0].ImagePullPolicy = *cc.Spec.ImagePullPolicy
		}
		if len(cc.Spec.Ports) > 0 {
			d.Spec.Template.Spec.Containers[0].Ports = cc.Spec.Ports
		}
		if cc.Spec.NodeSelector != nil {
			d.Spec.Template.Spec.NodeSelector = cc.Spec.NodeSelector
		}
		if cc.Spec.ServiceAccountName != nil {
			d.Spec.Template.Spec.ServiceAccountName = *cc.Spec.ServiceAccountName
		}
		if cc.Spec.NodeName != nil {
			d.Spec.Template.Spec.NodeName = *cc.Spec.NodeName
		}
		if cc.Spec.PodSecurityContext != nil {
			d.Spec.Template.Spec.SecurityContext = cc.Spec.PodSecurityContext
		}
		if cc.Spec.SecurityContext != nil {
			d.Spec.Template.Spec.Containers[0].SecurityContext = cc.Spec.SecurityContext
		}
		if len(cc.Spec.ImagePullSecrets) > 0 {
			d.Spec.Template.Spec.ImagePullSecrets = cc.Spec.ImagePullSecrets
		}
		if cc.Spec.Affinity != nil {
			d.Spec.Template.Spec.Affinity = cc.Spec.Affinity
		}
		if len(cc.Spec.Tolerations) > 0 {
			d.Spec.Template.Spec.Tolerations = cc.Spec.Tolerations
		}
		if cc.Spec.PriorityClassName != nil {
			d.Spec.Template.Spec.PriorityClassName = *cc.Spec.PriorityClassName
		}
		if cc.Spec.RuntimeClassName != nil {
			d.Spec.Template.Spec.RuntimeClassName = cc.Spec.RuntimeClassName
		}
		if cc.Spec.ResourceRequirements != nil {
			d.Spec.Template.Spec.Containers[0].Resources = *cc.Spec.ResourceRequirements
		}
		if len(cc.Spec.Args) > 0 {
			d.Spec.Template.Spec.Containers[0].Args = cc.Spec.Args
		}
		if len(cc.Spec.EnvFrom) > 0 {
			d.Spec.Template.Spec.Containers[0].EnvFrom = cc.Spec.EnvFrom
		}
		if len(cc.Spec.Env) > 0 {
			d.Spec.Template.Spec.Containers[0].Env = cc.Spec.Env
		}
	}
	for k, v := range d.Spec.Selector.MatchLabels { // ensure the template matches the selector
		templateLabels[k] = v
	}
	// IBM Patch: Add label for NSS operator
	templateLabels["intent"] = "projected"
	// IBM Patch end: Add label for NSS operator
	d.Spec.Template.Labels = templateLabels

	return s, d
}
