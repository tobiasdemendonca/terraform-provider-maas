package maas_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"terraform-provider-maas/maas"
	"terraform-provider-maas/maas/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Split the state ID of a tag in the format system_id:interface_id into its component ids, where system_id is the system ID of the machine or device, and interface_id is the ID of the network interface.
func SplitTagStateId(stateId string) (string, int, error) {
	splitId := strings.SplitN(stateId, ":", 2)
	if len(splitId) != 2 {
		return "", 0, fmt.Errorf("invalid resource ID: %s", stateId)
	}
	interfaceId, err := strconv.Atoi(splitId[1])
	if err != nil {
		return "", 0, err
	}
	return splitId[0], interfaceId, nil
}

func TestSplitTagStateId(t *testing.T) {
	expectedSystemId := "acb123"
	expectedInterfaceId := 12
	stateId := fmt.Sprintf("%s:%d", expectedSystemId, expectedInterfaceId)
	systemId, interfaceId, err := SplitTagStateId(stateId)
	if err != nil {
		t.Fatalf("Error splitting state ID: %s", err)
	}
	if systemId != expectedSystemId || interfaceId != expectedInterfaceId {
		t.Fatalf("Expected system ID %s and interface ID %d, got system ID %s and interface ID %d", expectedSystemId, expectedInterfaceId, systemId, interfaceId)
	}
}

func TestAccNetworkInterfaceTag(t *testing.T) {
	hostname := acctest.RandomWithPrefix("tf")
	macAddress := testutils.RandomMAC()
	tagName := acctest.RandomWithPrefix("tag")
	tagName2 := acctest.RandomWithPrefix("tag2")
	tagName3 := acctest.RandomWithPrefix("tag3")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		Providers:         testutils.TestAccProviders,
		ErrorCheck:        func(err error) error { return err },
		CheckDestroy:      testAccCheckMaasNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			// Test creation.
			{
				Config: testAccMaasNetworkInterfaceTagConfig(hostname, macAddress, tagName, tagName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMaasNetworkInterfaceTagExists("maas_network_interface_tag.test"),
					resource.TestCheckResourceAttr("maas_network_interface_tag.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("maas_network_interface_tag.test", "tags.0", tagName),
					resource.TestCheckResourceAttr("maas_network_interface_tag.test", "tags.1", tagName2),
				),
			},
			// Test update. Expected behaviour is that the previous tag is removed and the new tag is added.
			{
				Config: testAccMaasNetworkInterfaceTagConfig(hostname, macAddress, tagName2, tagName3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMaasNetworkInterfaceTagExists("maas_network_interface_tag.test"),
					resource.TestCheckResourceAttr("maas_network_interface_tag.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("maas_network_interface_tag.test", "tags.0", tagName2),
					resource.TestCheckResourceAttr("maas_network_interface_tag.test", "tags.1", tagName3),
				),
			},
			// // Test import.
			// {
			// 	ResourceName: "maas_network_interface_tag.test",
			// 	ImportState: true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccMaasNetworkInterfaceTagConfig(hostname, macAddress string, tagNames ...string) string {
	return fmt.Sprintf(`
resource "maas_device" "test" {
  hostname = %q
  network_interfaces {
	mac_address = %q 
  }
}

resource "maas_network_interface_tag" "test" {
  device = maas_device.test.hostname
  interface_id = [for iface in maas_device.test.network_interfaces : iface.id if iface.mac_address == %q][0]
  tags = %s
}
	`,hostname, macAddress, macAddress, fmt.Sprintf("[\"%s\"]", strings.Join(tagNames, "\", \"")))
} 

func testAccCheckMaasNetworkInterfaceTagExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		systemId, interfaceId, err := SplitTagStateId(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client
		response, err := conn.NetworkInterface.Get(systemId, interfaceId)
		if err != nil {
			return err
		}
		if response == nil {
			return fmt.Errorf("MAAS Network Interface (%s) not found.", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckMaasNetworkInterfaceDestroy(s *terraform.State) error {
	conn := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client

	// loop through the resources in state, verifying each maas_network_interface_tag is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "maas_network_interface_tag" {
			continue
		}

		// Retrieve the system and interface ID from the state ID
		systemId, interfaceId, err := SplitTagStateId(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Check the interface doesn't exist
		response, err := conn.NetworkInterface.Get(systemId, interfaceId)
		if err == nil {
			if response != nil && response.ID == interfaceId {
				return fmt.Errorf("MAAS Network Interface (%s) still exists.", rs.Primary.ID)
			}
		}
		// If the error is a 404, the interface is destroyed as expected
		if !strings.Contains(err.Error(), "404 Not Found") {
			return err
		}
	}
	return nil
}
