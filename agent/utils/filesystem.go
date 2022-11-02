package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	overlayFsMountIdLocation = "/var/lib/docker/image/overlay2/layerdb/mounts/%s/mount-id"
	// The actual path on the host is /run/containerd/io... but we don't change the path within the agent to not break backwards compatibility with the flag
	containerdMountIdLocation = "/var/lib/docker/io.containerd.runtime.v2.task/k8s.io/%s/rootfs"
	targetFileSystemLocation  = "/var/lib/docker/overlay2/%s/merged"
)

func GetTargetFileSystemLocation(containerId string) (string, error) {
	fileName := fmt.Sprintf(overlayFsMountIdLocation, containerId)

	// Check if the node uses overlayfs mount location
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		// retry with containerd path
		// remove containerd:// prefix if necessary
		containerIdSanitized := strings.Replace(containerId, "containerd://", "", 1)
		fileName = fmt.Sprintf(containerdMountIdLocation, containerIdSanitized)
		_, err := os.Stat(fileName)
		if err == nil {
			// file exists, must be a containerd node
			// the path here is already the rootfs of the container, can return immediately
			return fileName, nil
		}
	}
	mountId, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(targetFileSystemLocation, string(mountId)), nil
}
