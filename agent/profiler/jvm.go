package profiler

import (
	"bytes"
	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/agent/utils"
	"os"
	"os/exec"
	"path"
	"strconv"
)

const (
	profilerDir = "/tmp/async-profiler"
	fileName    = profilerDir + "/flamegraph.svg"
	profilerSh  = profilerDir + "/profiler.sh"
)

type JvmProfiler struct{}

func (j *JvmProfiler) SetUp(job *details.ProfilingJob) error {
	targetFs, err := utils.GetTargetFileSystemLocation(job.ContainerID)
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

	return j.copyProfilerToTempDir()
}

func (j *JvmProfiler) Invoke(job *details.ProfilingJob) error {
	pid, err := utils.FindProcessId(job)
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

	return utils.PublishFlameGraph(fileName)
}

func (j *JvmProfiler) copyProfilerToTempDir() error {
	cmd := exec.Command("cp", "-r", "/app/async-profiler", "/tmp")
	return cmd.Run()
}
