package v1beta1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindBackupConfiguration     = "BackupConfiguration"
	ResourceSingularBackupConfiguration = "backupconfiguration"
	ResourcePluralBackupConfiguration   = "backupconfigurations"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackupConfiguration struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BackupConfigurationSpec `json:"spec,omitempty"`
}

type BackupConfigurationSpec struct {
	Schedule string `json:"schedule,omitempty"`
	// Task specify the Task crd that specifies the steps to take backup
	// +optional
	Task *TaskRef `json:"task,omitempty"`
	// Repository refer to the Repository crd that holds backend information
	Repository core.LocalObjectReference `json:"repository"`
	// Target specify the backup target
	// +optional
	Target *Target `json:"target,omitempty"`
	// RetentionPolicy indicates the policy to follow to clean old backup snapshots
	RetentionPolicy `json:"retentionPolicy,omitempty"`
	// Indicates that the BackupConfiguration is paused from taking backup. Default value is 'false'
	// +optional
	Paused bool `json:"paused,omitempty"`
	// ExecutionEnvironment allow to specify Resources, NodeSelector, Affinity, Toleration, ReadinessProbe etc.
	//+optional
	ExecutionEnvironment *ExecutionEnvironment `json:"executionEnvironment,omitempty"`
	// SecurityScheme allow to specify ServiceAccountName, SecurityContext etc.
	//+optional
	SecurityScheme *SecurityScheme `json:"securityScheme,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackupConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupConfiguration `json:"items,omitempty"`
}