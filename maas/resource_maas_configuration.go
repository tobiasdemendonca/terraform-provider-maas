package maas

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Key corresponding to the configuration setting.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringDoesNotMatch(regexp.MustCompile("maas_auto_ipmi_workaround_flags"), "Key 'maas_auto_ipmi_workaround_flags' cannot currently be set through Terraform due to a [bug in MAAS](https://bugs.launchpad.net/maas/+bug/2112191). A fix will be worked on for a future MAAS releases.")),
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value for the configuration setting, always specified as a string",
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if d.Get("key").(string) == "remote_syslog" {
						if newValue == "" {
							return oldValue == "null"
						}
					}
					return oldValue == newValue
				},
			},
		},
	}
}

func resourceMAASConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("running resourceMAASConfigurationCreate")

	client := meta.(*ClientConfig).Client

	log.Println("Getting values from state")

	key := d.Get("key").(string)
	value := d.Get("value").(string)

	log.Println("Setting key and value in MAAS")

	err := client.MAASServer.Post(key, value)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting key %v with value %v in MAAS: %v", key, value, err))
	}

	log.Println("Setting ID in state")
	d.SetId(key)

	return resourceMAASConfigurationRead(ctx, d, meta)
}

func resourceMAASConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client

	key := d.Id()

	value, err := client.MAASServer.Get(key)
	if err != nil {
		return diag.FromErr(err)
	}

	tfState := map[string]any{
		"value": strings.Trim(string(value), "\""),
		"key":   key,
	}
	if err := setTerraformState(d, tfState); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMAASConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("running resourceMAASConfigurationUpdate")

	client := meta.(*ClientConfig).Client

	err := client.MAASServer.Post(d.Get("key").(string), d.Get("value").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMAASConfigurationRead(ctx, d, meta)
}

func resourceMAASConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("running resourceMAASConfigurationDelete")
	// Currently, settings cannot be reverted to their default values through the API. For now, just remove the key from the state.
	d.SetId("")

	return nil
}
