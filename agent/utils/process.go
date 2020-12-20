package utils

import (
	"bufio"
	"errors"
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
		api.Java:   "java",
		api.Python: "python",
	}
)

func getProcessName(job *details.ProfilingJob) string {
	if job.TargetProcessName != "" {
		return job.TargetProcessName
	}

	if val, ok := defaultProcessNames[job.Language]; ok {
		return val
	}

	return ""
}

func FindProcessId(job *details.ProfilingJob) (string, error) {
	name := getProcessName(job)
	foundProc := ""
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

					if name != "" {
						// search by process name
						if strings.Contains(exeName, name) {
							return dname, nil
						}
					} else {
						if foundProc != "" {
							return "", errors.New("found more than one process on container," +
								" specify process name using --pgrep flag")
						} else {
							foundProc = dname
						}
					}
				}
			}
		}
	}

	if foundProc != "" {
		return foundProc, nil
	}

	return "", errors.New("could not find any process")
}

func FindRootProcessId(job *details.ProfilingJob) (string, error) {
	name := getProcessName(job)
	procsAndParents := make(map[string]string)
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

					ppid, err := getProcessPPID(dname)
					if err != nil {
						return "", err
					}

					if name != "" {
						// search by process name
						if strings.Contains(exeName, name) {
							procsAndParents[dname] = ppid
						}
					} else {
						procsAndParents[dname] = ppid
					}
				}
			}
		}
	}

	return findRootProcess(procsAndParents)
}

func getProcessPPID(pid string) (string, error) {
	ppidKey := "PPid"
	statusFile, err := os.Open(path.Join("/proc", pid, "status"))
	if err != nil {
		return "", err
	}

	defer statusFile.Close()
	scanner := bufio.NewScanner(statusFile)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, ppidKey) {
			return strings.Fields(text)[1], nil
		}
	}

	return "", errors.New("unable to get process ppid")
}

func findRootProcess(procsAndParents map[string]string) (string, error) {
	for process, ppid := range procsAndParents {
		if _, ok := procsAndParents[ppid]; !ok {
			// Found process with ppid that is not in the same programming language - this is the root
			return process, nil
		}
	}

	return "", errors.New("could not find root process")
}
