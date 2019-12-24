package injector

import (
	corev1 "k8s.io/api/core/v1"
)

func MutateContainer(container *corev1.Container) error {
	args := container.Command
	args = append(args, container.Args...)

	container.Command = []string{"/cloud-secrets/bin/cloud-secrets", "exec"}
	container.Args = args

	container.VolumeMounts = append(container.VolumeMounts, []corev1.VolumeMount{
		{
			Name:      "cloud-secrets",
			MountPath: "/cloud-secrets/",
		},
	}...)

	return nil
}
