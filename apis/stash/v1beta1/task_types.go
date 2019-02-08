package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core"
)

const (
	ResourceKindTask     = "Task"
	ResourcePluralTask   = "procedures"
	ResourceSingularTask = "procedure"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Task struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TaskSpec `json:"spec,omitempty"`
}

type TaskSpec struct {
	Functions []FunctionSequence `json:"functions,omitempty"`
	// List of volumes that can be mounted by containers belonging to the pod created for this procedure.
	// +optional
	Volumes []core.Volume `json:"volumes,omitempty"`
}

type FunctionSequence struct {
	// Name indicates the name of Action crd
	Name string `json:"name,omitempty"`
	// Inputs specifies the inputs of respective Action
	// +optional
	Inputs map[string]string `json:"inputs,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Task `json:"items,omitempty"`
}
