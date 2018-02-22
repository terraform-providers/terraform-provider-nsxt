/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"strings"
	"testing"
)

func TestAccDataSourceNsxtSwitchingProfile_basic(t *testing.T) {
	profileName := getSwitchingProfileName()
	testResourceName := "data.nsxt_switching_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXSwitchingProfileReadTemplate(profileName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXSwitchingProfileExists(testResourceName, profileName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", profileName),
				),
			},
		},
	})
}

func testAccNSXSwitchingProfileExists(resourceName string, displayNamePrefix string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX switching profile data source %s not found", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX switching profile data source ID not set")
		}

		object, responseCode, err := nsxClient.LogicalSwitchingApi.GetSwitchingProfile(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving switching profile ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if switching profile %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if strings.HasPrefix(object.DisplayName, displayNamePrefix) {
			return nil
		}
		return fmt.Errorf("NSX switching profile data source '%s' wasn't found", displayNamePrefix)
	}
}

func testAccNSXSwitchingProfileReadTemplate(profileName string) string {
	return fmt.Sprintf(`
data "nsxt_switching_profile" "test" {
     display_name = "%s"
}`, profileName)
}
