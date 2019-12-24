package injector

import (
	"encoding/json"
	"strconv"

	"github.com/mattbaird/jsonpatch"
	corev1 "k8s.io/api/core/v1"
)

type Injector struct {
	pod *corev1.Pod
}

func NewInjector(pod *corev1.Pod) *Injector {
	return &Injector{
		pod: pod,
	}
}

func (i *Injector) GeneratePatch() ([]jsonpatch.JsonPatchOperation, error) {
	shouldMutate, err := i.ShouldMutate()

	if err != nil {
		return nil, err
	}

	if !shouldMutate {
		return nil, nil
	}

	muteted, err := i.Mutate()
	if err != nil {
		return nil, err
	}

	origJSON, err := json.Marshal(i.pod)
	if err != nil {
		return nil, err
	}

	mutatedJSON, err := json.Marshal(muteted)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.CreatePatch(origJSON, mutatedJSON)
	if err != nil {
		return nil, err
	}

	return patch, nil
}

func (i *Injector) ShouldMutate() (bool, error) {
	if i.pod.Annotations == nil {
		return false, nil
	}

	raw, ok := i.pod.Annotations[AnnotationEnabled]
	if !ok {
		return false, nil
	}

	enabled, err := strconv.ParseBool(raw)
	if err != nil {
		return false, err
	}

	if !enabled {
		return false, nil
	}

	raw, ok = i.pod.Annotations[AnnotationMutationStatus]
	if !ok {
		return true, nil
	}

	if raw == "mutated" {
		return false, nil
	}

	return true, nil
}

func (i *Injector) Mutate() (*corev1.Pod, error) {
	pod := i.pod.DeepCopy()

	// Annotation
	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}

	pod.Annotations[AnnotationMutationStatus] = "mutated"

	// Init container
	if pod.Spec.InitContainers == nil {
		pod.Spec.InitContainers = []corev1.Container{}
	}

	pod.Spec.InitContainers = append(pod.Spec.InitContainers, CreateInitContainers()...)

	// Containers
	if pod.Spec.Containers == nil {
		pod.Spec.Containers = []corev1.Container{}
	}

	for i, container := range pod.Spec.Containers {
		c := container.DeepCopy()

		if err := MutateContainer(c); err != nil {
			return nil, err
		}

		pod.Spec.Containers[i] = *c
	}

	// Volume
	if pod.Spec.Volumes == nil {
		pod.Spec.Volumes = []corev1.Volume{}
	}

	pod.Spec.Volumes = append(pod.Spec.Volumes, CreateVolumes()...)

	return pod, nil
}
