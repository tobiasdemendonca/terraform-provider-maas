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

// {
// 	"active_discovery_interval", "10000"
// 	"auto_vlan_creation", "false"
// 	"boot_images_auto_import", "false"
// 	"boot_images_no_proxy", "true"
// 	"commissioning_distro_series", "jammy"
// 	"completed_intro", "false"
// 	"curtin_verbose", "false"
// 	"default_boot_interface_link_type",
// 	"default_distro_series",
// 	"default_dns_ttl",
// 	"default_min_hwe_kernel",
// 	"default_osystem",
// 	"default_storage_layout",
// 	"disk_erase_with_quick_erase",
// 	"disk_erase_with_secure_erase",
// 	"dns_trusted_acl",
// 	"dnssec_validation",
// 	"enable_analytics",
// 	"enable_disk_erasing_on_release",
// 	"enable_http_proxy",
// 	"enable_kernel_crash_dump",
// 	"enable_third_party_drivers",
// 	"enlist_commissioning",
// 	"force_v1_network_yaml",
// 	"hardware_sync_interval",
// 	"http_proxy",
// 	"kernel_opts",
// 	"maas_auto_ipmi_cipher_suite_id",
// 	"maas_auto_ipmi_k_g_bmc_key",
// 	"maas_auto_ipmi_user",
// 	"maas_auto_ipmi_user_privilege_level",
// 	"maas_auto_ipmi_workaround_flags",
// 	"maas_internal_domain",
// 	"maas_name",
// 	"maas_proxy_port",
// 	"maas_syslog_port",
// 	"max_node_commissioning_results",
// 	"max_node_installation_results",
// 	"max_node_release_results",
// 	"max_node_testing_results",
// 	"network_discovery",
// 	"node_timeout",
// 	"ntp_external_only",
// 	"ntp_servers",
// 	"prefer_v4_proxy",
// 	"prometheus_enabled",
// 	"prometheus_push_gateway",
// 	"prometheus_push_interval",
// 	"promtail_enabled",
// 	"promtail_port",
// 	"release_notifications",
// 	"remote_syslog",
// 	"session_length",
// 	"subnet_ip_exhaustion_threshold_count",
// 	"theme",
// 	"tls_cert_expiration_notification_enabled",
// 	"tls_cert_expiration_notification_interval",
// 	"upstream_dns",
// 	"use_peer_proxy",
// 	"use_rack_proxy",
// 	"vcenter_datacenter",
// 	"vcenter_password",
// 	"vcenter_server",
// 	"vcenter_username",
// 	"windows_kms_host",
// }

func getOriginalValue(key string) (string, error) {
	client := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client

	val, err := client.MAASServer.Get(key)
	if err != nil {
		return "", fmt.Errorf("error getting original value for key %s: %v", key, err)
	}

	return strings.Trim(string(val),"\""), nil

}

