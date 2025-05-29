package maas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMAASConfiguration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource to manage MAAS configuration settings through key value pairs. See MAAS server in the [MAAS API documentation](https://maas.io/docs/api) for a full list of keys and their values. Upon destroy, the specified key will no longer be managed by Terraform, but the set value will remain set in MAAS.",
		CreateContext: resourceMAASConfigurationCreate,
		ReadContext:   resourceMAASConfigurationRead,
		UpdateContext: resourceMAASConfigurationUpdate,
		DeleteContext: resourceMAASConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key corresponding to the configuration setting.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value for the configuration setting.",
			},
		},
	}
}

func resourceMAASConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client

	err := client.MAASServer.Post(d.Get("key").(string), d.Get("value").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("key").(string))
	return resourceBlockDeviceRead(ctx, d, meta)

}

func resourceMAASConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client

	value, err := client.MAASServer.Get(d.Get("key").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	tfState := map[string]any{
		"value": string(value),
	}
	if err := setTerraformState(d, tfState); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMAASConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client

	err := client.MAASServer.Post(d.Get("key").(string), d.Get("value").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceMAASConfigurationRead(ctx, d, meta)
}

func resourceMAASConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Currently, settings cannot be reverted to their default values through the API. For now, just remove the key from the state.
	d.SetId("")
	return nil
}
