package admission

import (
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	mcv1alpha1 "github.com/AccessibleAI/cnvrg-shim/apis/metacloud/v1alpha1"
	"io"
	v1 "k8s.io/api/admission/v1"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	patchClusterDomainTpl = `[{"op":"replace","path":"/spec/clusterDomain","value":"%s"}]`
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

	clusterDomain, err := h.DiscoverClusterDomain(cnvrgApp)
	if err != nil {
		endWithError(err, w)
		return
	}
	resp, err := h.mutationResponse(ar.Request.UID, clusterDomain)
	if err != nil {
		endWithError(err, w)
		return
	}

	endWithOk(resp, w)
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

func (h *AICloudDomainHandler) DiscoverClusterDomain(cap *mlopsv1.CnvrgApp) (clusterDomain string, err error) {

	// do nothing if cluster domain already set
	if cap.Spec.ClusterDomain != "" {
		return cap.Spec.ClusterDomain, nil
	}
	// discover cluster domain based on Domain & DomainPool & Release Name
	return h.clusterDomain(cap.Name)
}

func (h *AICloudDomainHandler) clusterDomain(releaseName string) (string, error) {

	// compose list options based on label selector domainFromReleaseNameDomainpool-pool=<release-name>
	opts, err := h.domainListOptions(releaseName)
	if err != nil {
		return "", err
	}
	// list all the domains who match the selector
	domains := &mcv1alpha1.DomainList{}
	if err := KubeClient().List(context.Background(), domains, opts...); err != nil {
		return "", err
	}
	// return an error in case of empty list
	if len(domains.Items) == 0 {
		return "", fmt.Errorf("empty domains list, unable to detect MLOps clusterDomain")
	}
	clusterDomain := strings.Replace(domains.Items[0].Spec.CommonName, "*.", "", 1)
	// log and return
	klog.Infof("going to use %s for MLOps instance", clusterDomain)
	return clusterDomain, nil

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

func (h *AICloudDomainHandler) mutationResponse(uuid types.UID, clusterDomain string) ([]byte, error) {
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
		Patch:     []byte(fmt.Sprintf(patchClusterDomainTpl, clusterDomain)),
		Result:    &metav1.Status{Message: "ok"},
	}

	return json.Marshal(ar)
}
