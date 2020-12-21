package job

import (
	"fmt"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/version"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

type pythonCreator struct{}

func (p *pythonCreator) create(targetPod *apiv1.Pod, targetDetails *data.TargetDetails) (string, *batchv1.Job) {
	id := string(uuid.NewUUID())
	var imageName string
	if targetDetails.Image != "" {
		imageName = targetDetails.Image
	} else {
		imageName = fmt.Sprintf("%s:%s-python", baseImageName, version.GetCurrent())
	}

	commonMeta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("kubectl-flame-%s", id),
		Namespace: targetDetails.Namespace,
		Labels: map[string]string{
			"kubectl-flame/id": id,
		},
		Annotations: map[string]string{
			"sidecar.istio.io/inject": "false",
		},
	}

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: commonMeta,
		Spec: batchv1.JobSpec{
			Parallelism:             int32Ptr(1),
			Completions:             int32Ptr(1),
			TTLSecondsAfterFinished: int32Ptr(5),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: commonMeta,
				Spec: apiv1.PodSpec{
					HostPID:        true,
					InitContainers: nil,
					Containers: []apiv1.Container{
						{
							ImagePullPolicy: apiv1.PullAlways,
							Name:            ContainerName,
							Image:           imageName,
							Command:         []string{"/app/agent"},
							Args: []string{id,
								string(targetPod.UID),
								targetDetails.ContainerName,
								targetDetails.ContainerId,
								targetDetails.Duration.String(),
								string(targetDetails.Language),
								targetDetails.Pgrep,
							},
							SecurityContext: &apiv1.SecurityContext{
								Privileged: boolPtr(true),
								Capabilities: &apiv1.Capabilities{
									Add: []apiv1.Capability{"SYS_PTRACE"},
								},
							},
						},
					},
					RestartPolicy: "Never",
					NodeName:      targetPod.Spec.NodeName,
				},
			},
		},
	}

	return id, job
}
