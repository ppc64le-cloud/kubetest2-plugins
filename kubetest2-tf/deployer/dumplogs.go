package deployer

import (
	"fmt"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"path/filepath"
)

func (d *deployer) DumpClusterLogs() error {
	klog.Infof("Collecting cluster logs under %s", d.logsDir)
	// create a directory based on the generated path: _rundir/dump-cluster-logs
	if _, err := os.Stat(d.logsDir); os.IsNotExist(err) {
		if err := os.Mkdir(d.logsDir, os.ModePerm); err != nil {
			klog.Errorf("cannot create a directory in path %q. Err: %v", d.logsDir, err)
			return err
		}
	} else if err == nil {
		klog.Errorf("%q already exists. Please clean up directory", d.logsDir)
		return err
	} else {
		return fmt.Errorf("an error occured while obtaining directory stats. Err: %v", err)
	}
	outfile, err := os.Create(filepath.Join(d.logsDir, "cluster-info.log"))
	if err != nil {
		klog.Errorf("Failed to create a log file. Err: %v", err)
		return err
	}
	defer outfile.Close()
	command := []string{
		"kubectl",
		"cluster-info",
		"dump",
	}
	klog.Infof("About to run: %s", command)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = outfile
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("couldn't use kubectl to dump cluster info: %s", err)
	}
	klog.Infof("Executed %s successfully", command)
	// Todo: Include provider specific logic in this section. (Includes node level information/CRI/Services, etc.)
	return nil
}
