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
)

func Test_directField_ResourceDataHandler(t *testing.T) {
	tests := resourceDataHandlerTestCases{
		{
			name:        "invalid key",
			inputKey:    invalidFieldKey,
			inputStruct: newDirectFieldTestStorage(stringFieldObjectValue),
			wantError:   true,
		},
		{
			name:        "valid key, type mismatch",
			inputKey:    intFieldKey,
			inputStruct: newDirectFieldTestStorage(stringFieldObjectValue),
			wantValue:   intFieldResourceValue,
			wantError:   true,
		},
		{
			name:        "valid key, type match (string)",
			inputKey:    stringFieldKey,
			inputStruct: newDirectFieldTestStorage(stringFieldObjectValue),
			wantValue:   stringFieldObjectValue,
		},
		{
			name:        "valid key, type match (int)",
			inputKey:    intFieldKey,
			inputStruct: newDirectFieldTestStorage(intFieldObjectValue),
			wantValue:   intFieldObjectValue,
		},
		{
			// this is a confusing test, because a schema type List won't work for the object handler,
			// but it does work for the resource handler.
			name:        "list of strings",
			inputKey:    listStringFieldKey,
			inputStruct: newDirectFieldTestStorage([]string{stringFieldObjectValue}),
			wantValue:   []interface{}{stringFieldObjectValue},
		},
	}

	tests.test(t)
}

func Test_directField_ObjectHandler(t *testing.T) {
	tests := resourceDataObjectHandlerTestCases{
		{
			name:        "invalid key",
			inputKey:    invalidFieldKey,
			inputHelper: newDirectFieldTestStorage(stringFieldObjectValue),
			wantValue:   stringFieldObjectValue,
			wantError:   true,
		},
		{
			name:        "valid key, type mismatch",
			inputKey:    intFieldKey,
			inputHelper: newDirectFieldTestStorage(stringFieldObjectValue),
			wantValue:   stringFieldObjectValue,
			wantError:   true,
		},
		{
			name:        "valid key, type match (string)",
			inputKey:    stringFieldKey,
			inputHelper: newDirectFieldTestStorage(stringFieldObjectValue),
			wantValue:   stringFieldResourceValue,
		},
		{
			name:        "valid key, type match (int)",
			inputKey:    intFieldKey,
			inputHelper: newDirectFieldTestStorage(intFieldObjectValue),
			wantValue:   intFieldResourceValue,
		},
		{
			name:        "list of strings",
			inputKey:    listStringFieldKey,
			inputHelper: newDirectFieldTestStorage([]string{stringFieldObjectValue}),
			wantValue:   []string{stringFieldObjectValue},
			wantError:   true,
		},
		{
			name:        "unset field",
			inputKey:    unsetStringFieldKey,
			inputHelper: newDirectFieldTestStorage(stringFieldObjectValue),
			wantValue:   stringFieldObjectValue,
		},
	}

	tests.test(t)
}
