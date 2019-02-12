package v1beta1

import (
	"github.com/appscode/go/encoding/json/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core"
)

const (
	ResourceKindRestoreSession     = "RestoreSession"
	ResourcePluralRestoreSession   = "restoresessions"
	ResourceSingularRestoreSession = "restoresession"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RestoreSession struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RestoreSessionSpec   `json:"spec,omitempty"`
	Status            RestoreSessionStatus `json:"status,omitempty"`
}

type RestoreSessionSpec struct {
	// Repository refer to the Repository crd that hold backend information
	Repository *core.LocalObjectReference `json:"repository,omitempty"`
	// Task specify the Task crd that specifies the steps for recovery process
	// +optional
	Task *TaskRef `json:"task,omitempty"`
	// Snapshot to recover. Default is latest snapshot.
	// +optional
	Snapshot string `json:"snapshot,omitempty"`
	// Target indicates the target where the recovered data will be stored
	// +optional
	Target *Target `json:"target,omitempty"`
	// ExecutionEnvironment allow to specify Resources, NodeSelector, Affinity, Toleration, ReadinessProbe etc.
	//+optional
	ExecutionEnvironment *ExecutionEnvironment `json:"executionEnvironment,omitempty"`
	// SecurityScheme allow to specify SecurityContext
	//+optional
	SecurityScheme *SecurityScheme `json:"securityScheme,omitempty"`
}

type RestorePolicy string

const (
	RestorePolicyIfNotRecovered RestorePolicy = "IfNotRecovered"
	RestorePolicyOnRestart      RestorePolicy = "OnRestart"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RestoreSessionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RestoreSession `json:"items,omitempty"`
}

type RestorePhase string

const (
	RestorePending   RestorePhase = "Pending"
	RestoreRunning   RestorePhase = "Running"
	RestoreSucceeded RestorePhase = "Succeeded"
	RestoreFailed    RestorePhase = "Failed"
	RestoreUnknown   RestorePhase = "Unknown"
)

type RestoreSessionStatus struct {
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
	Phase              RestorePhase   `json:"phase,omitempty"`
	Stats              []RestoreStats `json:"stats,omitempty"`
}

type RestoreStats struct {
	Path     string       `json:"path,omitempty"`
	Phase    RestorePhase `json:"phase,omitempty"`
	Duration string       `json:"duration,omitempty"`
}
