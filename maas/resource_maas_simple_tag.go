package maas

import (
	"context"
	"fmt"
	"slices"

	// "github.com/canonical/gomaasclient/client"
	// "github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMaasInterfaceTag() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource to manage a MAAS tag as a string.",
		CreateContext: resourceInterfaceTagCreate,
		ReadContext:   resourceInterfaceTagRead,
		UpdateContext: resourceInterfaceTagUpdate,
		DeleteContext: resourceInterfaceTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil // TODO: implement
			},
		},
		Schema: map[string]*schema.Schema{
			"device": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"machine", "device"},
				Description:  "The identifier (system ID, hostname, or FQDN) of the device with the network interface. Either `machine` or `device` must be provided.",
			},
			"machine": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"machine", "device"},
				Description:  "The identifier (system ID, hostname, or FQDN) of the machine with the network interface. Either `machine` or `device` must be provided.",
			},
			"interface_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The network interface ID to tag.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The tags to assign to the network interface. It should be short and without spaces.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceInterfaceTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client

	interfaceId := d.Get("interface_id ").(int)
	desiredTags := convertToStringSlice(d.Get("tags").(*schema.Set).List())
	systemId, err := getMachineOrDeviceSystemID(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get the existing interface
	existingInterface, err := client.NetworkInterface.Get(systemId, interfaceId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Remove tags that are not in the desired set
	if existingInterface.Tags == nil {
		existingInterface.Tags = []string{}
	}
	for _, tag := range existingInterface.Tags {
		if !slices.Contains(desiredTags, tag) {
			client.NetworkInterface.RemoveTag(systemId, interfaceId, tag)
		} 
	}

	// Add tags that are in the desired set. AddTag will not add duplicates.
	for _, tag := range desiredTags {
		client.NetworkInterface.AddTag(systemId, interfaceId, tag)
	}

	// Create the resource ID in state. A unique resource for every interface.
	d.SetId(fmt.Sprintf("%v:%v", systemId, interfaceId))

	// Read the resource to update state
	return resourceInterfaceTagRead(ctx, d, meta)
}

func resourceInterfaceTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	
	return nil
}

func resourceInterfaceTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Add tags if not present
	// Since name is ForceNew, there's nothing to update in our dummy implementation
	return nil
}

func resourceInterfaceTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	
	// In a real implementation, we would delete from MAAS
	// For testing, we just let Terraform remove it from state
	d.SetId("")
	return nil
}


