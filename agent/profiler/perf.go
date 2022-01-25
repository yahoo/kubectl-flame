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
	perfLocation                    = "/app/perf"
	flameGraphPlLocation            = "/app/FlameGraph/flamegraph.pl"
	flameGraphStackCollapseLocation = "/app/FlameGraph/stackcollapse-perf.pl"
	rawPerfOutputFile               = "/tmp/perf.out"
	flameGraphPerfOutputFile        = "/tmp/perf.svg"
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
	cmd := exec.Command(perfLocation, "record", "-F99", "-p", pid, "-g", "--", "sleep", duration)

	return cmd.Run()
}

func (p *PerfProfiler) generateRawOutput(job *details.ProfilingJob) error {
	f, err := os.Create(rawPerfOutputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	generateRawCmd := exec.Command(perfLocation, "script")
	foldCmd := exec.Command(flameGraphStackCollapseLocation)
	foldCmd.Stdout = f

	foldCmd.Stdin, err = generateRawCmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = foldCmd.Start()
	if err != nil {
		return err
	}

	err = generateRawCmd.Run()
	if err != nil {
		return err
	}

	return foldCmd.Wait()
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

	flameGraphCmd := exec.Command(flameGraphPlLocation)
	flameGraphCmd.Stdin = inputFile
	flameGraphCmd.Stdout = outputFile

	return flameGraphCmd.Run()
}
