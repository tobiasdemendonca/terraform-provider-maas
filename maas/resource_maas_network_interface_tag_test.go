package maas_test

import (
	"fmt"
	"testing"

	"terraform-provider-maas/maas/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkInterfaceTag(t *testing.T) {
	hostname := acctest.RandomWithPrefix("tf")
	macAddress := testutils.RandomMAC()
	tagName := acctest.RandomWithPrefix("tag")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		Providers:         testutils.TestAccProviders,
		ErrorCheck:        func(err error) error { return err },
		CheckDestroy:      func(s *terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config: testAccMaasNetworkInterfaceTagConfig(hostname, macAddress, tagName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("maas_network_interface_tag.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("maas_network_interface_tag.test", "tags.0", tagName),
				),
			},
		},
	})
}

func testAccMaasNetworkInterfaceTagConfig(hostname, macAddress, tagName string) string {
	return fmt.Sprintf(`
resource "maas_device" "test" {
  hostname = %q
  network_interfaces {
    mac_address = %q 
  }
}

resource "maas_network_interface_tag" "test" {
  device = maas_device.test.id
  interface_id = [for iface in maas_device.test.network_interfaces : iface.id if iface.mac_address == %q][0]
  tags = [%q]
}
	`,hostname, macAddress, macAddress, tagName)
}
