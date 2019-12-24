package injector

import (
	corev1 "k8s.io/api/core/v1"
)

func CreateInitContainers() []corev1.Container {
	return []corev1.Container{
		{
			Name:    "cloud-secrets-init",
			Image:   "daisaru11/cloud-secrets:0.0.1",
			Command: []string{"sh", "-c", "mkdir -p /cloud-secrets/bin/ && cp /usr/local/bin/cloud-secrets /cloud-secrets/bin/"},
			Args:    []string{},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "cloud-secrets",
					MountPath: "/cloud-secrets/",
				},
			},
		},
	}
}
