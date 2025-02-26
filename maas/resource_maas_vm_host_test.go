package maas_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"terraform-provider-maas/maas"
	"terraform-provider-maas/maas/testutils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccMAASVMHost_DeployParams(t *testing.T) {
	// A VM host identifier. Used to create a VM, which is deployed as a VM host in this test.
	vmHostIdentifier := os.Getenv("TF_ACC_VM_HOST_MACHINE")
	// A random string to be used for test names
	rs := acctest.RandString(8)	

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.TestAccProviders,
		CheckDestroy: testAccCheckMAASVMHostDestroy,
		ErrorCheck:   func(err error) error { return err },
		Steps: []resource.TestStep{
			{
				Config: testAccMaasVMHostDeployParams(vmHostIdentifier, rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("maas_vm_host.test-vm-host", "type", "lxd"),
					resource.TestCheckResourceAttr("maas_vm_host.test-vm-host", "machine", vmHostIdentifier),
					resource.TestCheckResourceAttr("maas_vm_host.test-vm-host", "deploy_params.0.distro_series", "noble"),
					resource.TestCheckResourceAttr("maas_vm_host.test-vm-host", "deploy_params.0.enable_hw_sync", "true"),
				),
			},
		},
	})
}

func testAccMaasVMHostDeployParams(vmHostIdentifier string, rs string) string {
	return fmt.Sprintf(`
	resource "maas_vm_host_machine" "test-vm-host-machine-%s" {
	  vm_host = %q
	  cores   = 1
	  memory  = 2048

	  storage_disks {
	    size_gigabytes = 15
	  }
	}
	resource "maas_vm_host" "test-vm-host-%s" {
	  machine = maas_vm_host_machine.test-vm-host-machine-%s.id
	  type    = "lxd"

	  deploy_params {
		  distro_series    = "ubuntu/noble"
		  enable_hw_sync   = true
		  user_data        = "#!/bin/bash\necho 'Hello from cloud-init'"
	  }
	}
	`, rs, vmHostIdentifier, rs, rs)
}

func testAccCheckMAASVMHostDestroy(s *terraform.State) error {
	client := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "maas_vm_host" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		response, err := client.VMHost.Get(id)
		if err == nil {
			if response != nil && response.ID == id {
				return fmt.Errorf("VM host still exists")
			}
		}

		// If the error is a 404 not found error, the VM host is destroyed
		if err != nil && strings.Contains(err.Error(), "404 Not Found") {
			continue
		}

		return err
	}

	return nil
}
