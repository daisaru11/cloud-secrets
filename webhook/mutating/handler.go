package mutating

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var deserializer = func() runtime.Decoder {
	codecs := serializer.NewCodecFactory(runtime.NewScheme())
	return codecs.UniversalDeserializer()
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// nolint: funlen
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Handling webhook. Method: %s, URL: %s", r.Method, r.URL)

	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		msg := fmt.Sprintf("Invalid Content-Type: %q", ct)
		code := http.StatusBadRequest
		http.Error(w, msg, code)
		logrus.Warnf("Returns error. Code: %d, Message: %s", code, msg)

		return
	}

	var body []byte

	if r.Body != nil {
		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			msg := fmt.Sprintf("Failed to read the request body")
			code := http.StatusBadRequest
			http.Error(w, msg, code)
			logrus.Warnf("Returns error. Code: %d, Message: %s, Error: %s", code, msg, err)

			return
		}
	}

	if len(body) == 0 {
		msg := fmt.Sprintf("Request body is empty")
		code := http.StatusBadRequest
		http.Error(w, msg, code)
		logrus.Warnf("Returns error. Code: %d, Message: %s", code, msg)

		return
	}

	var admReq v1beta1.AdmissionReview

	var admResp v1beta1.AdmissionReview

	if _, _, err := deserializer().Decode(body, nil, &admReq); err != nil {
		msg := fmt.Sprintf("Failed to decode the admission request")
		code := http.StatusInternalServerError
		http.Error(w, msg, http.StatusInternalServerError)
		logrus.Errorf("Returns error. Code: %d, Message: %s, Error: %s", code, msg, err)

		return
	}

	admResp.Response = Mutate(admReq.Request)
	resp, err := json.Marshal(&admResp)

	if err != nil {
		msg := fmt.Sprintf("Failed to marshal admission response")
		code := http.StatusInternalServerError
		http.Error(w, msg, http.StatusInternalServerError)
		logrus.Errorf("Returns error. Code: %d, Message: %s, Error: %s", code, msg, err)

		return
	}

	if _, err := w.Write(resp); err != nil {
		logrus.Errorf("Failed to write body: %s", err)
	}
}
