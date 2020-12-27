package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestResourceConfig_ParseResources(t *testing.T) {
	tt := []struct {
		name       string
		resConf    ResourceConfig
		want       apiv1.ResourceList
		wantErrMsg string
	}{
		{
			name: "empty config yields no resource list",
		},
		{
			name: "invalid CPU yields error",
			resConf: ResourceConfig{
				CPU: "test",
			},
			wantErrMsg: "unable to parse CPU value",
		},
		{
			name: "invalid memory yields error",
			resConf: ResourceConfig{
				Memory: "test",
			},
			wantErrMsg: "unable to parse memory value",
		},
		{
			name: "invalid CPU shortcircuits",
			resConf: ResourceConfig{
				CPU:    "test",
				Memory: "200Mi",
			},
			wantErrMsg: "unable to parse CPU value",
		},
		{
			name: "only CPU is parsed correctly",
			resConf: ResourceConfig{
				CPU: "200m",
			},
			want: apiv1.ResourceList{
				apiv1.ResourceCPU: resource.MustParse("200m"),
			},
		},
		{
			name: "only memory is parsed correctly",
			resConf: ResourceConfig{
				Memory: "200Mi",
			},
			want: apiv1.ResourceList{
				apiv1.ResourceMemory: resource.MustParse("200Mi"),
			},
		},
		{
			name: "both CPU and memory are parsed correctly",
			resConf: ResourceConfig{
				CPU:    "200m",
				Memory: "200Mi",
			},
			want: apiv1.ResourceList{
				apiv1.ResourceCPU:    resource.MustParse("200m"),
				apiv1.ResourceMemory: resource.MustParse("200Mi"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.resConf.ParseResources()

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestJobDetails_ToResourceRequirements(t *testing.T) {
	tt := []struct {
		name       string
		jobDetails *JobDetails
		want       apiv1.ResourceRequirements
		wantErrMsg string
	}{
		{
			name:       "empty resources yields empty requirements",
			jobDetails: &JobDetails{},
			want:       apiv1.ResourceRequirements{},
		},
		{
			name: "invalid request CPU yields error",
			jobDetails: &JobDetails{
				RequestConfig: ResourceConfig{
					CPU: "test",
				},
			},
			wantErrMsg: "unable to generate container requests",
		},
		{
			name: "invalid request mem yields error",
			jobDetails: &JobDetails{
				RequestConfig: ResourceConfig{
					Memory: "test",
				},
			},
			wantErrMsg: "unable to generate container requests",
		},
		{
			name: "valid requests yields requests only",
			jobDetails: &JobDetails{
				RequestConfig: ResourceConfig{
					CPU:    "100m",
					Memory: "200Mi",
				},
			},
			want: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("100m"),
					apiv1.ResourceMemory: resource.MustParse("200Mi"),
				},
			},
		},
		{
			name: "valid requests & invalid cpu limits yields error",
			jobDetails: &JobDetails{
				RequestConfig: ResourceConfig{
					CPU:    "100m",
					Memory: "200Mi",
				},
			},
			want: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("100m"),
					apiv1.ResourceMemory: resource.MustParse("200Mi"),
				},
			},
		},
		{
			name: "valid requests & invalid memory limits yields error",
			jobDetails: &JobDetails{
				RequestConfig: ResourceConfig{
					CPU:    "100m",
					Memory: "200Mi",
				},
				LimitConfig: ResourceConfig{
					CPU: "test",
				},
			},
			wantErrMsg: "unable to generate container limits",
		},
		{
			name: "valid requests & invalid memory limits yields error",
			jobDetails: &JobDetails{
				RequestConfig: ResourceConfig{
					CPU:    "100m",
					Memory: "200Mi",
				},
				LimitConfig: ResourceConfig{
					Memory: "test",
				},
			},
			wantErrMsg: "unable to generate container limits",
		},
		{
			name: "valid requests & memory yields both correctly",
			jobDetails: &JobDetails{
				RequestConfig: ResourceConfig{
					CPU:    "100m",
					Memory: "200Mi",
				},
				LimitConfig: ResourceConfig{
					CPU:    "100m",
					Memory: "200Mi",
				},
			},
			want: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("100m"),
					apiv1.ResourceMemory: resource.MustParse("200Mi"),
				},
				Limits: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("100m"),
					apiv1.ResourceMemory: resource.MustParse("200Mi"),
				},
			},
		},
		{
			name: "missing cpu limits yields requirements without cpu limits",
			jobDetails: &JobDetails{
				RequestConfig: ResourceConfig{
					CPU:    "100m",
					Memory: "200Mi",
				},
				LimitConfig: ResourceConfig{
					Memory: "200Mi",
				},
			},
			want: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceCPU:    resource.MustParse("100m"),
					apiv1.ResourceMemory: resource.MustParse("200Mi"),
				},
				Limits: apiv1.ResourceList{
					apiv1.ResourceMemory: resource.MustParse("200Mi"),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.jobDetails.ToResourceRequirements()

			if tc.wantErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}
