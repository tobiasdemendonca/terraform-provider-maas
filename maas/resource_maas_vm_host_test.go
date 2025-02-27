package maas_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"terraform-provider-maas/maas"
	"terraform-provider-maas/maas/testutils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMAASVMHost_DeployParams(t *testing.T) {
	// A VM host identifier. Used to create a VM, which is deployed as a VM host in this test.
	vmHostIdentifier := os.Getenv("TF_ACC_VM_HOST_ID")
	// A random string to be used for test
	rs := acctest.RandString(8)
	testMachineName := fmt.Sprintf("test-vm-host-machine-%s", rs)
	testVMHostName := fmt.Sprintf("test-vm-host-%s", rs)
	resourceName := fmt.Sprintf("maas_vm_host.%s", testVMHostName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, []string{"TF_ACC_VM_HOST_ID"}) },
		Providers:    testutils.TestAccProviders,
		CheckDestroy: testAccCheckMAASVMHostDestroy,
		ErrorCheck:   func(err error) error { return err },
		Steps: []resource.TestStep{
			{
				Config: testAccMaasVMHostDeployParams(vmHostIdentifier, testMachineName, testVMHostName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "lxd"),
					resource.TestCheckResourceAttr(resourceName, "deploy_params.0.enable_hw_sync", "true"),
					resource.TestCheckResourceAttr(resourceName, "deploy_params.0.user_data", "#!/bin/bash\necho 'Hello from cloud-init'"),
				),
			},
		},
	})
}

func testAccMaasVMHostDeployParams(vmHostIdentifier string, testMachineName string, testVMHostName string) string {
	return fmt.Sprintf(`
	resource "maas_vm_host_machine" "%s" {
	  vm_host = %q
	  cores   = 1
	  memory  = 2048

	  storage_disks {
	    size_gigabytes = 15
	  }
	}
	resource "maas_vm_host" "%s" {
	  machine = maas_vm_host_machine.%s.id
	  type    = "lxd"

	  deploy_params {
		  enable_hw_sync   = true
		  user_data        = "#!/bin/bash\necho 'Hello from cloud-init'"
	  }
	}
	`, testMachineName, vmHostIdentifier, testVMHostName, testMachineName)
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
			machine, err := client.Machine.Get(rs.Primary.Attributes["machine"])
			if err != nil {
				return fmt.Errorf("machine %s not found after VM host deletion, with error:\n%s", rs.Primary.Attributes["machine"], err)
			}
			if machine.StatusName != "Ready" {
				return fmt.Errorf("machine %s is not in Ready state after VM host deletion but in state %s", machine.SystemID, machine.StatusName)
			}
			continue
		}

		return err
	}

	return nil
}
