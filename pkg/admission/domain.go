package admission

import (
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	simetadata "github.com/AccessibleAI/cnvrg-shim/pkg/serviceinstaller/metadata"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	v1 "k8s.io/api/admission/v1"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type AICloudDomainHandler struct {
	deserializer runtime.Decoder
}

func NewAICloudDomainHandler() *AICloudDomainHandler {
	return &AICloudDomainHandler{
		deserializer: serializer.NewCodecFactory(
			runtime.NewScheme()).
			UniversalDeserializer(),
	}
}

func (h *AICloudDomainHandler) Handler(w http.ResponseWriter, r *http.Request) {

	ar, err := h.admissionReviewDecode(h.body(r))
	if err != nil {
		endWithError(err, w)
		return
	}

	cnvrgApp, err := h.cnvrgAppDecode(ar.Request.Object.Raw)
	if err != nil {
		endWithError(err, w)
		return
	}

	patchBody, err := h.patchBody(cnvrgApp)
	if err != nil {
		endWithError(err, w)
		return
	}

	resp, err := h.mutationResponse(ar.Request.UID, patchBody)
	if err != nil {
		endWithError(err, w)
		return
	}

	endWithOk(resp, w)
}

func (h *AICloudDomainHandler) patchBody(cap *mlopsv1.CnvrgApp) (string, error) {
	patchCnvrgAppTpl := `[
			{"op":"replace","path":"/spec/clusterDomain","value":"%s"},
			{"op":"replace","path":"/spec/networking/ingress/istioIngressSelectorKey","value":"%s"},
			{"op":"replace","path":"/spec/networking/ingress/istioIngressSelectorValue","value":"%s"}]`

	serviceInstanceMetadata, err := h.ServiceInstanceMetadata(cap.Namespace)
	if err != nil {
		return "", err
	}
	clusterDomain, err := h.ClusterDomain(cap, serviceInstanceMetadata)
	if err != nil {
		return "", err
	}

	patch := fmt.Sprintf(patchCnvrgAppTpl,
		clusterDomain,
		serviceInstanceMetadata.Ingress.SelectorKey,
		serviceInstanceMetadata.Ingress.SelectorValue,
	)
	zap.S().Info(patch)
	return patch, nil
}

func (h *AICloudDomainHandler) HookCfg(ns, svc string, caBundle []byte) *admissionv1.MutatingWebhookConfiguration {

	hookName := fmt.Sprintf("%s.%s.svc", svc, ns)
	path := h.HandlerPath()
	failPolicy := admissionv1.Fail
	nsScope := admissionv1.NamespacedScope
	sideEffect := admissionv1.SideEffectClassNone

	return &admissionv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{Name: hookName},
		Webhooks: []admissionv1.MutatingWebhook{
			{
				Name: hookName,
				NamespaceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"name": ns},
				},
				ClientConfig: admissionv1.WebhookClientConfig{
					Service: &admissionv1.ServiceReference{
						Namespace: ns,
						Name:      svc,
						Path:      &path,
					},
					CABundle: caBundle,
				},
				Rules: []admissionv1.RuleWithOperations{
					{
						Operations: []admissionv1.OperationType{
							admissionv1.Create,
						},
						Rule: admissionv1.Rule{
							APIGroups:   []string{"mlops.cnvrg.io"},
							APIVersions: []string{"v1"},
							Resources:   []string{"cnvrgapps"},
							Scope:       &nsScope,
						},
					},
				},
				FailurePolicy:           &failPolicy,
				SideEffects:             &sideEffect,
				AdmissionReviewVersions: []string{"v1"},
			},
		},
	}
}

func (h *AICloudDomainHandler) HandlerPath() string {
	return "/aicloud/domain-discovery"
}

func (h *AICloudDomainHandler) body(r *http.Request) (body []byte) {
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	return
}

func (h *AICloudDomainHandler) admissionReviewDecode(body []byte) (*v1.AdmissionReview, error) {
	ar := &v1.AdmissionReview{}
	_, _, err := h.deserializer.Decode(body, nil, ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (h *AICloudDomainHandler) cnvrgAppDecode(b []byte) (*mlopsv1.CnvrgApp, error) {

	cnvrgApp := &mlopsv1.CnvrgApp{}

	if err := json.Unmarshal(b, cnvrgApp); err != nil {
		return nil, err
	}
	return cnvrgApp, nil
}

func (h *AICloudDomainHandler) ClusterDomain(cap *mlopsv1.CnvrgApp, siMeta *simetadata.ServiceInstanceMetadata) (clusterDomain string, err error) {

	// do nothing if cluster domain already set
	if cap.Spec.ClusterDomain != "" {
		return cap.Spec.ClusterDomain, nil
	}
	return strings.Replace(siMeta.Domain, "*.", "", 1), nil
}

func (h *AICloudDomainHandler) ServiceInstanceMetadata(releaseNs string) (*simetadata.ServiceInstanceMetadata, error) {
	aiCloudMetadataCm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Namespace: releaseNs, Name: "serviceinstancemetadata"},
	}
	cmKey := types.NamespacedName{Namespace: releaseNs, Name: "serviceinstancemetadata"}
	if err := KubeClient().Get(context.Background(), cmKey, aiCloudMetadataCm, []client.GetOption{}...); err != nil {
		return nil, err
	}

	instanceMetadata := &simetadata.ServiceInstanceMetadata{}
	if err := yaml.Unmarshal([]byte(aiCloudMetadataCm.Data[simetadata.IngressKeyName]), instanceMetadata); err != nil {
		return nil, err
	}

	return instanceMetadata, nil

}

func (h *AICloudDomainHandler) domainListOptions(releaseName string) ([]client.ListOption, error) {

	selector, err := labels.Parse(fmt.Sprintf("domain-pool=%s", releaseName))
	if err != nil {
		return nil, err
	}
	return []client.ListOption{
		&client.ListOptions{
			LabelSelector: selector,
			Limit:         1,
		},
	}, nil

}

func (h *AICloudDomainHandler) mutationResponse(uuid types.UID, patchBody string) ([]byte, error) {
	pt := v1.PatchTypeJSONPatch
	ar := &v1.AdmissionReview{}
	ar.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "admission.k8s.io",
		Version: "v1",
		Kind:    "AdmissionReview",
	})
	ar.Response = &v1.AdmissionResponse{
		UID:       uuid,
		Allowed:   true,
		PatchType: &pt,
		Patch:     []byte(patchBody),
		Result:    &metav1.Status{Message: "ok"},
	}

	return json.Marshal(ar)
}
