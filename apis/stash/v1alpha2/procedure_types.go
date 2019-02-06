package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core"
)

const (
	ResourceKindProcedure     = "Procedure"
	ResourcePluralProcedure   = "procedures"
	ResourceSingularProcedure = "procedure"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Procedure struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProcedureSpec `json:"spec,omitempty"`
}

type ProcedureSpec struct {
	Actions []Steps `json:"actions,omitempty"`
	// List of volumes that can be mounted by containers belonging to the pod created for this procedure.
	// +optional
	Volumes []core.Volume `json:"volumes,omitempty"`
}

type Steps struct {
	// Name indicates the name of Action crd
	Name string `json:"name,omitempty"`
	// Inputs specifies the inputs of respective Action
	// +optional
	Inputs map[string]string `json:"inputs,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ProcedureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Procedure `json:"items,omitempty"`
}
