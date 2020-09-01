package job

import (
	"fmt"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/version"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

type bpfCreator struct{}

func (b *bpfCreator) create(targetPod *v1.Pod, targetDetails *data.TargetDetails) (string, *batchv1.Job) {
	id := string(uuid.NewUUID())
	var imageName string
	if targetDetails.Image != "" {
		imageName = targetDetails.Image
	} else {
		imageName = fmt.Sprintf("%s:%s-bpf", baseImageName, version.GetCurrent())
	}

	commonMeta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("kubectl-flame-%s", id),
		Namespace: targetDetails.Namespace,
		Labels: map[string]string{
			"kubectl-flame/id": id,
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
			Template: v1.PodTemplateSpec{
				ObjectMeta: commonMeta,
				Spec: v1.PodSpec{
					HostPID: true,
					Volumes: []apiv1.Volume{
						{
							Name: "sys",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/sys",
								},
							},
						},
						{
							Name: "modules",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/lib/modules",
								},
							},
						},
					},
					InitContainers: nil,
					Containers: []apiv1.Container{
						{
							ImagePullPolicy: v1.PullAlways,
							Name:            "kubectl-flame",
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
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "sys",
									MountPath: "/sys",
								},
								{
									Name:      "modules",
									MountPath: "/lib/modules",
								},
							},
							SecurityContext: &v1.SecurityContext{
								Privileged: boolPtr(true),
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
