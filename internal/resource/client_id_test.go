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

func Test_clientIdResourceDataHandler(t *testing.T) {
	tests := resourceIdHandlerTestCases{
		{
			name: "empty",
		},
		{
			name:       "valid",
			inputIdUrl: "https://localhost:8089/services/authentication/users/testuser",
			wantValue:  "https://localhost:8089/services/authentication/users/testuser",
		},
	}

	tests.test(t)
}

func Test_clientIdObjectValueHandler(t *testing.T) {
	tests := objectIdHandlerTestCases{
		{
			name: "empty",
		},
		{
			name:    "invalid resource Id",
			inputId: "invalid",
			// clientId.ManageObjectFunc doesn't actually return any errors, as invalid URLs
			// are assumed to be due to migration from the legacy client.
		},
		{
			name:      "valid resource Id",
			inputId:   "https://localhost:8089/services/authentication/users/testuser",
			wantValue: "https://localhost:8089/services/authentication/users/testuser",
		},
	}

	tests.test(t)
}
