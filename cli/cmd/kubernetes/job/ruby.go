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

type rubyCreator struct{}

func (r *rubyCreator) create(targetPod *apiv1.Pod, cfg *data.FlameConfig) (string, *batchv1.Job, error) {
	id := string(uuid.NewUUID())
	var imageName string
	var imagePullSecret []apiv1.LocalObjectReference
	args := []string{
		id,
		string(targetPod.UID),
		cfg.TargetConfig.ContainerName,
		cfg.TargetConfig.ContainerId,
		cfg.TargetConfig.Duration.String(),
		string(cfg.TargetConfig.Language),
		cfg.TargetConfig.Pgrep,
	}

	if cfg.TargetConfig.Image != "" {
		imageName = cfg.TargetConfig.Image
	} else {
		imageName = fmt.Sprintf("%s:%s-ruby", baseImageName, version.GetCurrent())
	}

	if cfg.TargetConfig.ImagePullSecret != "" {
		imagePullSecret = []apiv1.LocalObjectReference{{Name: cfg.TargetConfig.ImagePullSecret}}
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
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: commonMeta,
				Spec: apiv1.PodSpec{
					HostPID: true,
					Volumes: []apiv1.Volume{
						{
							Name: "target-filesystem",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: cfg.TargetConfig.DockerPath,
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
									Name:      "target-filesystem",
									MountPath: "/var/lib/docker",
								},
							},
							SecurityContext: &apiv1.SecurityContext{
								Privileged: boolPtr(true),
								Capabilities: &apiv1.Capabilities{
									Add: []apiv1.Capability{"SYS_PTRACE"},
								},
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
