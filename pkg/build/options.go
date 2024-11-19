/*
Copyright 2021 The Kubernetes Authors.

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

package build

import (
	"fmt"
	"regexp"

	"sigs.k8s.io/kubetest2/pkg/build"
)

// ignore package name stutter
type BuildAndStageStrategy string //nolint:revive

const (
	// bazelStrategy builds and (optionally) stages using bazel
	bazelStrategy BuildAndStageStrategy = "bazel"
	// MakeStrategy builds using make and (optionally) stages using krel
	MakeStrategy BuildAndStageStrategy = "make"
)

type Options struct {
	Strategy           string `flag:"~strategy" desc:"Determines the build strategy to use either make or bazel."`
	StageLocation      string `flag:"~stage" desc:"Upload binaries to storage location if set, rightnow it supports cos://us/bucket123/<PATH> format for the IBM COS"`
	RepoRoot           string `flag:"-"`
	ImageLocation      string `flag:"~image-location" desc:"Image registry where built images are stored."`
	StageExtraGCPFiles bool   `flag:"-"`
	VersionSuffix      string `flag:"-"`
	UpdateLatest       bool   `flag:"~update-latest" desc:"Whether should upload the build number to the GCS"`
	TargetBuildArch    string `flag:"~target-build-arch" desc:"Target architecture for the test artifacts for dockerized build"`
	COSCredType        string `flag:"~cos-cred-type" desc:"IBM COS credential type(supported options: shared, cos_hmac)"`
	Builder
	Stager
}

func (o *Options) Validate() error {
	return o.implementationFromStrategy()
}

func (o *Options) implementationFromStrategy() error {
	switch BuildAndStageStrategy(o.Strategy) {
	case bazelStrategy:
		bazel := &build.Bazel{
			RepoRoot:      o.RepoRoot,
			StageLocation: o.StageLocation,
			ImageLocation: o.ImageLocation,
		}
		o.Builder = bazel
		o.Stager = bazel
	case MakeStrategy:
		o.Builder = &MakeBuilder{
			RepoRoot:        o.RepoRoot,
			TargetBuildArch: o.TargetBuildArch,
		}
		// skip the staging if stage is empty
		if o.StageLocation == "" {
			break
		}
		re := regexp.MustCompile(`^([a-zA-Z]+):\/\/([a-zA-Z0-9-]+)\/([a-zA-Z0-9-]+)(\/.*)?$`)
		matches := re.FindStringSubmatch(o.StageLocation)
		if len(matches) < 1 {
			return fmt.Errorf("invalid stage URL")
		}
		if matches[1] == "cos" {
			stager, err := NewIBMCOSStager(o.StageLocation, o.RepoRoot, o.TargetBuildArch, o.COSCredType)
			if err != nil {
				return err
			}
			o.Stager = stager
		} else {
			return fmt.Errorf("unsupported stage: %s", o.StageLocation)
		}
	default:
		return fmt.Errorf("unknown build strategy: %v", o.Strategy)
	}
	return nil
}
