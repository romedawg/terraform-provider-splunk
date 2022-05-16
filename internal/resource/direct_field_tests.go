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
	"reflect"
	"testing"
)

// directFieldTestStorage is a simple struct that stores a value for
// a generic type. It exists to enable testing that is otherwise difficult
// due to needing pointers to values.
type directFieldTestStorage[T any] struct {
	value T
}

// newDirectFieldTestStorage returns a new directFieldTestStruct for the given
// value.
func newDirectFieldTestStorage[T any](value T) *directFieldTestStorage[T] {
	return &directFieldTestStorage[T]{
		value: value,
	}
}

// resourceObjectHandler returns a ResourceObjectHandler for the directFieldTestStorage.
func (test *directFieldTestStorage[T]) resourceObjectHandler(key string) ResourceObjectManager {
	return NewDirectField(key, &test.value)
}

// equals returns a boolean indicating if the stored value is deeply equal to the given value.
func (test directFieldTestStorage[T]) equals(value interface{}) bool {
	return reflect.DeepEqual(test.value, value)
}

// storedValue returns the stored value.
func (test directFieldTestStorage[T]) storedValue() interface{} {
	return test.value
}

// directListFieldTestStorage is a simple struct that stores a slice of values
// for a generic type. It exists to enable testing, similar to directFieldTestStorage.
type directListFieldTestStorage[T any] struct {
	value []T
}

// newDirectListFieldTestStorage returns a new directListFieldTestStorage for the given
// slice of values.
func newDirectListFieldTestStorage[T any](value []T) *directListFieldTestStorage[T] {
	return &directListFieldTestStorage[T]{
		value: value,
	}
}

// resourceObjectHandler returns a ResourceObjectHandler for the directListFieldTestStorage.
func (test *directListFieldTestStorage[T]) resourceObjectHandler(key string) ResourceObjectManager {
	return NewDirectListField(key, &test.value)
}

// equals returns a boolean indicating if the stored value is deeply equal to the given value.
func (test directListFieldTestStorage[T]) equals(value interface{}) bool {
	return reflect.DeepEqual(test.value, value)
}

// storedValue returns the stored value.
func (test directListFieldTestStorage[T]) storedValue() interface{} {
	return test.value
}

// resourceDataTestHelper is an interface to simplify testing of directFieldTestStorage.
type resourceDataTestHelper interface {
	resourceObjectHandler(string) ResourceObjectManager
	equals(interface{}) bool
	storedValue() interface{}
}

// resourceDataHandlerTestCase defines a test case for ResourceDataHandlers.
type resourceDataHandlerTestCase struct {
	// name of the test, logged with failures
	name string

	// resourceDataTestHelper stores the object value that corresponds to a resource data field
	inputStruct resourceDataTestHelper

	// inputKey is the name of the resource data field the stored value corresponds to
	inputKey string

	// wantValue is the expected value in the resource data after inputHelper's ResourceDataHandler is called
	wantValue interface{}

	// wantError determines if an error is expected to be encountered
	wantError bool
}

// test performs the test defined by resourceDataHandlerTestCase.
func (test resourceDataHandlerTestCase) test(t *testing.T) {
	handler := test.inputStruct.resourceObjectHandler(test.inputKey)
	d := resourceData(t)
	err := handler.ManageResourceFunc()(d)
	gotError := err != nil
	gotValue := d.Get(test.inputKey)

	if gotError != test.wantError {
		t.Errorf("%s: ResourceDataHandler()(d) returned error? %v (%s)", test.name, gotError, err)
	}

	if !reflect.DeepEqual(gotValue, test.wantValue) {
		t.Errorf("%s: ResourceDataHandler()(d) result:\n%#v, want\n%#v", test.name, gotValue, test.wantValue)
	}
}

// resourceDataHandlerTestCases is a slice of resourceDataHandlerTestCase.
type resourceDataHandlerTestCases []resourceDataHandlerTestCase

// test performs the tests defined by resourceDataHandlerTestCases.
func (tests resourceDataHandlerTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}

// resourceDataObjectHandlerTestCase defines a test case for ObjectHandlers.
type resourceDataObjectHandlerTestCase struct {
	// name of the test, logged with failures
	name string

	// resourceDataTestHelper stores the object value that corresponds to a resource data field
	inputHelper resourceDataTestHelper

	// inputKey is the name of the resource data field the stored value corresponds to
	inputKey string

	// wantValue is the expected value in the inputHelper after inputHelper's ObjectHandler is called
	wantValue interface{}

	// wantError determines if an error is expected to be encountered
	wantError bool
}

// test performs the test defined by resourceDataObjectHandlerTestCase.
func (test resourceDataObjectHandlerTestCase) test(t *testing.T) {
	handler := test.inputHelper.resourceObjectHandler(test.inputKey)
	d := resourceData(t)
	err := handler.ManageObjectFunc()(d)
	gotError := err != nil

	if gotError != test.wantError {
		t.Errorf("%s: ObjectHandler()(d) returned error? %v (%s)", test.name, gotError, err)
	}

	if !test.inputHelper.equals(test.wantValue) {
		t.Errorf("%s: ObjectHandler()(d) result:\n%#v, want\n%#v", test.name, test.inputHelper.storedValue(), test.wantValue)
	}
}

// resourceDataObjectHandlerTestCases is a slice of resourceDataObjectHandlerTestCase.
type resourceDataObjectHandlerTestCases []resourceDataObjectHandlerTestCase

// test performs the tests defined by resourceDataObjectHandlerTestCases.
func (tests resourceDataObjectHandlerTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}
