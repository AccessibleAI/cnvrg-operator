//go:build !ignore_autogenerated

/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BaseConfig) DeepCopyInto(out *BaseConfig) {
	*out = *in
	if in.FeatureFlags != nil {
		in, out := &in.FeatureFlags, &out.FeatureFlags
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BaseConfig.
func (in *BaseConfig) DeepCopy() *BaseConfig {
	if in == nil {
		return nil
	}
	out := new(BaseConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CentralProxy) DeepCopyInto(out *CentralProxy) {
	*out = *in
	out.Limits = in.Limits
	out.Requests = in.Requests
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CentralProxy.
func (in *CentralProxy) DeepCopy() *CentralProxy {
	if in == nil {
		return nil
	}
	out := new(CentralProxy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CentralSSO) DeepCopyInto(out *CentralSSO) {
	*out = *in
	if in.EmailDomain != nil {
		in, out := &in.EmailDomain, &out.EmailDomain
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Limits = in.Limits
	out.Requests = in.Requests
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CentralSSO.
func (in *CentralSSO) DeepCopy() *CentralSSO {
	if in == nil {
		return nil
	}
	out := new(CentralSSO)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CleanupPolicy) DeepCopyInto(out *CleanupPolicy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CleanupPolicy.
func (in *CleanupPolicy) DeepCopy() *CleanupPolicy {
	if in == nil {
		return nil
	}
	out := new(CleanupPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterDomainPrefix) DeepCopyInto(out *ClusterDomainPrefix) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterDomainPrefix.
func (in *ClusterDomainPrefix) DeepCopy() *ClusterDomainPrefix {
	if in == nil {
		return nil
	}
	out := new(ClusterDomainPrefix)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnvrgApp) DeepCopyInto(out *CnvrgApp) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnvrgApp.
func (in *CnvrgApp) DeepCopy() *CnvrgApp {
	if in == nil {
		return nil
	}
	out := new(CnvrgApp)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CnvrgApp) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnvrgAppList) DeepCopyInto(out *CnvrgAppList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CnvrgApp, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnvrgAppList.
func (in *CnvrgAppList) DeepCopy() *CnvrgAppList {
	if in == nil {
		return nil
	}
	out := new(CnvrgAppList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CnvrgAppList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnvrgAppSpec) DeepCopyInto(out *CnvrgAppSpec) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.ControlPlane.DeepCopyInto(&out.ControlPlane)
	out.Registry = in.Registry
	in.Dbs.DeepCopyInto(&out.Dbs)
	in.Networking.DeepCopyInto(&out.Networking)
	in.SSO.DeepCopyInto(&out.SSO)
	out.Tenancy = in.Tenancy
	out.PriorityClass = in.PriorityClass
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnvrgAppSpec.
func (in *CnvrgAppSpec) DeepCopy() *CnvrgAppSpec {
	if in == nil {
		return nil
	}
	out := new(CnvrgAppSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnvrgClusterProvisionerOperator) DeepCopyInto(out *CnvrgClusterProvisionerOperator) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnvrgClusterProvisionerOperator.
func (in *CnvrgClusterProvisionerOperator) DeepCopy() *CnvrgClusterProvisionerOperator {
	if in == nil {
		return nil
	}
	out := new(CnvrgClusterProvisionerOperator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnvrgScheduler) DeepCopyInto(out *CnvrgScheduler) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnvrgScheduler.
func (in *CnvrgScheduler) DeepCopy() *CnvrgScheduler {
	if in == nil {
		return nil
	}
	out := new(CnvrgScheduler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControlPlane) DeepCopyInto(out *ControlPlane) {
	*out = *in
	out.WebApp = in.WebApp
	out.Sidekiq = in.Sidekiq
	out.Searchkiq = in.Searchkiq
	out.Systemkiq = in.Systemkiq
	out.Hyper = in.Hyper
	out.CnvrgScheduler = in.CnvrgScheduler
	in.BaseConfig.DeepCopyInto(&out.BaseConfig)
	out.Ldap = in.Ldap
	out.SMTP = in.SMTP
	out.ObjectStorage = in.ObjectStorage
	out.Nomex = in.Nomex
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControlPlane.
func (in *ControlPlane) DeepCopy() *ControlPlane {
	if in == nil {
		return nil
	}
	out := new(ControlPlane)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Dbs) DeepCopyInto(out *Dbs) {
	*out = *in
	in.Pg.DeepCopyInto(&out.Pg)
	in.Redis.DeepCopyInto(&out.Redis)
	in.Minio.DeepCopyInto(&out.Minio)
	in.Es.DeepCopyInto(&out.Es)
	in.Prom.DeepCopyInto(&out.Prom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Dbs.
func (in *Dbs) DeepCopy() *Dbs {
	if in == nil {
		return nil
	}
	out := new(Dbs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Elastalert) DeepCopyInto(out *Elastalert) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Elastalert.
func (in *Elastalert) DeepCopy() *Elastalert {
	if in == nil {
		return nil
	}
	out := new(Elastalert)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Es) DeepCopyInto(out *Es) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.CleanupPolicy = in.CleanupPolicy
	out.Kibana = in.Kibana
	in.Elastalert.DeepCopyInto(&out.Elastalert)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Es.
func (in *Es) DeepCopy() *Es {
	if in == nil {
		return nil
	}
	out := new(Es)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtraScrapeConfigs) DeepCopyInto(out *ExtraScrapeConfigs) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtraScrapeConfigs.
func (in *ExtraScrapeConfigs) DeepCopy() *ExtraScrapeConfigs {
	if in == nil {
		return nil
	}
	out := new(ExtraScrapeConfigs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Grafana) DeepCopyInto(out *Grafana) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Grafana.
func (in *Grafana) DeepCopy() *Grafana {
	if in == nil {
		return nil
	}
	out := new(Grafana)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPS) DeepCopyInto(out *HTTPS) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPS.
func (in *HTTPS) DeepCopy() *HTTPS {
	if in == nil {
		return nil
	}
	out := new(HTTPS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Hpa) DeepCopyInto(out *Hpa) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Hpa.
func (in *Hpa) DeepCopy() *Hpa {
	if in == nil {
		return nil
	}
	out := new(Hpa)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HugePages) DeepCopyInto(out *HugePages) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HugePages.
func (in *HugePages) DeepCopy() *HugePages {
	if in == nil {
		return nil
	}
	out := new(HugePages)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Hyper) DeepCopyInto(out *Hyper) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Hyper.
func (in *Hyper) DeepCopy() *Hyper {
	if in == nil {
		return nil
	}
	out := new(Hyper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ingress.
func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Jwks) DeepCopyInto(out *Jwks) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Jwks.
func (in *Jwks) DeepCopy() *Jwks {
	if in == nil {
		return nil
	}
	out := new(Jwks)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Kibana) DeepCopyInto(out *Kibana) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Kibana.
func (in *Kibana) DeepCopy() *Kibana {
	if in == nil {
		return nil
	}
	out := new(Kibana)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ldap) DeepCopyInto(out *Ldap) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ldap.
func (in *Ldap) DeepCopy() *Ldap {
	if in == nil {
		return nil
	}
	out := new(Ldap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Limits) DeepCopyInto(out *Limits) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Limits.
func (in *Limits) DeepCopy() *Limits {
	if in == nil {
		return nil
	}
	out := new(Limits)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Minio) DeepCopyInto(out *Minio) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Minio.
func (in *Minio) DeepCopy() *Minio {
	if in == nil {
		return nil
	}
	out := new(Minio)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Networking) DeepCopyInto(out *Networking) {
	*out = *in
	out.Ingress = in.Ingress
	out.HTTPS = in.HTTPS
	in.Proxy.DeepCopyInto(&out.Proxy)
	out.ClusterDomainPrefix = in.ClusterDomainPrefix
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Networking.
func (in *Networking) DeepCopy() *Networking {
	if in == nil {
		return nil
	}
	out := new(Networking)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Nomex) DeepCopyInto(out *Nomex) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nomex.
func (in *Nomex) DeepCopy() *Nomex {
	if in == nil {
		return nil
	}
	out := new(Nomex)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectStorage) DeepCopyInto(out *ObjectStorage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectStorage.
func (in *ObjectStorage) DeepCopy() *ObjectStorage {
	if in == nil {
		return nil
	}
	out := new(ObjectStorage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pg) DeepCopyInto(out *Pg) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	out.HugePages = in.HugePages
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pg.
func (in *Pg) DeepCopy() *Pg {
	if in == nil {
		return nil
	}
	out := new(Pg)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pki) DeepCopyInto(out *Pki) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pki.
func (in *Pki) DeepCopy() *Pki {
	if in == nil {
		return nil
	}
	out := new(Pki)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PriorityClass) DeepCopyInto(out *PriorityClass) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PriorityClass.
func (in *PriorityClass) DeepCopy() *PriorityClass {
	if in == nil {
		return nil
	}
	out := new(PriorityClass)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Prom) DeepCopyInto(out *Prom) {
	*out = *in
	if in.ExtraScrapeConfigs != nil {
		in, out := &in.ExtraScrapeConfigs, &out.ExtraScrapeConfigs
		*out = make([]*ExtraScrapeConfigs, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ExtraScrapeConfigs)
				**out = **in
			}
		}
	}
	out.Grafana = in.Grafana
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Prom.
func (in *Prom) DeepCopy() *Prom {
	if in == nil {
		return nil
	}
	out := new(Prom)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Proxy) DeepCopyInto(out *Proxy) {
	*out = *in
	if in.HttpProxy != nil {
		in, out := &in.HttpProxy, &out.HttpProxy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.HttpsProxy != nil {
		in, out := &in.HttpsProxy, &out.HttpsProxy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NoProxy != nil {
		in, out := &in.NoProxy, &out.NoProxy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Proxy.
func (in *Proxy) DeepCopy() *Proxy {
	if in == nil {
		return nil
	}
	out := new(Proxy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Redis) DeepCopyInto(out *Redis) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Redis.
func (in *Redis) DeepCopy() *Redis {
	if in == nil {
		return nil
	}
	out := new(Redis)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Registry) DeepCopyInto(out *Registry) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Registry.
func (in *Registry) DeepCopy() *Registry {
	if in == nil {
		return nil
	}
	out := new(Registry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Requests) DeepCopyInto(out *Requests) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Requests.
func (in *Requests) DeepCopy() *Requests {
	if in == nil {
		return nil
	}
	out := new(Requests)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SMTP) DeepCopyInto(out *SMTP) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SMTP.
func (in *SMTP) DeepCopy() *SMTP {
	if in == nil {
		return nil
	}
	out := new(SMTP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SSO) DeepCopyInto(out *SSO) {
	*out = *in
	out.Pki = in.Pki
	in.Jwks.DeepCopyInto(&out.Jwks)
	in.Central.DeepCopyInto(&out.Central)
	in.Proxy.DeepCopyInto(&out.Proxy)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SSO.
func (in *SSO) DeepCopy() *SSO {
	if in == nil {
		return nil
	}
	out := new(SSO)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Searchkiq) DeepCopyInto(out *Searchkiq) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	out.Hpa = in.Hpa
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Searchkiq.
func (in *Searchkiq) DeepCopy() *Searchkiq {
	if in == nil {
		return nil
	}
	out := new(Searchkiq)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Sidekiq) DeepCopyInto(out *Sidekiq) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	out.Hpa = in.Hpa
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Sidekiq.
func (in *Sidekiq) DeepCopy() *Sidekiq {
	if in == nil {
		return nil
	}
	out := new(Sidekiq)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Status) DeepCopyInto(out *Status) {
	*out = *in
	if in.StackReadiness != nil {
		in, out := &in.StackReadiness, &out.StackReadiness
		*out = make(map[string]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Status.
func (in *Status) DeepCopy() *Status {
	if in == nil {
		return nil
	}
	out := new(Status)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Systemkiq) DeepCopyInto(out *Systemkiq) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	out.Hpa = in.Hpa
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Systemkiq.
func (in *Systemkiq) DeepCopy() *Systemkiq {
	if in == nil {
		return nil
	}
	out := new(Systemkiq)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tenancy) DeepCopyInto(out *Tenancy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tenancy.
func (in *Tenancy) DeepCopy() *Tenancy {
	if in == nil {
		return nil
	}
	out := new(Tenancy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebApp) DeepCopyInto(out *WebApp) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	out.Hpa = in.Hpa
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebApp.
func (in *WebApp) DeepCopy() *WebApp {
	if in == nil {
		return nil
	}
	out := new(WebApp)
	in.DeepCopyInto(out)
	return out
}
