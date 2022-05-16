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

import "testing"

func Test_directListField_ResourceDataHandler(t *testing.T) {
	tests := resourceDataHandlerTestCases{
		{
			name:        "invalid key",
			inputKey:    invalidFieldKey,
			inputStruct: newDirectListFieldTestStorage([]string{stringFieldObjectValue}),
			wantError:   true,
		},
		{
			name:        "valid key, type mismatch",
			inputKey:    listIntFieldKey,
			inputStruct: newDirectListFieldTestStorage([]string{stringFieldObjectValue}),
			wantValue:   []interface{}{intFieldResourceValue},
			wantError:   true,
		},
		{
			name:        "valid key, type match (string)",
			inputKey:    listStringFieldKey,
			inputStruct: newDirectListFieldTestStorage([]string{stringFieldObjectValue}),
			wantValue:   []interface{}{stringFieldObjectValue},
		},
		{
			name:        "valid key, type match (int)",
			inputKey:    listStringFieldKey,
			inputStruct: newDirectListFieldTestStorage([]int{intFieldObjectValue}),
			wantValue:   []interface{}{stringFieldResourceValue},
			wantError:   true,
		},
	}

	tests.test(t)
}

func Test_directListField_ObjectHandler(t *testing.T) {
	tests := resourceDataObjectHandlerTestCases{
		{
			name:        "invalid key",
			inputKey:    invalidFieldKey,
			inputHelper: newDirectListFieldTestStorage([]string{stringFieldObjectValue}),
			wantValue:   []string{stringFieldObjectValue},
			wantError:   true,
		},
		{
			name:        "valid key, type mismatch",
			inputKey:    listIntFieldKey,
			inputHelper: newDirectListFieldTestStorage([]string{stringFieldObjectValue}),
			wantValue:   []string{stringFieldObjectValue},
			wantError:   true,
		},
		{
			name:        "valid key, type match (string)",
			inputKey:    listStringFieldKey,
			inputHelper: newDirectListFieldTestStorage([]string{stringFieldObjectValue}),
			wantValue:   []string{stringFieldResourceValue},
		},
		{
			name:        "valid key, type match (int)",
			inputKey:    listIntFieldKey,
			inputHelper: newDirectListFieldTestStorage([]int{intFieldObjectValue}),
			wantValue:   []int{intFieldResourceValue},
		},
		{
			name:        "unset field",
			inputKey:    unsetListStringFieldKey,
			inputHelper: newDirectListFieldTestStorage([]string{stringFieldObjectValue}),
			wantValue:   []string{stringFieldObjectValue},
		},
	}

	tests.test(t)
}
