//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package profiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
)

const (
	mountIdLocation          = "/var/lib/docker/image/overlay2/layerdb/mounts/%s/mount-id"
	targetFileSystemLocation = "/var/lib/docker/overlay2/%s/merged"
)

func SetUp(job *details.ProfilingJob) error {
	targetFs, err := getTargetFileSystemLocation(job.ContainerID)
	if err != nil {
		return err
	}

	err = os.RemoveAll("/tmp")
	if err != nil {
		return err
	}

	err = os.Symlink(path.Join(targetFs, "tmp"), "/tmp")
	if err != nil {
		return err
	}

	return copyProfilerToTempDir()
}

func copyProfilerToTempDir() error {
	cmd := exec.Command("cp", "-r", "/app/async-profiler", "/tmp")
	return cmd.Run()
}

func getTargetFileSystemLocation(containerId string) (string, error) {
	fileName := fmt.Sprintf(mountIdLocation, containerId)
	mountId, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(targetFileSystemLocation, string(mountId)), nil
}
