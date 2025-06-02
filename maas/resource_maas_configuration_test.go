package maas_test

import (
	"fmt"
	"log"
	"strings"
	"terraform-provider-maas/maas/testutils"
	"testing"

	"terraform-provider-maas/maas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	// "google.golang.org/grpc/keepalive"
	// "github.com/canonical/gomaasclient/client"
)


func getOriginalValue(key string) (string, error) {
	client := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client

	val, err := client.MAASServer.Get(key)
	if err != nil {
		return "", fmt.Errorf("error getting original value for key %s: %v", key, err)
	}

	return strings.Trim(string(val),"\""), nil

}



func TestAccResourceMAASConfiguration_basic(t * testing.T) {
	// Test all MAAS Settings cases. Set twice to ensure a change is actually made
	// incase defaults change, and value 2 is selected to be a typical value to 
	// ensure dev environments aren't greatly changed between test runs.
	testCases := []struct{
		key string
		value_1 string
		value_2 string
	} {
		{
			key: "active_discovery_interval",
			value_1: "21600",
			value_2: "10800",
		},
		{
			key:"auto_vlan_creation",
			value_1: "false",
			value_2: "true",
		},
		{
			key:"boot_images_auto_import",
			value_1: "false",
			value_2: "true",
		},
		{
			key:"boot_images_no_proxy",
			value_1: "true",
			value_2: "false",
		},
		{
			key:"commissioning_distro_series",
			value_1: "noble",
			value_2: "noble",    // Dev environment does not necessarily have multiple distros
		},
		{
			key:"completed_intro",
			value_1: "false",
			value_2: "true",
		},
		{
			key:"curtin_verbose",
			value_1: "false",
			value_2: "true",
		},
		{
			key: "default_boot_interface_link_type",
			value_1:  "static",
			value_2:  "auto", 
		},
		{
			key: "default_distro_series",
			value_1:  "noble", 
			value_2:  "noble",    // Dev environment does not necessarily have multiple distros
		},
		{
			key: "default_dns_ttl",
			value_1:  "60", 
			value_2:  "30",
		},
		{
			key: "default_min_hwe_kernel",
			value_1:  "ga-24.04", 
			value_2:  "",
		},
		{
			key: "default_osystem",
			value_1:  "",
			value_2:  "ubuntu",
		},
		{
			key: "default_storage_layout",
			value_1:  "bcache", 
			value_2:  "blank",
		},
		{
			key: "disk_erase_with_quick_erase",
			value_1:  "true",
			value_2:  "false",
		},
		{
			key: "disk_erase_with_secure_erase",
			value_1:  "true",
			value_2:  "false",
		},
		{
			key: "dns_trusted_acl",
			value_1:  "192.168.1.1 192.168.1.2",
			value_2:  "",
		},
		{
			key: "dnssec_validation",
			value_1:  "yes",
			value_2:  "auto", 
		},
		{
			key: "enable_analytics",
			value_1:  "false",
			value_2: "true",
		},
		{
			key: "enable_disk_erasing_on_release",
			value_1:  "true", 
			value_2:  "false",
		},
		{
			key: "enable_http_proxy",
			value_1:  "false",
			value_2:  "true", 
		},
		{
			key: "enable_kernel_crash_dump",
			value_1:  "true", 
			value_2:  "false",
		},
		{
			key: "enable_third_party_drivers",
			value_1:  "false", 
			value_2:  "true",
		},
		{
			key: "enlist_commissioning",
			value_1:  "false", 
			value_2:  "true",
		},
		{
			key: "force_v1_network_yaml",
			value_1:  "true",
			value_2:  "false", 
		},
		{
			key: "hardware_sync_interval",
			value_1:  "10m", 
			value_2:  "15m",
		},
		{
			key: "http_proxy",
			value_1:  "http://proxy.example.com:8080", 
			value_2:  "",
		},
		{
			key: "kernel_opts",
			value_1:  "console=ttyS0", 
			value_2:  "",
		},
		{
			key: "maas_auto_ipmi_cipher_suite_id",
			value_1:  "8", 
			value_2:  "3",
		},
		{
			key: "maas_auto_ipmi_k_g_bmc_key",
			value_1:  "12345678901234567890", // Must be 20 characters
			value_2:  "",
		},
		{
			key: "maas_auto_ipmi_user",
			value_1:  "admin", 
			value_2:  "maas",
		},
		{
			key: "maas_auto_ipmi_user_privilege_level",
			value_1:  "USER", 
			value_2:  "ADMIN",
		},
		// {
		// 	key: "maas_auto_ipmi_workaround_flags",   // Does not currently work 
		// 	value_1:  , 
		// },
		{
			key: "maas_internal_domain",
			value_1:  "maas-internal-alt", 
			value_2:  "maas-internal",
		},
		{
			key: "maas_name",
			value_1:  "maas-alt", 
			value_2:  "maas-dev",
		},
		{
			key: "maas_proxy_port",
			value_1:  "8080", 
			value_2:  "8000",
		},
		{
			key: "maas_syslog_port",
			value_1:  "49152",
			value_2:  "5247",
		},
		{
			key: "max_node_commissioning_results",
			value_1:  "100",
			value_2:  "10",
		},
		{
			key: "max_node_installation_results",
			value_1:  "5",
			value_2:  "3",
		},
		{
			key: "max_node_release_results",
			value_1:  "5",
			value_2:  "3",
		},
		{
			key: "max_node_testing_results",
			value_1:  "100",
			value_2:  "10",
		},
		{
			key: "network_discovery",
			value_1:  "disabled", 
			value_2:  "enabled",
		},
		{
			key: "node_timeout",
			value_1:  "10",
			value_2:  "30",
		},
		{
			key: "ntp_external_only",
			value_1:  "true",
			value_2: "false",
		},
		{
			key: "ntp_servers",
			value_1:  "ntp.example.com",
			value_2: "ntp.ubuntu.com"  , 
		},
		{
			key: "prefer_v4_proxy",
			value_1:  "true",
			value_2: "false",
		},
		{
			key: "prometheus_enabled",
			value_1:  "true",
			value_2: "false",
		},
		{
			key: "prometheus_push_gateway",
			value_1:  "http://prometheus.example.com:9090",
			value_2:  "",
		},
		{
			key: "prometheus_push_interval",
			value_1:  "80", 
			value_2:  "60",
		},
		{
			key: "promtail_enabled",
			value_1:  "true", 
			value_2:  "false",
		},
		{
			key: "promtail_port",
			value_1:  "49153", 
			value_2:  "5238",
		},
		{
			key: "release_notifications",
			value_1:  "false", 
			value_2:  "true",
		},
		{
			key: "remote_syslog",
			value_1:  "example.com:514", 
			value_2:  "",
		},
		{
			key: "session_length",
			value_1:  "604800",  // 7 days
			value_2:  "1209600", // 14 days
		},
		{
			key: "subnet_ip_exhaustion_threshold_count",
			value_1:  "24", 
			value_2:  "16",
		},
		{
			key: "theme",
			value_1:  "sage", 
			value_2:  "",
		},
		{
			key: "tls_cert_expiration_notification_enabled",
			value_1:  "true", 
			value_2:  "false",
		},
		{
			key: "tls_cert_expiration_notification_interval",
			value_1:  "60", 
			value_2:  "30",
		},
		{
			key: "upstream_dns",
			value_1:  "8.8.8.8", 
			value_2:  "",
		},
		{
			key: "use_peer_proxy",
			value_1:  "true", 
			value_2:  "false",
		},
		{
			key: "use_rack_proxy",
			value_1:  "false", 
			value_2:  "true",
		},
		{
			key: "vcenter_datacenter",
			value_1:  "maas-vcenter", 
			value_2:  "",
		},
		{
			key: "vcenter_password",
			value_1:  "dummy-password", 
			value_2:  "",
		},
		{
			key: "vcenter_server",
			value_1:  "vcenter.example.com", 
			value_2:  "",
		},
		{
			key: "vcenter_username",
			value_1:  "maas", 
			value_2:  "",
		},
		{
			key: "windows_kms_host",
			value_1:  "kms.example.com", 
			value_2:  "",
		},
		
	}
	log.Print("Test cases for MAAS Configuration:")
	for _, testCase := range testCases {
		// val, err := getOriginalValue(testCase.key)
		// if err != nil {
		// 	t.Fatalf("error getting original value for key %s: %v", testCase.key, err)
		// }
		// testCase.originalValue = val

		t.Run(testCase.key, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:     func() { testutils.PreCheck(t, nil) },
				Providers:   testutils.TestAccProviders,
				CheckDestroy: testAccMAASConfigurationCheckDestroy,
				ErrorCheck:   func(err error) error {return err},
				Steps: []resource.TestStep{
					{	
						Config: testAccMAASConfigurationConfigBasic(testCase.key, testCase.value_1),
						Check: resource.ComposeTestCheckFunc(
							testAccMAASConfigurationCheckExists("maas_configuration.test"),
							resource.TestCheckResourceAttr("maas_configuration.test", "key", testCase.key),
							resource.TestCheckResourceAttr("maas_configuration.test", "value", testCase.value_1),
						),
					},
					{
						Config: testAccMAASConfigurationConfigBasic(testCase.key, testCase.value_2),
						Check: resource.ComposeTestCheckFunc(
							testAccMAASConfigurationCheckExists("maas_configuration.test"),
							resource.TestCheckResourceAttr("maas_configuration.test", "key", testCase.key),
							resource.TestCheckResourceAttr("maas_configuration.test", "value", testCase.value_2),
						),
					},
				},
			})
		})
	}
}

func testAccMAASConfigurationCheckExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if it exists in state
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s\n %v", rn, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		conn := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client

		value, err := conn.MAASServer.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting configuration value: %s", err)
		}
		cleanedVal := strings.Trim(string(value), "\"")
		if cleanedVal != rs.Primary.Attributes["value"] {
			return fmt.Errorf("configuration value does not match: expected %s, got %s", rs.Primary.Attributes["value"], value)
		}
		return nil
	}
}

func testAccMAASConfigurationCheckDestroy(s *terraform.State) error {
	_ = testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client
	return nil
}

func testAccMAASConfigurationConfigBasic(key string, value string) string {
	return fmt.Sprintf(`
resource "maas_configuration" "test" {
  key   = %q
  value = %q
}
`, key, value)
}