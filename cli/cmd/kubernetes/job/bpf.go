//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package job

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"

	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/version"
)

type bpfCreator struct{}

func (b *bpfCreator) create(targetPod *apiv1.Pod, cfg *data.FlameConfig) (string, *batchv1.Job, error) {
	id := string(uuid.NewUUID())
	var imageName string
	var imagePullSecret []apiv1.LocalObjectReference
	if cfg.TargetConfig.Image != "" {
		imageName = cfg.TargetConfig.Image
	} else {
		imageName = fmt.Sprintf("%s:%s-bpf", baseImageName, version.GetCurrent())
	}

	if cfg.TargetConfig.ImagePullSecret != "" {
		imagePullSecret = []apiv1.LocalObjectReference{{Name: cfg.TargetConfig.ImagePullSecret}}
	}

	args := []string{
		id, string(targetPod.UID),
		cfg.TargetConfig.ContainerName, cfg.TargetConfig.ContainerId,
		cfg.TargetConfig.Duration.String(), string(cfg.TargetConfig.Language),
		string(cfg.TargetConfig.Event),
	}

	if cfg.TargetConfig.Pgrep != "" {
		args = append(args, cfg.TargetConfig.Pgrep)
	}

	commonMeta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("kubectl-flame-%s", id),
		Namespace: cfg.TargetConfig.Namespace,
		Labels: map[string]string{
			"kubectl-flame/id": id,
		},
		Annotations: map[string]string{
			"sidecar.istio.io/inject": "false",
		},
	}

	resources, err := cfg.JobConfig.ToResourceRequirements()
	if err != nil {
		return "", nil, fmt.Errorf("unable to generate resource requirements: %w", err)
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
			BackoffLimit:            int32Ptr(2),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: commonMeta,
				Spec: apiv1.PodSpec{
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
					ImagePullSecrets: imagePullSecret,
					InitContainers:   nil,
					Containers: []apiv1.Container{
						{
							ImagePullPolicy: apiv1.PullAlways,
							Name:            ContainerName,
							Image:           imageName,
							Command:         []string{"/app/agent"},
							Args:            args,
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
							SecurityContext: &apiv1.SecurityContext{
								Privileged: boolPtr(true),
							},
							Resources: resources,
						},
					},
					RestartPolicy: "Never",
					NodeName:      targetPod.Spec.NodeName,
				},
			},
		},
	}

	if cfg.TargetConfig.ServiceAccountName != "" {
		job.Spec.Template.Spec.ServiceAccountName = cfg.TargetConfig.ServiceAccountName
	}

	return id, job, nil
}
