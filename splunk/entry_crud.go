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

package splunk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/go-splunk-client/pkg/client"
	"github.com/splunk/terraform-provider-splunk/internal/resource"
)

// useLegacyClient returns:
//
// * the resource's use_legacy_client value, if false or explicitly true
//
// * the provider's use_legacy_client_default value
func useLegacyClient(provider *SplunkProvider, d *schema.ResourceData) bool {
	resourceLegacyClientI, ok := d.GetOk("use_legacy_client")
	resourceLegacyClient := resourceLegacyClientI.(bool)
	// GetOk only returns true if the fetched value is not the zero value for its type,
	// so we can only determine if use_legacy_client was explicitly true. but because
	// true is our default value, we know that it can only be false if explicitly set.
	if ok || !resourceLegacyClient {
		return resourceLegacyClient
	}

	return provider.useLegacyClientDefault
}

// createFuncForHandler returns a schema.CreateFunc for the ResourceObjectHandler returned by
// the given function.
//
// During the transition between the legacy and external Splunk clients, it will return legacyFunction
// if the provider configuration or the resource configuration sets use_legacy_client=true.
func createFuncForHandler[T any](fn func(*T) resource.ResourceObjectManager, legacyFunction schema.CreateFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		provider := meta.(*SplunkProvider)

		if useLegacyClient(provider, d) {
			return legacyFunction(d, meta)
		}

		c := provider.ExternalClient

		entry := new(T)
		if err := fn(entry).ManageObjectFunc()(d); err != nil {
			return err
		}

		if err := c.Create(entry); err != nil {
			return err
		}

		return readFuncForHandler(fn, legacyFunction)(d, meta)
	}
}

// createFuncForHandlerFunc returns a schema.CreateFunc for the ResourceObjectHandler returned by
// the given function.
//
// During the transition between the legacy and external Splunk clients, it will return legacyFunction
// if the provider configuration or the resource configuration sets use_legacy_client=true.
func readFuncForHandler[T any](fn func(*T) resource.ResourceObjectManager, legacyFunction schema.CreateFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		provider := meta.(*SplunkProvider)

		if useLegacyClient(provider, d) {
			return legacyFunction(d, meta)
		}

		c := provider.ExternalClient

		entry := new(T)
		if err := fn(entry).ManageObjectFunc()(d); err != nil {
			return err
		}

		if err := c.Read(entry); err != nil {
			if clientErr, ok := err.(client.Error); ok {
				if clientErr.Code == client.ErrorNotFound {
					d.SetId("")

					return nil
				}
			}
			return err
		}

		if err := fn(entry).ManageResourceFunc()(d); err != nil {
			return err
		}

		return nil
	}
}

// updateFuncForHandler returns a schema.UpdateFunc for the ResourceObjectHandler returned by
// the given function.
//
// During the transition between the legacy and external Splunk clients, it will return legacyFunction
// if the provider configuration or the resource configuration sets use_legacy_client=true.
func updateFuncForHandler[T any](fn func(*T) resource.ResourceObjectManager, legacyFunction schema.CreateFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		provider := meta.(*SplunkProvider)

		if useLegacyClient(provider, d) {
			return legacyFunction(d, meta)
		}

		c := provider.ExternalClient

		entry := new(T)
		if err := fn(entry).ManageObjectFunc()(d); err != nil {
			return err
		}

		if err := c.Update(entry); err != nil {
			return err
		}

		return readFuncForHandler(fn, legacyFunction)(d, meta)
	}
}

// deleteFuncForHandler returns a schema.DeleteFunc for the ResourceObjectHandler returned by
// the given function.
//
// During the transition between the legacy and external Splunk clients, it will return legacyFunction
// if the provider configuration or the resource configuration sets use_legacy_client=true.
func deleteFuncForHandler[T any](fn func(*T) resource.ResourceObjectManager, legacyFunction schema.CreateFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		provider := meta.(*SplunkProvider)

		if useLegacyClient(provider, d) {
			return legacyFunction(d, meta)
		}

		c := provider.ExternalClient

		entry := new(T)
		if err := fn(entry).ManageObjectFunc()(d); err != nil {
			return err
		}

		if err := c.Delete(entry); err != nil {
			return err
		}

		return nil
	}
}
