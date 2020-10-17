//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package job

import (
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
	jvm = jvmCreator{}
	bpf = bpfCreator{}
)

type creator interface {
	create(targetPod *apiv1.Pod, targetDetails *data.TargetDetails) (string, *batchv1.Job)
}

func Create(targetPod *apiv1.Pod, targetDetails *data.TargetDetails) (string, *batchv1.Job) {
	switch targetDetails.Language {
	case api.Java:
		return jvm.create(targetPod, targetDetails)
	case api.Go:
		return bpf.create(targetPod, targetDetails)
	}

	// Should not happen
	panic("got language without job creator")
}
