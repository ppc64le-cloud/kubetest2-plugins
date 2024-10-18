package options

import (
	"sigs.k8s.io/kubetest2/pkg/build"
)

type BuildOptions struct {
	CommonBuildOptions *build.Options
}

var _ build.Builder = &BuildOptions{}
var _ build.Stager = &BuildOptions{}

func (bo *BuildOptions) Validate() error {
	return bo.CommonBuildOptions.Validate()
}

func (bo *BuildOptions) Build() (string, error) {
	return bo.CommonBuildOptions.Build()
}

func (bo *BuildOptions) Stage(version string) error {
	return bo.CommonBuildOptions.Stage(version)
}
