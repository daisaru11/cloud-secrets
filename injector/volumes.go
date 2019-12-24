package injector

import (
	corev1 "k8s.io/api/core/v1"
)

func CreateVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "cloud-secrets",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: corev1.StorageMediumMemory,
				},
			},
		},
	}
}
