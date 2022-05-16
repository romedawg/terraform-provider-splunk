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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/go-splunk-client/pkg/client"
)

// clientID implements the ResourceObjectHandler interface for client.ID.
type clientID struct {
	id *client.ID
}

// NewClientID returns a ResourceObjectHandler that maps a client.ID to a Terraform
// resource ID.
func NewClientID(id *client.ID) ResourceObjectManager {
	return clientID{
		id: id,
	}
}

// ManageResourceFunc returns a ResourceDataFunc that sets the schema.ResourceData
// Id to the stored client.ID's URL, if client.ID was previously set by URL.
func (id clientID) ManageResourceFunc() ResourceDataFunc {
	return func(d *schema.ResourceData) error {
		if idURL, err := id.id.URL(); err == nil {
			d.SetId(idURL)
		}

		return nil
	}
}

// ManageObjectFunc returns a ResourceDataFunc that sets client.ID by the URL ID stored in
// schema.ResourceData.
func (id clientID) ManageObjectFunc() ResourceDataFunc {
	return func(d *schema.ResourceData) error {
		if d.Id() != "" {
			// an unparseable ID URL should be ignored, so it can instead be determined at read-time.
			// this is a likely scenario when moving from the legacy client to the external client.
			if parsedId, err := client.ParseID(d.Id()); err == nil {
				*id.id = parsedId
			}
		}

		return nil
	}
}
