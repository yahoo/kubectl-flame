//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package kubernetes

import (
	"context"
	"fmt"
	"os"

	"github.com/VerizonMedia/kubectl-flame/cli/cmd/kubernetes/job"

	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func LaunchFlameJob(targetPod *v1.Pod, cfg *data.FlameConfig, ctx context.Context) (string, *batchv1.Job, error) {
	id, flameJob, err := job.Create(targetPod, cfg)
	if err != nil {
		return "", nil, fmt.Errorf("unable to create job: %w", err)
	}

	if cfg.TargetConfig.DryRun {
		err := printJob(flameJob)
		return "", nil, err
	}

	createJob, err := clientSet.
		BatchV1().
		Jobs(cfg.JobConfig.Namespace).
		Create(ctx, flameJob, metav1.CreateOptions{})
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

func DeleteProfilingJob(job *batchv1.Job, targetDetails *data.TargetDetails, ctx context.Context) error {
	deleteStrategy := metav1.DeletePropagationForeground
	return clientSet.
		BatchV1().
		Jobs(job.Namespace).
		Delete(ctx, job.Name, metav1.DeleteOptions{
			PropagationPolicy: &deleteStrategy,
		})
}
