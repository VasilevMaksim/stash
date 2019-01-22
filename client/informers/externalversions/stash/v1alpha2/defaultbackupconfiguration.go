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
	time "time"

	stashv1alpha2 "github.com/appscode/stash/apis/stash/v1alpha2"
	versioned "github.com/appscode/stash/client/clientset/versioned"
	internalinterfaces "github.com/appscode/stash/client/informers/externalversions/internalinterfaces"
	v1alpha2 "github.com/appscode/stash/client/listers/stash/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// DefaultBackupConfigurationInformer provides access to a shared informer and lister for
// DefaultBackupConfigurations.
type DefaultBackupConfigurationInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha2.DefaultBackupConfigurationLister
}

type defaultBackupConfigurationInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewDefaultBackupConfigurationInformer constructs a new informer for DefaultBackupConfiguration type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDefaultBackupConfigurationInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDefaultBackupConfigurationInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredDefaultBackupConfigurationInformer constructs a new informer for DefaultBackupConfiguration type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDefaultBackupConfigurationInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.StashV1alpha2().DefaultBackupConfigurations().List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.StashV1alpha2().DefaultBackupConfigurations().Watch(options)
			},
		},
		&stashv1alpha2.DefaultBackupConfiguration{},
		resyncPeriod,
		indexers,
	)
}

func (f *defaultBackupConfigurationInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDefaultBackupConfigurationInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *defaultBackupConfigurationInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&stashv1alpha2.DefaultBackupConfiguration{}, f.defaultInformer)
}

func (f *defaultBackupConfigurationInformer) Lister() v1alpha2.DefaultBackupConfigurationLister {
	return v1alpha2.NewDefaultBackupConfigurationLister(f.Informer().GetIndexer())
}
