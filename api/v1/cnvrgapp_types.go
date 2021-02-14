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

package v1

import (
	"github.com/cnvrg-operator/pkg/pg"
	"github.com/cnvrg-operator/pkg/storage"
	"github.com/cnvrg-operator/pkg/tenancy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CnvrgAppSpec defines the desired state of CnvrgApp
type CnvrgAppSpec struct {
	Pg      pg.Pg           `json:"pg,omitempty"`
	Storage storage.Storage `json:"storage,omitempty"`
	Tenancy tenancy.Tenancy `json:"tenancy,omitempty"`
}

// CnvrgAppStatus defines the observed state of CnvrgApp
type CnvrgAppStatus struct {
	Message string `json:"message,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// CnvrgApp is the Schema for the cnvrgapps API
type CnvrgApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgAppSpec   `json:"spec,omitempty"`
	Status CnvrgAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CnvrgAppList contains a list of CnvrgApp
type CnvrgAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CnvrgApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CnvrgApp{}, &CnvrgAppList{})
}

func DefaultCnvrgAppSpec() CnvrgAppSpec {
	return CnvrgAppSpec{
		Pg:      pg.Defaults(),
		Storage: storage.Defaults(),
		Tenancy: tenancy.Defaults(),
	}
}
