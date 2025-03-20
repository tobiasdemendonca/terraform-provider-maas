package maas_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"terraform-provider-maas/maas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"terraform-provider-maas/maas/testutils"
)

func TestAccBlockDeviceTag_basic(t *testing.T) {
	machine := os.Getenv("TF_ACC_BLOCK_DEVICE_MACHINE")

	blockDeviceName := acctest.RandomWithPrefix("tf")
	tagName := acctest.RandomWithPrefix("tag")
	tagName2 := acctest.RandomWithPrefix("tag2")
	tagName3 := acctest.RandomWithPrefix("tag3")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.TestAccProviders,
		ErrorCheck:   func(err error) error { return err },
		CheckDestroy: testAccCheckMaasBlockDeviceTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockDeviceTagConfig(machine, blockDeviceName, tagName, tagName2, tagName3),
			},
		},
	})
}

func testAccBlockDeviceTagConfig(hostname string, name string, tagNames ...string) string {
	return fmt.Sprintf(`

data "maas_machine" "machine" {
  hostname = %q
}

resource "maas_block_device" "test" {
  machine        = data.maas_machine.machine.id
  name           = %q
  size_gigabytes = 1
  id_path        = "/dev/test"
}

resource "maas_block_device_tag" "test" {
  block_device_id = maas_block_device.test.id
  machine = maas_block_device.test.machine
  tags = %s
  }
	`, hostname, name, fmt.Sprintf("[\"%s\"]", strings.Join(tagNames, "\", \"")))
}


func testAccCheckMaasBlockDeviceTagDestroy(s *terraform.State) error {
	conn := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client

	// loop through the resources in state, verifying each maas_block_device_tag is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "maas_block_device_tag" {
			continue
		}

		// Retrieve the system and block device ID from the state ID
		systemId, blockDeviceId, err := maas.SplitTagStateId(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Check the block device doesn't exist
		response, err := conn.BlockDevice.Get(systemId, blockDeviceId)
		if err == nil {
			if response != nil && response.ID == blockDeviceId {
				return fmt.Errorf("MAAS Block Device (%s) still exists.", rs.Primary.ID)
			}
		}
		// If the error is a 404, the block device is destroyed as expected
		if !strings.Contains(err.Error(), "404 Not Found") {
			return err
		}
	}
	return nil
}