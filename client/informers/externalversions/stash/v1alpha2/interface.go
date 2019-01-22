/*
Copyright 2019 The Stash Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha2

import (
	internalinterfaces "github.com/appscode/stash/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Actions returns a ActionInformer.
	Actions() ActionInformer
	// BackupConfigurations returns a BackupConfigurationInformer.
	BackupConfigurations() BackupConfigurationInformer
	// BackupInstances returns a BackupInstanceInformer.
	BackupInstances() BackupInstanceInformer
	// BackupTemplates returns a BackupTemplateInformer.
	BackupTemplates() BackupTemplateInformer
	// DefaultBackupConfigurations returns a DefaultBackupConfigurationInformer.
	DefaultBackupConfigurations() DefaultBackupConfigurationInformer
	// Recoveries returns a RecoveryInformer.
	Recoveries() RecoveryInformer
	// Repositories returns a RepositoryInformer.
	Repositories() RepositoryInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Actions returns a ActionInformer.
func (v *version) Actions() ActionInformer {
	return &actionInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// BackupConfigurations returns a BackupConfigurationInformer.
func (v *version) BackupConfigurations() BackupConfigurationInformer {
	return &backupConfigurationInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// BackupInstances returns a BackupInstanceInformer.
func (v *version) BackupInstances() BackupInstanceInformer {
	return &backupInstanceInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// BackupTemplates returns a BackupTemplateInformer.
func (v *version) BackupTemplates() BackupTemplateInformer {
	return &backupTemplateInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// DefaultBackupConfigurations returns a DefaultBackupConfigurationInformer.
func (v *version) DefaultBackupConfigurations() DefaultBackupConfigurationInformer {
	return &defaultBackupConfigurationInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Recoveries returns a RecoveryInformer.
func (v *version) Recoveries() RecoveryInformer {
	return &recoveryInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Repositories returns a RepositoryInformer.
func (v *version) Repositories() RepositoryInformer {
	return &repositoryInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
