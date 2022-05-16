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

package resource

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// directListField implements the ResourceObjectHandler for simple types
// that are set directly between the resource and object.
type directListField[T any] struct {
	value *[]T
	key   string
}

// NewDirectListField returns a ResourceObjectHandler that directly associates
// a list (slice) value with a list resource field.
func NewDirectListField[T any](key string, value *[]T) ResourceObjectManager {
	return directListField[T]{
		key:   key,
		value: value,
	}
}

// ManageResourceFunc returns the ResourceDataFunc that manages resource data.
func (field directListField[T]) ManageResourceFunc() ResourceDataFunc {
	return func(d *schema.ResourceData) error {
		return d.Set(field.key, *field.value)
	}
}

// ManageObjectFunc returns the ResourceDataFunc that manages the object value.
func (field directListField[T]) ManageObjectFunc() ResourceDataFunc {
	return func(d *schema.ResourceData) error {
		resourceValueI, ok := d.GetOk(field.key)
		// nil interface = unknown key
		if resourceValueI == nil {
			return fmt.Errorf("resource: likely unknown key %q", field.key)
		}
		// ok==false = unset, ignore
		if !ok {
			return nil
		}

		resourceValue, ok := resourceValueI.([]interface{})
		if !ok {
			return fmt.Errorf("resource: key %q not a slice (%T)", field.key, resourceValueI)
		}

		newValues := make([]T, len(resourceValue))
		for i, resourceValue := range resourceValue {
			newValue, ok := resourceValue.(T)
			if !ok {
				return fmt.Errorf("resource: key %q not a slice of type %T", field.key, *new(T))
			}

			newValues[i] = newValue
		}

		*field.value = newValues

		return nil
	}
}
