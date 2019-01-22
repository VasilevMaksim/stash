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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha2

import (
	v1alpha2 "github.com/appscode/stash/apis/stash/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ActionLister helps list Actions.
type ActionLister interface {
	// List lists all Actions in the indexer.
	List(selector labels.Selector) (ret []*v1alpha2.Action, err error)
	// Get retrieves the Action from the index for a given name.
	Get(name string) (*v1alpha2.Action, error)
	ActionListerExpansion
}

// actionLister implements the ActionLister interface.
type actionLister struct {
	indexer cache.Indexer
}

// NewActionLister returns a new ActionLister.
func NewActionLister(indexer cache.Indexer) ActionLister {
	return &actionLister{indexer: indexer}
}

// List lists all Actions in the indexer.
func (s *actionLister) List(selector labels.Selector) (ret []*v1alpha2.Action, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.Action))
	})
	return ret, err
}

// Get retrieves the Action from the index for a given name.
func (s *actionLister) Get(name string) (*v1alpha2.Action, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha2.Resource("action"), name)
	}
	return obj.(*v1alpha2.Action), nil
}
