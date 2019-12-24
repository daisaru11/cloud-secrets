package injector

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestShouldMutate(t *testing.T) {
	type testcase struct {
		pod  *corev1.Pod
		want bool
	}

	testcases := []testcase{
		{
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{},
				},
			},
			want: false,
		},
		{
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"cloud-secrets.daisaru11.dev/enabled": "true",
					},
				},
			},
			want: true,
		},
		{
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"cloud-secrets.daisaru11.dev/enabled":         "true",
						"cloud-secrets.daisaru11.dev/mutation-status": "mutated",
					},
				},
			},
			want: false,
		},
	}

	for _, tc := range testcases {
		injector := NewInjector(tc.pod)
		got, err := injector.ShouldMutate()

		if !assert.NoError(t, err) {
			continue
		}

		assert.Equal(t, tc.want, got)
	}
}

// nolint: funlen
func TestMutate(t *testing.T) {
	orig := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"cloud-secrets.daisaru11.dev/enabled": "true",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{

					Name:  "redis",
					Image: "redis:latest",
					Args:  []string{"redis-server"},
				},
			},
		},
	}
	injector := NewInjector(orig)

	mutated, err := injector.Mutate()
	if !assert.NoError(t, err) {
		return
	}

	expected := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"cloud-secrets.daisaru11.dev/enabled":         "true",
				"cloud-secrets.daisaru11.dev/mutation-status": "mutated",
			},
		},
		Spec: corev1.PodSpec{
			InitContainers: []corev1.Container{
				{
					Name:  "cloud-secrets-init",
					Image: "daisaru11/cloud-secrets:0.0.1",
					Command: []string{
						"sh",
						"-c",
						"mkdir -p /cloud-secrets/bin/ && cp /usr/local/bin/cloud-secrets /cloud-secrets/bin/",
					},
					Args: []string{},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "cloud-secrets",
							MountPath: "/cloud-secrets/",
						},
					},
				},
			},
			Containers: []corev1.Container{
				{

					Name:    "redis",
					Image:   "redis:latest",
					Command: []string{"/cloud-secrets/bin/cloud-secrets", "exec"},
					Args:    []string{"redis-server"},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "cloud-secrets",
							MountPath: "/cloud-secrets/",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "cloud-secrets",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium: corev1.StorageMediumMemory,
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, mutated)

	if diff := cmp.Diff(mutated, expected); diff != "" {
		t.Errorf("Mutate differs: (-got +want)\n%s", diff)
	}
}
