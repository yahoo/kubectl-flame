package profiler

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/agent/utils"
)

const (
	perfLocation             = "/app/perf"
	flameGraphVizLocation    = "stackvis"
	rawPerfOutputFile        = "/tmp/perf.out"
	flameGraphPerfOutputFile = "/tmp/perf.svg"
)

type PerfProfiler struct{}

func (p *PerfProfiler) SetUp(job *details.ProfilingJob) error {
	return nil
}

func (p *PerfProfiler) Invoke(job *details.ProfilingJob) error {
	err := p.runProfiler(job)
	if err != nil {
		return fmt.Errorf("profiling failed: %s", err)
	}

	err = p.generateRawOutput(job)
	if err != nil {
		return fmt.Errorf("raw output generation failed: %s", err)
	}

	err = p.generateFlameGraph(job)
	if err != nil {
		return fmt.Errorf("flamegraph generation failed: %s", err)
	}

	return utils.PublishFlameGraph(flameGraphPerfOutputFile)
}

func (p *PerfProfiler) runProfiler(job *details.ProfilingJob) error {
	pid, err := utils.FindRootProcessId(job)
	if err != nil {
		return err
	}

	duration := strconv.Itoa(int(job.Duration.Seconds()))
	cmd := exec.Command(perfLocation, "record", "-p", pid, "-g", "--", "sleep", duration)

	return cmd.Run()
}

func (p *PerfProfiler) generateRawOutput(job *details.ProfilingJob) error {
	f, err := os.Create(rawProfilerOutputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command(perfLocation, "script")
	cmd.Stdout = f

	return cmd.Run()
}

func (p *PerfProfiler) generateFlameGraph(job *details.ProfilingJob) error {
	inputFile, err := os.Open(rawPerfOutputFile)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(flameGraphPerfOutputFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	flameGraphCmd := exec.Command(flameGraphVizLocation, "perf")
	flameGraphCmd.Stdin = inputFile
	flameGraphCmd.Stdout = outputFile

	return flameGraphCmd.Run()
}
