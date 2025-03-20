package maas

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMaasNetworkInterfaceTag() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource to manage a MAAS tags as strings on a network interface not managed by Terraform.",
		CreateContext: resourceNetworkInterfaceTagCreate,
		ReadContext:   resourceNetworkInterfaceTagRead,
		UpdateContext: resourceNetworkInterfaceTagUpdate,
		DeleteContext: resourceNetworkInterfaceTagDelete,
		// Importer: &schema.ResourceImporter{
		// 	StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		// 		client := meta.(*ClientConfig).Client

		// 		systemId, err := getMachineOrDeviceSystemID(client, d)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		interfaceId := d.Get("interface_id").(int)

		// 		// Gets the existing interface, thereby ensuring that the node and interface exists.
		// 		existingInterface, err := client.NetworkInterface.Get(systemId, interfaceId)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		existingTags := existingInterface.Tags

		// 		tfState := map[string]interface{}{
		// 			"id":       fmt.Sprintf("%s:%d", systemId, interfaceId),
		// 			"interface_id": interfaceId,
		// 			"tags":     existingTags,
		// 			"machine": d.Get("machine").(string),
		// 			"device": d.Get("device").(string),
		// 		}
		// 		if err := setTerraformState(d, tfState); err != nil {
		// 			return nil, err
		// 		}
		// 		return []*schema.ResourceData{d}, nil
		// 	},
		// },
		Schema: map[string]*schema.Schema{
			"device": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"machine", "device"},
				Description:  "The identifier (system ID, hostname, or FQDN) of the device with the network interface. Either `machine` or `device` must be provided.",
			},
			"machine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"machine", "device"},
				Description:  "The identifier (system ID, hostname, or FQDN) of the machine with the network interface. Either `machine` or `device` must be provided.",
			},
			"interface_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:     true,
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

func resourceNetworkInterfaceTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client

	interfaceId := d.Get("interface_id").(int)
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
	return resourceNetworkInterfaceTagRead(ctx, d, meta)
}

func resourceNetworkInterfaceTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client
	
	systemId, err := getMachineOrDeviceSystemID(client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	interfaceId := d.Get("interface_id").(int)

	// Get the existing interface
	existingInterface, err := client.NetworkInterface.Get(systemId, interfaceId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the tags in state
	d.Set("tags", existingInterface.Tags)
	d.Set("interface_id", interfaceId)
	d.Set("machine", d.Get("machine").(string))
	d.Set("device", d.Get("device").(string))

	return nil
}

func resourceNetworkInterfaceTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client
	systemId, err := getMachineOrDeviceSystemID(client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	interfaceId := d.Get("interface_id").(int)
	// Get the existing interface
	existingInterface, err := client.NetworkInterface.Get(systemId, interfaceId)
	if err != nil {
		return diag.FromErr(err)
	}
	
	existingTags := existingInterface.Tags
	desiredTags := convertToStringSlice(d.Get("tags").(*schema.Set).List())
	
	// Remove tags that are not in the specified set
	for _, tag := range existingTags {
		if !slices.Contains(desiredTags, tag) {
			client.NetworkInterface.RemoveTag(systemId, interfaceId, tag)
			} 
	}

	// Add tags that are in the specified set. AddTag will not add duplicates.
	for _, tag := range desiredTags {
		client.NetworkInterface.AddTag(systemId, interfaceId, tag)
	}
	return resourceNetworkInterfaceTagRead(ctx, d, meta)
}

func resourceNetworkInterfaceTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client
	systemId, err := getMachineOrDeviceSystemID(client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	interfaceId := d.Get("interface_id").(int)
	tags := convertToStringSlice(d.Get("tags").(*schema.Set).List())
	for _, t := range tags {
		_, err := client.NetworkInterface.RemoveTag(systemId, interfaceId, t)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return nil
}
