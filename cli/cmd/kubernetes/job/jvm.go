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

type jvmCreator struct{}

func (c *jvmCreator) create(targetPod *v1.Pod, targetDetails *data.TargetDetails) (string, *batchv1.Job) {
	id := string(uuid.NewUUID())
	imageName := c.getAgentImage(targetDetails)
	args := []string{id, string(targetPod.UID),
		targetDetails.ContainerName, targetDetails.ContainerId,
		targetDetails.Duration.String(), string(targetDetails.Language)}

	if targetDetails.Pgrep != "" {
		args = append(args, targetDetails.Pgrep)
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
							Name: "target-filesystem",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/var/lib/docker",
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
							Args:            args,
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "target-filesystem",
									MountPath: "/var/lib/docker",
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

func (c *jvmCreator) getAgentImage(targetDetails *data.TargetDetails) string {
	if targetDetails.Image != "" {
		return targetDetails.Image
	}

	tag := fmt.Sprintf("%s-jvm", version.GetCurrent())
	if targetDetails.Alpine {
		tag = fmt.Sprintf("%s-alpine", tag)
	}

	return fmt.Sprintf("%s:%s", baseImageName, tag)
}
