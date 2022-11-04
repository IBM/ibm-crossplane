module github.com/crossplane/crossplane

go 1.16

require (
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/Masterminds/semver v1.5.0
	github.com/alecthomas/kong v0.2.17
	github.com/aws/aws-sdk-go v1.31.6 // indirect
	github.com/crossplane/crossplane-runtime v0.15.0
	github.com/docker/cli v20.10.12+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20200926000217-2617742802f6+incompatible // indirect
	github.com/google/go-cmp v0.5.5
	github.com/google/go-containerregistry v0.4.1
	github.com/google/go-containerregistry/pkg/authn/k8schain v0.0.0-20210330174036-3259211c1f24
	github.com/imdario/mergo v0.3.12
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.4.1
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
	golang.org/x/tools v0.1.8
	k8s.io/api v0.21.2
	k8s.io/apiextensions-apiserver v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/code-generator v0.21.2
	k8s.io/utils v0.0.0-20210527160623-6fdb442a123b
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/controller-tools v0.3.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/docker/cli => github.com/docker/cli v20.10.12+incompatible
	github.com/docker/distribution => github.com/docker/distribution v2.8.1+incompatible
	github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.0.2
)
