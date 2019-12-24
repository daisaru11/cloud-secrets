package mutating

import (
	"encoding/json"
	"fmt"

	"github.com/daisaru11/cloud-secrets/injector"
	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var systemNamespaces = []string{
	metav1.NamespaceSystem,
	metav1.NamespacePublic,
}

func Mutate(req *v1beta1.AdmissionRequest) *v1beta1.AdmissionResponse {
	var pod corev1.Pod

	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		logrus.Errorf("Failed to unmarshal object to pod: %s", err)
		return errorResponse(req, fmt.Sprintf("Failed to unmarshal object to pod: %s", err))
	}

	logrus.Debugln("Checking namespaces")

	for _, ns := range systemNamespaces {
		if req.Namespace == ns {
			return errorResponse(req, fmt.Sprintf("Cannot mutate objects in system namespaces: %s", req.Namespace))
		}
	}

	i := injector.NewInjector(&pod)

	logrus.Debugln("Checking whether to mutate the object")

	shouldMutate, err := i.ShouldMutate()

	if err != nil {
		logrus.Errorf("Failed to check whether to mutate the object: %s", err)
		return errorResponse(req, fmt.Sprintf("Failed to check whether to mutate the object: %s", err))
	}

	if !shouldMutate {
		return successResponse(req)
	}

	logrus.Debugln("Generating patches")

	patches, err := i.GeneratePatch()

	if err != nil {
		logrus.Errorf("Failed to generate mutation patches: %s", err)
		return errorResponse(req, fmt.Sprintf("Failed to generate mutation patches: %s", err))
	}

	logrus.Debugf("Generated patches: %s", patches)

	if patches == nil {
		return successResponse(req)
	}

	var patchData []byte

	if len(patches) > 0 {
		var err error

		patchData, err = json.Marshal(patches)
		if err != nil {
			logrus.Errorf("Failed to marshal mutation patches: %s", err)
			return errorResponse(req, fmt.Sprintf("Failed to marshal mutation patches: %s", err))
		}
	}

	return successResponseWithPatches(req, patchData)
}

func errorResponse(req *v1beta1.AdmissionRequest, message string) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		UID: req.UID,
		Result: &metav1.Status{
			Message: message,
			Status:  metav1.StatusFailure,
		},
	}
}

func successResponse(req *v1beta1.AdmissionRequest) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Allowed: true,
		UID:     req.UID,
	}
}

func successResponseWithPatches(req *v1beta1.AdmissionRequest, patchData []byte) *v1beta1.AdmissionResponse {
	patchType := v1beta1.PatchTypeJSONPatch

	return &v1beta1.AdmissionResponse{
		Allowed:   true,
		UID:       req.UID,
		PatchType: &patchType,
		Patch:     patchData,
	}
}
