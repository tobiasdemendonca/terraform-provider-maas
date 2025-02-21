package maas_test

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-maas/maas"
	"terraform-provider-maas/maas/testutils"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceMaasSubnetIPRange(t *testing.T) {
	// Setup ip range attrs
	var ipRange entity.IPRange
	range_type := "reserved"
	comment := "test-comment"
	start_ip := "10.88.88.1"
	end_ip := "10.88.88.50"

	pre_check := func() { testutils.PreCheck(t, nil) }

	// Test create
	checks := []resource.TestCheckFunc{
		testAccMAASSubnetIPRangeCheckExists("maas_subnet_ip_range.test_ip_range", &ipRange),
		resource.TestCheckResourceAttr("maas_subnet_ip_range.test_ip_range", "type", range_type),
		resource.TestCheckResourceAttr("maas_subnet_ip_range.test_ip_range", "comment", comment),
		resource.TestCheckResourceAttr("maas_subnet_ip_range.test_ip_range", "start_ip", start_ip),
		resource.TestCheckResourceAttr("maas_subnet_ip_range.test_ip_range", "end_ip", end_ip),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: pre_check,
		Providers: testutils.TestAccProviders,
		CheckDestroy: testAccCheckMAASSubnetIPRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMAASSubnetIPRange(range_type, comment, start_ip, end_ip),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

func testAccMAASSubnetIPRangeCheckExists(rn string, ipRange *entity.IPRange) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s\n %#v", rn, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		conn := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		gotIPRange, err := conn.IPRange.Get(id)
		if err != nil {
			return fmt.Errorf("error getting ip range: %s", err)
		}

		*ipRange = *gotIPRange

		return nil
	}
}

func testAccMAASSubnetIPRange(
	range_type string,
	comment string,
	start_ip string,
	end_ip string,
) string {
	// A subnet is required to create an IP range
	return fmt.Sprintf(`
		resource "maas_subnet" "test_subnet" { 
			cidr = "10.88.88.0/26" 
			name = "test-tf-subnet"
			gateway_ip = "10.88.88.1"
			dns_servers = ["8.8.8.8"]
		}

		resource "maas_subnet_ip_range" "test_ip_range" {
			subnet = maas_subnet.test_subnet.id
			type = "%s"
			start_ip = "%s"
			end_ip = "%s"
			comment = "%s"
		}
	`, range_type, start_ip, end_ip, comment)
}

func testAccCheckMAASSubnetIPRangeDestroy(s *terraform.State) error {
	// retrieve the connection established in Provider configuration
	conn := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client

	// loop through the resources in state, verifying each maas_subnet_ip_range
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "maas_subnet_ip_range" && rs.Type != "maas_subnet" {
			continue
		}
		
		// retrieve the maas_subnet_ip_range by referencing it's state ID for API lookup
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		var exists bool
		if rs.Type == "maas_subnet_ip_range" {
			response, err := conn.IPRange.Get(id)
			if err == nil {
				if response != nil && response.ID == id {
					exists = true
				}
			}
		} else {
			response, err := conn.Subnet.Get(id)
			if err == nil {
				if response != nil && response.ID == id {
					exists = true
				}
			}
		}

		if exists {
			return fmt.Errorf("MAAS %s (%s) still exists.", rs.Type, rs.Primary.ID)
		}

		// If the error is equivalent to 404 not found, the maas_resource_pool is destroyed.
		// Otherwise return the error
		if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
			return err
		}
	}

	return nil
}
