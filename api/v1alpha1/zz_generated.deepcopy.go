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

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetaStorageProvisioner) DeepCopyInto(out *MetaStorageProvisioner) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetaStorageProvisioner.
func (in *MetaStorageProvisioner) DeepCopy() *MetaStorageProvisioner {
	if in == nil {
		return nil
	}
	out := new(MetaStorageProvisioner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MetaStorageProvisioner) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetaStorageProvisionerList) DeepCopyInto(out *MetaStorageProvisionerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MetaStorageProvisioner, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetaStorageProvisionerList.
func (in *MetaStorageProvisionerList) DeepCopy() *MetaStorageProvisionerList {
	if in == nil {
		return nil
	}
	out := new(MetaStorageProvisionerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MetaStorageProvisionerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NFSProvisioner) DeepCopyInto(out *NFSProvisioner) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NFSProvisioner.
func (in *NFSProvisioner) DeepCopy() *NFSProvisioner {
	if in == nil {
		return nil
	}
	out := new(NFSProvisioner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProvisionerType) DeepCopyInto(out *ProvisionerType) {
	*out = *in
	if in.NFSProvisioner != nil {
		in, out := &in.NFSProvisioner, &out.NFSProvisioner
		*out = new(NFSProvisioner)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProvisionerType.
func (in *ProvisionerType) DeepCopy() *ProvisionerType {
	if in == nil {
		return nil
	}
	out := new(ProvisionerType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageProvisionerSpec) DeepCopyInto(out *StorageProvisionerSpec) {
	*out = *in
	in.ProvisionerType.DeepCopyInto(&out.ProvisionerType)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageProvisionerSpec.
func (in *StorageProvisionerSpec) DeepCopy() *StorageProvisionerSpec {
	if in == nil {
		return nil
	}
	out := new(StorageProvisionerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageProvisionerStatus) DeepCopyInto(out *StorageProvisionerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageProvisionerStatus.
func (in *StorageProvisionerStatus) DeepCopy() *StorageProvisionerStatus {
	if in == nil {
		return nil
	}
	out := new(StorageProvisionerStatus)
	in.DeepCopyInto(out)
	return out
}