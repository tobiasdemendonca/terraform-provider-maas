package maas

import (
	"context"
	// "fmt"
	// "slices"
	// "strconv"
	// "strings"

	// "github.com/canonical/gomaasclient/client"
	// "github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMAASConfiguration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource to manage MAAS configuration settings through key value pairs.",
		CreateContext: resourceMAASConfigurationCreate,
		ReadContext:   resourceMAASConfigurationRead,
		UpdateContext: resourceMAASConfigurationUpdate,
		DeleteContext: resourceMAASConfigurationDelete,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:		schema.TypeString,
				Required:	 true,
				Description: "Key corresponding to the configuration setting.",
			},
			"value": {
				Type:		schema.TypeString,
				Required:	 true,
				Description: "Value for the configuration setting.",
			},
		},
	}
}

func resourceMAASConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	 client := meta.(*ClientConfig).Client

	//  client.MAASServer.Post()
	 
}
func resourceMAASConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	 return nil 
}
func resourceMAASConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceMAASConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}