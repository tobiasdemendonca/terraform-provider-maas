package maas

import (
	"context"
	// "fmt"
	// "slices"

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
			"interfaces": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of MAAS machine or device interfaces that will be tagged with the new tag.",
				Elem: &schema.Resource{
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
						"ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The network interface IDs to tag.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The new tag name. Because the name will be used in urls, it should be short.",
			},
		},
	}
}

func resourceInterfaceTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client
	
	// get relevant system id from inputs.
	interfaces := d.Get("interfaces").(*schema.Set)

	for _, iface := range interfaces.List() {
		iface := iface.(map[string]interface{})
		systemId, err := getSystemIDFromInterfaceMap(client, iface)
		if err != nil {
			return diag.FromErr(err)
		}
		// Add tag for each interface
		for _, id := range iface["ids"].([]int) {
			client.NetworkInterface.AddTag(systemId, id, d.Get("name").(string))
		}
	}

	// update state
	name := d.Get("name").(string)
	
	d.SetId(name)
	return resourceSimpleTagUpdate(ctx, d, meta)
}

func resourceInterfaceTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	
	return nil
}

func resourceInterfaceTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client
	
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