{
	key: "default_boot_interface_link_type",
	value: ,
},
{
	key: "default_distro_series",
	value: ,
},
{
	key: "default_dns_ttl",
	value: ,
},
{
	key: "default_min_hwe_kernel",
	value: ,
},
{
	key: "default_osystem",
	value: ,
},
{
	key: "default_storage_layout",
	value: ,
},
{
	key: "disk_erase_with_quick_erase",
	value: ,
},
{
	key: "disk_erase_with_secure_erase",
	value: ,
},
{
	key: "dns_trusted_acl",
	value: ,
},
{
	key: "dnssec_validation",
	value: ,
},
{
	key: "enable_analytics",
	value: ,
},
{
	key: "enable_disk_erasing_on_release",
	value: ,
},
{
	key: "enable_http_proxy",
	value: ,
},
{
	key: "enable_kernel_crash_dump",
	value: ,
},
{
	key: "enable_third_party_drivers",
	value: ,
},
{
	key: "enlist_commissioning",
	value: ,
},
{
	key: "force_v1_network_yaml",
	value: ,
},
{
	key: "hardware_sync_interval",
	value: ,
},
{
	key: "http_proxy",
	value: ,
},
{
	key: "kernel_opts",
	value: ,
},
{
	key: "maas_auto_ipmi_cipher_suite_id",
	value: ,
},
{
	key: "maas_auto_ipmi_k_g_bmc_key",
	value: ,
},
{
	key: "maas_auto_ipmi_user",
	value: ,
},
{
	key: "maas_auto_ipmi_user_privilege_level",
	value: ,
},
{
	key: "maas_auto_ipmi_workaround_flags",
	value: ,
},
{
	key: "maas_internal_domain",
	value: ,
},
{
	key: "maas_name",
	value: ,
},
{
	key: "maas_proxy_port",
	value: ,
},
{
	key: "maas_syslog_port",
	value: ,
},
{
	key: "max_node_commissioning_results",
	value: ,
},
{
	key: "max_node_installation_results",
	value: ,
},
{
	key: "max_node_release_results",
	value: ,
},
{
	key: "max_node_testing_results",
	value: ,
},
{
	key: "network_discovery",
	value: ,
},
{
	key: "node_timeout",
	value: ,
},
{
	key: "ntp_external_only",
	value: ,
},
{
	key: "ntp_servers",
	value: ,
},
{
	key: "prefer_v4_proxy",
	value: ,
},
{
	key: "prometheus_enabled",
	value: ,
},
{
	key: "prometheus_push_gateway",
	value: ,
},
{
	key: "prometheus_push_interval",
	value: ,
},
{
	key: "promtail_enabled",
	value: ,
},
{
	key: "promtail_port",
	value: ,
},
{
	key: "release_notifications",
	value: ,
},
{
	key: "remote_syslog",
	value: ,
},
{
	key: "session_length",
	value: ,
},
{
	key: "subnet_ip_exhaustion_threshold_count",
	value: ,
},
{
	key: "theme",
	value: ,
},
{
	key: "tls_cert_expiration_notification_enabled",
	value: ,
},
{
	key: "tls_cert_expiration_notification_interval",
	value: ,
},
{
	key: "upstream_dns",
	value: ,
},
{
	key: "use_peer_proxy",
	value: ,
},
{
	key: "use_rack_proxy",
	value: ,
},
{
	key: "vcenter_datacenter",
	value: ,
},
{
	key: "vcenter_password",
	value: ,
},
{
	key: "vcenter_server",
	value: ,
},
{
	key: "vcenter_username",
	value: ,
},
{
	key: "windows_kms_host",
	value: ,
},

func TestAccResourceMAASConfiguration_basic(t * testing.T) {

	testCases := []struct{
		key string
		value string
		originalValue string
	} {
		{
			key: "active_discovery_interval",
			value: "10800",
		},
		{
			key:"auto_vlan_creation",
			value: "false",
		},
		{
			key:"boot_images_auto_import",
			value: "false",
		},
		{
			key:"boot_images_no_proxy",
			value: "true",
		},
		{
			key:"commissioning_distro_series",
			value: "jammy",
		},
		{
			key:"completed_intro",
			value: "false",
		},
		{
			key:"curtin_verbose",
			value: "false",
		},
	}
	log.Print("Test cases for MAAS Configuration:")
	for _, testCase := range testCases {
		// val, err := getOriginalValue(testCase.key)
		// if err != nil {
		// 	t.Fatalf("error getting original value for key %s: %v", testCase.key, err)
		// }
		// testCase.originalValue = val

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testutils.PreCheck(t, nil) },
			Providers:   testutils.TestAccProviders,
			CheckDestroy: testAccMAASConfigurationCheckDestroy,
			ErrorCheck:   func(err error) error {return err},
			Steps: []resource.TestStep{
				{	
					Config: testAccMAASConfigurationConfigBasic(testCase.key, testCase.value),
					Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("maas_configuration.test", "value", testCase.value),
							testAccMAASConfigurationCheckExists("maas_configuration.test"),
							resource.TestCheckResourceAttr("maas_configuration.test", "key", testCase.key),
					),
				},
			},
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