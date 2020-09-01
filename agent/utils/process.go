package utils

import (
	"errors"
	"fmt"
	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/api"
	"github.com/fntlnz/mountinfo"
	"io"
	"os"
	"path"
	"strings"
)

var (
	defaultProcessNames = map[api.ProgrammingLanguage]string{
		api.Java: "java",
	}
)

func getProcessName(job *details.ProfilingJob) (string, error) {
	if job.TargetProcessName != "" {
		return job.TargetProcessName, nil
	}

	if val, ok := defaultProcessNames[job.Language]; ok {
		return val, nil
	}

	return "", fmt.Errorf("could not find default process name for language %s", job.Language)
}

func FindProcessId(job *details.ProfilingJob) (string, error) {
	name, err := getProcessName(job)
	if err != nil {
		return "", err
	}

	proc, err := os.Open("/proc")
	if err != nil {
		return "", err
	}

	defer proc.Close()

	for {
		dirs, err := proc.Readdir(15)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		for _, di := range dirs {
			if !di.IsDir() {
				continue
			}

			dname := di.Name()
			if dname[0] < '0' || dname[0] > '9' {
				continue
			}

			mi, err := mountinfo.GetMountInfo(path.Join("/proc", dname, "mountinfo"))
			if err != nil {
				continue
			}

			for _, m := range mi {
				root := m.Root
				if strings.Contains(root, job.PodUID) &&
					strings.Contains(root, job.ContainerName) {

					exeName, err := os.Readlink(path.Join("/proc", dname, "exe"))
					if err != nil {
						continue
					}

					if strings.Contains(exeName, name) {
						return dname, nil
					}
				}
			}
		}
	}
	return "", errors.New("could not find any process")
}
