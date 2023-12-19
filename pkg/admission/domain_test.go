package admission

import (
	"context"
	"github.com/AccessibleAI/cnvrg-shim/apis/metacloud/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"testing"
)

func TestDomainListByLabelsSelector(t *testing.T) {
	domainPoolName := "foo-bar-domainFromReleaseNameDomainpool-pool"
	testDomain := &v1alpha1.Domain{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "foo-bar",
			Labels: map[string]string{"domainFromReleaseNameDomainpool-pool": domainPoolName},
		},
	}
	kc := KubeClient()
	if got := kc.Create(context.Background(), testDomain, []client.CreateOption{}...); got != nil {
		t.Errorf(got.Error())
	}

	discoveryHandler := NewAICloudDomainHandler()
	opts, got := discoveryHandler.domainListOptions(domainPoolName)
	if got != nil {
		t.Errorf(got.Error())
	}
	domains := &v1alpha1.DomainList{}
	if got := kc.List(context.Background(), domains, opts...); got != nil {
		t.Errorf(got.Error())
	}
	want := 1
	totalDomainsNum := len(domains.Items)
	if totalDomainsNum != want {
		t.Errorf("got wrong number of domains, expted %d got %d", want, totalDomainsNum)
	}

	defer func() {
		if got := kc.Delete(context.Background(), testDomain, []client.DeleteOption{}...); got != nil {
			t.Errorf(got.Error())
		}
	}()

}
