package utils

import (
	"fmt"
	"io/ioutil"
)

const (
	mountIdLocation          = "/var/lib/docker/image/overlay2/layerdb/mounts/%s/mount-id"
	targetFileSystemLocation = "/var/lib/docker/overlay2/%s/merged"
)

func GetTargetFileSystemLocation(containerId string) (string, error) {
	fileName := fmt.Sprintf(mountIdLocation, containerId)
	mountId, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(targetFileSystemLocation, string(mountId)), nil
}
