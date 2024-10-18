package deployer

import (
	"fmt"
	"os"

	"k8s.io/klog/v2"
)

func (d *deployer) Build() error {
	klog.V(1).Info("PowerVS deployer starting Build()")

	if err := d.init(); err != nil {
		return fmt.Errorf("build failed to init: %s", err)
	}
	_, err := d.BuildOptions.Build()
	if err != nil {
		return err
	}
	return nil
}

func (d *deployer) setRepoPathIfNotSet() error {
	if d.RepoRoot != "" {
		return nil
	}

	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory for setting Kubernetes root path: %s", err)
	}
	klog.V(1).Infof("defaulting repo root to the current directory: %s", path)
	d.RepoRoot = path

	return nil
}

// verifyBuildFlags only checks flags that are needed for Build
func (d *deployer) verifyBuildFlags() error {
	if err := d.setRepoPathIfNotSet(); err != nil {
		return err
	}
	d.BuildOptions.CommonBuildOptions.RepoRoot = d.RepoRoot
	return d.BuildOptions.Validate()
}
