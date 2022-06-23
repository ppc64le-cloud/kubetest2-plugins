package terraform

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ppc64le-cloud/kubetest2-plugins/data"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/terraform/exec"
)

const (
	// StateFileName is the default name for Terraform state files.
	StateFileName string = "terraform.tfstate"
)

func Apply(dir string, platform string, autoApprove bool, extraArgs ...string) (path string, err error) {
	err = unpackAndInit(dir, platform)
	if err != nil {
		return "", err
	}

	defaultArgs := []string{
		"-input=false",
		fmt.Sprintf("-state=%s", filepath.Join(dir, StateFileName)),
		fmt.Sprintf("-state-out=%s", filepath.Join(dir, StateFileName)),
	}
	if autoApprove {
		defaultArgs = append(defaultArgs, "-auto-approve")
	}
	args := append(defaultArgs, extraArgs...)
	sf := filepath.Join(dir, StateFileName)

	if exitCode := exec.Apply(dir, args); exitCode != 0 {
		return sf, errors.New("failed to apply Terraform")
	}
	return sf, nil
}

func Destroy(dir string, platform string, autoApprove bool, extraArgs ...string) (err error) {
	err = unpackAndInit(dir, platform)
	if err != nil {
		return err
	}

	defaultArgs := []string{
		"-input=false",
		fmt.Sprintf("-state=%s", filepath.Join(dir, StateFileName)),
		fmt.Sprintf("-state-out=%s", filepath.Join(dir, StateFileName)),
	}
	if autoApprove {
		defaultArgs = append(defaultArgs, "-auto-approve")
	}
	args := append(defaultArgs, extraArgs...)

	if exitCode := exec.Destroy(dir, args); exitCode != 0 {
		return errors.New("failed to destroy using Terraform")
	}
	return nil
}

func Output(dir string, platform string, extraArgs ...string) (output string, err error) {
	err = unpackAndInit(dir, platform)
	if err != nil {
		return "", err
	}

	defaultArgs := []string{
		fmt.Sprintf("-state=%s", filepath.Join(dir, StateFileName)),
		fmt.Sprint("-no-color"),
	}
	args := append(defaultArgs, extraArgs...)

	op, exitCode := exec.Output(dir, args)
	if exitCode != 0 {
		return "", errors.New("failed to terraform output")
	}
	return op, nil
}

// unpack unpacks the platform-specific Terraform modules into the
// given directory.
func unpack(dir string, platform string) (err error) {
	err = data.Unpack(dir, platform)
	if err != nil {
		return err
	}

	err = data.Unpack(filepath.Join(dir, "config.tf"), "config.tf")
	if err != nil {
		return err
	}

	return nil
}

// unpackAndInit unpacks the platform-specific Terraform modules into
// the given directory and then runs 'terraform init'.
func unpackAndInit(dir string, platform string) (err error) {
	err = unpack(dir, platform)
	if err != nil {
		return errors.Wrap(err, "failed to unpack Terraform modules")
	}

	args := []string{
		"-upgrade",
	}
	if exitCode := exec.Init(dir, args); exitCode != 0 {
		return errors.New("failed to initialize Terraform")
	}
	return nil
}
