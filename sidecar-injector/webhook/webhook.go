package hook

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io

// SidecarInjector annotates Pods
type SidecarInjector struct {
	Name          string
	Client        client.Client
	decoder       *admission.Decoder
	Annotations	  map[string]string
}

func shoudInject(pod *corev1.Pod) bool {
	shouldInjectSidecar, err := strconv.ParseBool(pod.Annotations["inject-sidecar"])

	if err != nil {
		shouldInjectSidecar = true
	}

	log.Info("Should Inject: ", shouldInjectSidecar)

	return shouldInjectSidecar
}

// SidecarInjector adds an annotation to every incoming pods.
func (si *SidecarInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := si.decoder.Decode(req, pod)
	if err != nil {
		log.Info("Sdecar-Injector: cannot decode")
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}

	shoudInjectSidecar := shoudInject(pod)

	if shoudInjectSidecar {
		log.Info("Injecting sidecar...")

		pod.Annotations = si.Annotations

		log.Info("Sidecar ", si.Name, " injected.")
	} else {
		log.Info("Inject not needed.")
	}

	marshaledPod, err := json.Marshal(pod)

	if err != nil {
		log.Info("Sdecar-Injector: cannot marshal")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// SidecarInjector implements admission.DecoderInjector.
// A decoder will be automatically inj1ected.

// InjectDecoder injects the decoder.
func (si *SidecarInjector) InjectDecoder(d *admission.Decoder) error {
	si.decoder = d
	return nil
}
