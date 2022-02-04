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
	perfRecordOutputFileName        = "/tmp/perf.data"
	flameGraphPlLocation            = "/app/FlameGraph/flamegraph.pl"
	flameGraphStackCollapseLocation = "/app/FlameGraph/stackcollapse-perf.pl"
	perfScriptOutputFileName        = "/tmp/perf.out"
	perfFoldedOutputFileName        = "/tmp/perf.folded"
	flameGraphPerfOutputFile        = "/tmp/perf.svg"
)

type PerfProfiler struct{}

func (p *PerfProfiler) SetUp(job *details.ProfilingJob) error {
	return nil
}

func (p *PerfProfiler) Invoke(job *details.ProfilingJob) error {
	err := p.runPerfRecord(job)
	if err != nil {
		return fmt.Errorf("perf record failed: %s", err)
	}

	err = p.runPerfScript(job)
	if err != nil {
		return fmt.Errorf("perf script failed: %s", err)
	}

	err = p.foldPerfOutput(job)
	if err != nil {
		return fmt.Errorf("folding perf output failed: %s", err)
	}

	err = p.generateFlameGraph(job)
	if err != nil {
		return fmt.Errorf("flamegraph generation failed: %s", err)
	}

	return utils.PublishFlameGraph(flameGraphPerfOutputFile)
}

func (p *PerfProfiler) runPerfRecord(job *details.ProfilingJob) error {
	pid, err := utils.FindRootProcessId(job)
	if err != nil {
		return err
	}

	duration := strconv.Itoa(int(job.Duration.Seconds()))
	cmd := exec.Command(perfLocation, "record", "-p", pid, "-o", perfRecordOutputFileName, "-g", "--", "sleep", duration)

	return cmd.Run()
}

func (p *PerfProfiler) runPerfScript(job *details.ProfilingJob) error {
	f, err := os.Create(perfScriptOutputFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command(perfLocation, "script", "-i", perfRecordOutputFileName)
	cmd.Stdout = f

	return cmd.Run()
}

func (p *PerfProfiler) foldPerfOutput(job *details.ProfilingJob) error {
	f, err := os.Create(perfFoldedOutputFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command(flameGraphStackCollapseLocation, perfScriptOutputFileName)
	cmd.Stdout = f

	return cmd.Run()
}

func (p *PerfProfiler) generateFlameGraph(job *details.ProfilingJob) error {
	inputFile, err := os.Open(perfFoldedOutputFileName)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(flameGraphPerfOutputFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	cmd := exec.Command(flameGraphPlLocation)
	cmd.Stdin = inputFile
	cmd.Stdout = outputFile

	return cmd.Run()
}
