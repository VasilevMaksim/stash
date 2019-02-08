package v1beta1

import (
	"github.com/appscode/go/encoding/json/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core"
)

const (
	ResourceKindRecoveryConfiguration     = "RecoveryConfiguration"
	ResourcePluralRecoveryConfiguration   = "recoveryconfigurations"
	ResourceSingularRecoveryConfiguration = "recoveryconfiguration"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RecoveryConfiguration struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RecoveryConfigurationSpec   `json:"spec,omitempty"`
	Status            RecoveryConfigurationStatus `json:"status,omitempty"`
}

type RecoveryConfigurationSpec struct {
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

type RecoveryPolicy string

const (
	RecoveryPolicyIfNotRecovered RecoveryPolicy = "IfNotRecovered"
	RecoveryPolicyOnRestart      RecoveryPolicy = "OnRestart"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RecoveryConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecoveryConfiguration `json:"items,omitempty"`
}

type RecoveryPhase string

const (
	RecoveryPending   RecoveryPhase = "Pending"
	RecoveryRunning   RecoveryPhase = "Running"
	RecoverySucceeded RecoveryPhase = "Succeeded"
	RecoveryFailed    RecoveryPhase = "Failed"
	RecoveryUnknown   RecoveryPhase = "Unknown"
)

type RecoveryConfigurationStatus struct {
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
	Phase              RecoveryPhase  `json:"phase,omitempty"`
	Stats              []RestoreStats `json:"stats,omitempty"`
}

type RestoreStats struct {
	Path     string        `json:"path,omitempty"`
	Phase    RecoveryPhase `json:"phase,omitempty"`
	Duration string        `json:"duration,omitempty"`
}
