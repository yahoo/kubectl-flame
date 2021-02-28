package profiler

import (
	"fmt"
	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/api"
)

type FlameGraphProfiler interface {
	SetUp(job *details.ProfilingJob) error
	Invoke(job *details.ProfilingJob) error
}

var (
	jvm    = JvmProfiler{}
	bpf    = BpfProfiler{}
	python = PythonProfiler{}
	ruby   = RubyProfiler{}
)

func ForLanguage(lang api.ProgrammingLanguage) (FlameGraphProfiler, error) {
	switch lang {
	case api.Java:
		return &jvm, nil
	case api.Go:
		return &bpf, nil
	case api.Python:
		return &python, nil
	case api.Ruby:
		return &ruby, nil
	default:
		return nil, fmt.Errorf("could not find profiler for language %s", lang)
	}
}
