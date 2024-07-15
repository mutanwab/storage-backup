package v1beta1

// NOTES: Harvester VM backup & restore is referenced from the Kubevirt's VM snapshot & restore,
// currently, we have decided to use custom VM backup and restore controllers because of the following issues:
// 1. live VM snapshot/backup should be supported, but it is prohibited on the Kubevirt side.
// 2. restore a VM backup to a new VM should be supported.
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:shortName=vmbackup;vmbackups,scope=Namespaced
// +kubebuilder:printcolumn:name="SOURCE_KIND",type=string,JSONPath=`.spec.source.kind`
// +kubebuilder:printcolumn:name="SOURCE_NAME",type=string,JSONPath=`.spec.source.name`
// +kubebuilder:printcolumn:name="TYPE",type=string,JSONPath=`.spec.type`
// +kubebuilder:printcolumn:name="READY_TO_USE",type=boolean,JSONPath=`.status.readyToUse`
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="ERROR",type=date,JSONPath=`.status.error.message`

type StorageBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec StorageBackupSpec `json:"spec"`

	// +optional
	Status *StorageBackupStatus `json:"status,omitempty" default:""`
}

type StorageBackupSpec struct {
	//Source corev1.TypedLocalObjectReference `json:"source"`
	//
	//// +kubebuilder:default:="backup"
	//// +kubebuilder:validation:Enum=backup;snapshot
	//// +kubebuilder:validation:Optional
	//Type BackupType `json:"type,omitempty" default:"backup"`
	VolumeBackups []*VolumeBackup `json:"volumeBackups"`
}

type VolumeBackup struct {
}

// StorageBackupStatus is the status for a VirtualMachineBackup resource
type StorageBackupStatus struct {
	//// +optional
	//SourceUID *types.UID `json:"sourceUID,omitempty"`
	//
	//// +optional
	//CreationTime *metav1.Time `json:"creationTime,omitempty"`
	//
	//// +optional
	//BackupTarget *BackupTarget `json:"backupTarget,omitempty"`
	//
	//// +optional
	//CSIDriverVolumeSnapshotClassNames map[string]string `json:"csiDriverVolumeSnapshotClassNames,omitempty"`
	//
	//// +kubebuilder:validation:Required
	//// SourceSpec contains the vm spec source of the backup target
	//SourceSpec *VirtualMachineSourceSpec `json:"source,omitempty"`
	//
	//// +optional
	//VolumeBackups []VolumeBackup `json:"volumeBackups,omitempty"`
	//
	//// +optional
	//SecretBackups []SecretBackup `json:"secretBackups,omitempty"`
	//
	//// +optional
	//Progress int `json:"progress,omitempty"`
	//
	//// +optional
	//ReadyToUse *bool `json:"readyToUse,omitempty"`
	//
	//// +optional
	//Error *Error `json:"error,omitempty"`
	//
	//// +optional
	//Conditions []Condition `json:"conditions,omitempty"`
}
