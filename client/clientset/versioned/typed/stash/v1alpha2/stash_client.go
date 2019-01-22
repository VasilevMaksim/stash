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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha2

import (
	v1alpha2 "github.com/appscode/stash/apis/stash/v1alpha2"
	"github.com/appscode/stash/client/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type StashV1alpha2Interface interface {
	RESTClient() rest.Interface
	ActionsGetter
	BackupConfigurationsGetter
	BackupInstancesGetter
	DefaultBackupConfigurationsGetter
	RecoveriesGetter
	RepositoriesGetter
	StashTemplatesGetter
}

// StashV1alpha2Client is used to interact with features provided by the stash.appscode.com group.
type StashV1alpha2Client struct {
	restClient rest.Interface
}

func (c *StashV1alpha2Client) Actions() ActionInterface {
	return newActions(c)
}

func (c *StashV1alpha2Client) BackupConfigurations(namespace string) BackupConfigurationInterface {
	return newBackupConfigurations(c, namespace)
}

func (c *StashV1alpha2Client) BackupInstances(namespace string) BackupInstanceInterface {
	return newBackupInstances(c, namespace)
}

func (c *StashV1alpha2Client) DefaultBackupConfigurations() DefaultBackupConfigurationInterface {
	return newDefaultBackupConfigurations(c)
}

func (c *StashV1alpha2Client) Recoveries(namespace string) RecoveryInterface {
	return newRecoveries(c, namespace)
}

func (c *StashV1alpha2Client) Repositories(namespace string) RepositoryInterface {
	return newRepositories(c, namespace)
}

func (c *StashV1alpha2Client) StashTemplates() StashTemplateInterface {
	return newStashTemplates(c)
}

// NewForConfig creates a new StashV1alpha2Client for the given config.
func NewForConfig(c *rest.Config) (*StashV1alpha2Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &StashV1alpha2Client{client}, nil
}

// NewForConfigOrDie creates a new StashV1alpha2Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *StashV1alpha2Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new StashV1alpha2Client for the given RESTClient.
func New(c rest.Interface) *StashV1alpha2Client {
	return &StashV1alpha2Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha2.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *StashV1alpha2Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
