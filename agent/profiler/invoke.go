//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package profiler

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/api"
	"github.com/fntlnz/mountinfo"
)

const (
	profilerDir = "/tmp/async-profiler"
	fileName    = profilerDir + "/flamegraph.svg"
	profilerSh  = profilerDir + "/profiler.sh"
)

func Invoke(job *details.ProfilingJob) error {
	pid, err := findJavaProcessId(job)
	if err != nil {
		return err
	}

	duration := strconv.Itoa(int(job.Duration.Seconds()))
	cmd := exec.Command(profilerSh, "-d", duration, "-f", fileName, "-e", "wall", pid)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return publishFlameGraph()
}

func findJavaProcessId(job *details.ProfilingJob) (string, error) {
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

					if strings.Contains(exeName, "java") {
						return dname, nil
					}
				}
			}
		}
	}
	return "", errors.New("Could not find any process")
}

func publishFlameGraph() error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	fgData := api.FlameGraphData{EncodedFile: encoded}

	return api.PublishEvent(api.FlameGraph, fgData)
}
