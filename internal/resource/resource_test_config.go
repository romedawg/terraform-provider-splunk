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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// test values used in stored objects and fake schema.ResourceData schema/data.
const (
	stringFieldKey          string = "string_field"
	unsetStringFieldKey     string = "unset_string_field"
	intFieldKey             string = "int_field"
	listStringFieldKey      string = "list_string_field"
	unsetListStringFieldKey string = "unset_list_string_field"
	listIntFieldKey         string = "list_int_field"
	invalidFieldKey         string = "invalid_field"

	stringFieldResourceValue string = "resource_value"
	stringFieldObjectValue   string = "object_value"

	intFieldResourceValue int = 1
	intFieldObjectValue   int = 2
)

// resourceData returns a new schema.ResourceData with a schema and values suitable
// for testing.
func resourceData(t *testing.T) *schema.ResourceData {
	return schema.TestResourceDataRaw(
		t,
		map[string]*schema.Schema{
			stringFieldKey: {
				Type: schema.TypeString,
			},
			unsetStringFieldKey: {
				Type: schema.TypeString,
			},
			intFieldKey: {
				Type: schema.TypeInt,
			},
			listStringFieldKey: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			unsetListStringFieldKey: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			listIntFieldKey: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
		map[string]interface{}{
			stringFieldKey: stringFieldResourceValue,
			intFieldKey:    intFieldResourceValue,
			listStringFieldKey: []interface{}{
				stringFieldResourceValue,
			},
			listIntFieldKey: []interface{}{
				intFieldResourceValue,
			},
		},
	)
}
