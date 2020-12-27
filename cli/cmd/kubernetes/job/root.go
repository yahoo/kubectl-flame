//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package job

import (
	"errors"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"

	"github.com/VerizonMedia/kubectl-flame/api"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
)

const (
	baseImageName = "verizondigital/kubectl-flame"
	ContainerName = "kubectl-flame"
)

var (
	jvm    = jvmCreator{}
	bpf    = bpfCreator{}
	python = pythonCreator{}
)

type creator interface {
	create(targetPod *apiv1.Pod, cfg *data.FlameConfig) (string, *batchv1.Job, error)
}

func Create(targetPod *apiv1.Pod, cfg *data.FlameConfig) (string, *batchv1.Job, error) {
	switch cfg.TargetConfig.Language {
	case api.Java:
		return jvm.create(targetPod, cfg)
	case api.Go:
		return bpf.create(targetPod, cfg)
	case api.Python:
		return python.create(targetPod, cfg)
	}

	// Should not happen
	return "", nil, errors.New("got language without job creator")
}
