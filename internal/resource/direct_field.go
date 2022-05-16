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

// directField implements the ResourceObjectHandler interface for simple types
// that are set directly between the resource and object.
type directField[T any] struct {
	key   string
	value *T
}

// NewDirectField returns a ResourceObjectHandler that directly associates
// a value with a resource field.
//
// This direct association is generally not sufficient for lists, sets, or maps,
// as the resource value is stored with stored elements as interface{}, not
// the concrete types the referenced value likely contains.
func NewDirectField[T any](key string, value *T) ResourceObjectManager {
	return directField[T]{
		key:   key,
		value: value,
	}
}

// ManageResourceFunc returns the ResourceDataFunc that manages resource data.
func (field directField[T]) ManageResourceFunc() ResourceDataFunc {
	return func(d *schema.ResourceData) error {
		return d.Set(field.key, *field.value)
	}
}

// ManageObjectFunc returns the ResourceDataFunc that manages the object value.
func (field directField[T]) ManageObjectFunc() ResourceDataFunc {
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

		resourceValueT, ok := resourceValueI.(T)
		if !ok {
			return fmt.Errorf("resource: unable to assign resource type %T to object type %T", resourceValueI, *new(T))
		}

		*field.value = resourceValueT

		return nil
	}
}
