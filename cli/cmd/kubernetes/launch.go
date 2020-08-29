//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package kubernetes

import (
	"context"
	"fmt"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/version"
	"os"

	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const imageName = "verizondigital/kubectl-flame"

func LaunchFlameJob(targetPod *v1.Pod, targetDetails *data.TargetDetails, ctx context.Context) (string, *batchv1.Job, error) {
	id := string(uuid.NewUUID())
	imageName := getAgentImage(targetDetails)

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
							Args: []string{id,
								string(targetPod.UID),
								targetDetails.ContainerName,
								targetDetails.ContainerId,
								targetDetails.Duration.String(),
							},
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

	if targetDetails.DryRun {
		err := printJob(job)
		return "", nil, err
	}

	createJob, err := clientSet.
		BatchV1().
		Jobs(targetDetails.Namespace).
		Create(ctx, job, metav1.CreateOptions{})

	if err != nil {
		return "", nil, err
	}

	return id, createJob, nil
}

func printJob(job *batchv1.Job) error {
	encoder := json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{
		Yaml: true,
	})

	return encoder.Encode(job, os.Stdout)
}

func getAgentImage(targetDetails *data.TargetDetails) string {
	if targetDetails.Image != "" {
		return targetDetails.Image
	}

	tag := fmt.Sprintf("%s-jvm", version.GetCurrent())
	if targetDetails.Alpine {
		tag = fmt.Sprintf("%s-alpine", tag)
	}

	return fmt.Sprintf("%s:%s", imageName, tag)
}

func DeleteProfilingJob(job *batchv1.Job, targetDetails *data.TargetDetails, ctx context.Context) error {
	deleteStrategy := metav1.DeletePropagationForeground
	return clientSet.
		BatchV1().
		Jobs(targetDetails.Namespace).
		Delete(ctx, job.Name, metav1.DeleteOptions{
			PropagationPolicy: &deleteStrategy,
		})
}
