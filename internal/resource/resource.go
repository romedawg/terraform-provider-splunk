// Copyright 2022 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package resource provides an abstraction between concrete types and schema.ResourceData.
package resource

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// ResourceDataFunc is a function that performs an operation against ResourceData.
type ResourceDataFunc func(*schema.ResourceData) error

// ResourceObjectManager is the interface for types that define ResourceDataFuncs to
// manage the relationship between schema.ResourceData and another object.
type ResourceObjectManager interface {
	ManageResourceFunc() ResourceDataFunc
	ManageObjectFunc() ResourceDataFunc
}

// resourceManagers is a slice of ResourceManagers.
type resourceManagers []ResourceObjectManager

// ManageResourceFunc returns a new ResourceDataFunc that runs each ResourceDataHandler
// defined by resourceManagers to implement ResourceManager.
func (handlers resourceManagers) ManageResourceFunc() ResourceDataFunc {
	return func(d *schema.ResourceData) error {
		for _, handler := range handlers {
			if err := handler.ManageResourceFunc()(d); err != nil {
				return err
			}
		}

		return nil
	}
}

// ManageObjectFunc returns a new ResourceDataFunc that runs each ObjectManagerFunc
// defined by resourceManagers to implement ResourceManager.
func (handlers resourceManagers) ManageObjectFunc() ResourceDataFunc {
	return func(d *schema.ResourceData) error {
		for _, handler := range handlers {
			if err := handler.ManageObjectFunc()(d); err != nil {
				return err
			}
		}

		return nil
	}
}

// ComposeResourceObjectHandler returns a new ResourceObjectHandler that is composed
// of the given handlers.
func ComposeResourceObjectHandler(handlers ...ResourceObjectManager) ResourceObjectManager {
	newHandlers := make(resourceManagers, len(handlers))

	for i, handler := range handlers {
		newHandlers[i] = handler
	}

	return newHandlers
}
