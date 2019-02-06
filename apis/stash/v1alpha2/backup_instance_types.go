package v1alpha2

import (
	"github.com/appscode/go/encoding/json/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindBackupInstance     = "BackupInstance"
	ResourcePluralBackupInstance   = "backupinstances"
	ResourceSingularBackupInstance = "backupinstance"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackupInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BackupInstanceSpec   `json:"spec,omitempty"`
	Status            BackupInstanceStatus `json:"status,omitempty"`
}

type BackupInstanceSpec struct {
	// TargetBackupConfiguration indicates the respective BackupConfiguration crd for target backup
	TargetBackupConfiguration string `json:"targetBackupConfiguration"`
}

type BackupInstancePhase string

const (
	BackupInstancePending   BackupInstancePhase = "Pending"
	BackupInstanceRunning   BackupInstancePhase = "Running"
	BackupInstanceSucceeded BackupInstancePhase = "Succeeded"
	BackupInstanceFailed    BackupInstancePhase = "Failed"
	BackupInstanceUnknown   BackupInstancePhase = "Unknown"
)

type BackupInstanceStatus struct {
	// ObservedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
	// Phase indicates the phase of the backup process for this BackupInstance
	Phase BackupInstancePhase `json:"phase,omitempty"`
	Stats BackupStats         `json:"stats,omitempty"`
}

type BackupStats struct {
	// Snapshot indicates the name of the backup snapshot created in this backup session
	Snapshot string `json:"snapshot,omitempty"`
	// Size indicates the size of target data to backup
	Size string `json:"size,omitempty"`
	// Uploaded indicates size of data uploaded to backend in this backup session
	Uploaded string `json:"uploaded,omitempty"`
	// ProcessingTime indicates time taken to process the target data
	ProcessingTime string `json:"processingTime,omitempty"`
	// FileStats shows statistics of files of backup session
	FileStats FileStats `json:"fileStats,omitempty"`
}

type FileStats struct {
	// TotalFiles shows total number of files that has been backed up
	TotalFiles *int `json:"totalFiles,omitempty"`
	// NewFiles shows total number of new files that has been created since last backup
	NewFiles *int `json:"newFiles,omitempty"`
	// ModifiedFiles shows total number of files that has been modified since last backup
	ModifiedFiles *int `json:"modifiedFiles,omitempty"`
	// UnmodifiedFiles shows total number of files that has not been changed since last backup
	UnmodifiedFiles *int `json:"unmodifiedFiles,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackupInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupInstance `json:"items,omitempty"`
}
