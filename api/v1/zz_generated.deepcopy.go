// +build !ignore_autogenerated

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
func (in *CnvrgApp) DeepCopyInto(out *CnvrgApp) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
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
	out.Pg = in.Pg
	out.Storage = in.Storage
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
func (in *CnvrgAppStatus) DeepCopyInto(out *CnvrgAppStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnvrgAppStatus.
func (in *CnvrgAppStatus) DeepCopy() *CnvrgAppStatus {
	if in == nil {
		return nil
	}
	out := new(CnvrgAppStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Hostpath) DeepCopyInto(out *Hostpath) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Hostpath.
func (in *Hostpath) DeepCopy() *Hostpath {
	if in == nil {
		return nil
	}
	out := new(Hostpath)
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
func (in *Nfs) DeepCopyInto(out *Nfs) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nfs.
func (in *Nfs) DeepCopy() *Nfs {
	if in == nil {
		return nil
	}
	out := new(Nfs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pg) DeepCopyInto(out *Pg) {
	*out = *in
	out.HugePages = in.HugePages
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
func (in *Storage) DeepCopyInto(out *Storage) {
	*out = *in
	out.Hostpath = in.Hostpath
	out.Nfs = in.Nfs
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Storage.
func (in *Storage) DeepCopy() *Storage {
	if in == nil {
		return nil
	}
	out := new(Storage)
	in.DeepCopyInto(out)
	return out
}
