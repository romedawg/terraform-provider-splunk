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

	"github.com/splunk/go-splunk-client/pkg/client"
)

// resourceIdHandlerTestCase defines a test case for the client.ID ResourceDataHandler's
// ResourceDataHandler.
type resourceIdHandlerTestCase struct {
	// name of the test, logged with failures
	name string

	// inputIdUrl is the string URL to parse to client.ID. If empty, an empty client.ID will be used
	// for the test.
	inputIdUrl string

	// wantValue is the expected value of the resource's Id after inputHelper's ResourceDataHandler is called
	wantValue string

	// wantError determines if an error is expected to be encountered
	wantError bool
}

// test performs the test defined by resourceIdHandlerTestCase.
func (test resourceIdHandlerTestCase) test(t *testing.T) {
	var id client.ID
	if test.inputIdUrl != "" {
		if parsedId, err := client.ParseID(test.inputIdUrl); err != nil {
			t.Fatalf("%s: ParseID returned error: %s", test.name, err)
		} else {
			id = parsedId
		}
	}

	handler := NewClientID(&id)
	d := resourceData(t)
	err := handler.ManageResourceFunc()(d)
	gotError := err != nil
	gotValue := d.Id()

	if gotError != test.wantError {
		t.Errorf("%s: ResourceDataHandler()(d) returned error? %v (%s)", test.name, gotError, err)
	}

	if gotValue != test.wantValue {
		t.Errorf("%s: ResourceDataHandler()(d) result:\n%s, want\n%s", test.name, gotValue, test.wantValue)
	}
}

// resourceIdHandlerTestCases is a slice of resourceIdHandlerTestCase.
type resourceIdHandlerTestCases []resourceIdHandlerTestCase

// test performs each test defined by resourceIdHandlerTestCases.
func (tests resourceIdHandlerTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}

// resourceIdHandlerTestCase defines a test case for the client.ID ResourceDataHandler's
// ObjectValueHandler.
type objectIdHandlerTestCase struct {
	// name of the test, logged with failures
	name string

	// inputId is the value stored as the resource's ID
	inputId string

	// wantValue is the expected value of the ID object, in URL form, after inputHelper's ObjectHandler is called
	wantValue string

	// wantError determines if an error is expected to be encountered
	wantError bool
}

// test performs the test defined by objectIdHandlerTestCase.
func (test objectIdHandlerTestCase) test(t *testing.T) {
	var wantId client.ID
	if test.wantValue != "" {
		if parsedId, err := client.ParseID(test.wantValue); err != nil {
			t.Fatalf("%s: ParseID returned error: %s", test.name, err)
		} else {
			wantId = parsedId
		}
	}

	var gotId client.ID
	handler := NewClientID(&gotId)
	d := resourceData(t)
	d.SetId(test.inputId)
	err := handler.ManageObjectFunc()(d)
	gotError := err != nil

	if gotError != test.wantError {
		t.Errorf("%s: ObjectValueHandler()(d) returned error? %v (%s)", test.name, gotError, err)
	}

	if gotId != wantId {
		t.Errorf("%s: ObjectValueHandler()(d) result:\n%#v, want\n%#v", test.name, gotId, wantId)
	}
}

// objectIdHandlerTestCases is a slice of objectIdHandlerTestCase.
type objectIdHandlerTestCases []objectIdHandlerTestCase

// test performs each test defined by objectIdHandlerTestCases.
func (tests objectIdHandlerTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}
