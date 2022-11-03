module github.com/crossplane/crossplane

go 1.16

require (
	github.com/Masterminds/semver v1.5.0
	github.com/alecthomas/kong v0.2.17
	github.com/crossplane/crossplane-runtime v0.15.0
	github.com/google/go-cmp v0.5.9
	github.com/google/go-containerregistry v0.12.0
	github.com/google/go-containerregistry/pkg/authn/k8schain v0.0.0-20221030203717-1711cefd7eec
	github.com/imdario/mergo v0.3.12
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.6.0
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
	golang.org/x/tools v0.1.12
	k8s.io/api v0.25.3
	k8s.io/apiextensions-apiserver v0.21.2
	k8s.io/apimachinery v0.25.3
	k8s.io/client-go v0.25.3
	k8s.io/code-generator v0.21.2
	k8s.io/utils v0.0.0-20221012122500-cfd413dd9e85
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/controller-tools v0.3.0
	sigs.k8s.io/yaml v1.3.0
)

replace github.com/docker/cli => github.com/docker/cli v20.10.12+incompatible
