package job

import (
	"github.com/VerizonMedia/kubectl-flame/api"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
)

const baseImageName = "verizondigital/kubectl-flame"

var (
	jvm = jvmCreator{}
	bpf = bpfCreator{}
)

type creator interface {
	create(targetPod *v1.Pod, targetDetails *data.TargetDetails) (string, *batchv1.Job)
}

func Create(targetPod *v1.Pod, targetDetails *data.TargetDetails) (string, *batchv1.Job) {
	switch targetDetails.Language {
	case api.Java:
		return jvm.create(targetPod, targetDetails)
	case api.Go:
		return bpf.create(targetPod, targetDetails)
	}

	// Should not happen
	panic("got language without job creator")
}
