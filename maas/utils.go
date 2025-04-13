package maas

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"strconv"
	"strings"

	"github.com/canonical/gomaasclient/client"
	"github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/gocty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func base64Encode(data []byte) string {
	if isBase64Encoded(data) {
		return string(data)
	}

	return base64.StdEncoding.EncodeToString(data)
}

func isBase64Encoded(data []byte) bool {
	_, err := base64.StdEncoding.DecodeString(string(data))
	return err == nil
}

func convertToStringSlice(field any) []string {
	if field == nil {
		return nil
	}

	fieldSlice := field.([]any)
	result := make([]string, len(fieldSlice))

	for i, value := range fieldSlice {
		result[i] = value.(string)
	}

	return result
}

func isElementIPAddress(i any, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attr := p[len(p)-1].(cty.IndexStep)

	var index int

	if err := gocty.FromCtyValue(attr.Key, &index); err != nil {
		return diag.FromErr(err)
	}

	ws, es := validation.IsIPAddress(i, fmt.Sprintf("element %v", index))

	for _, w := range ws {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Warning,
			Summary:       w,
			AttributePath: p,
		})
	}

	for _, e := range es {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       e.Error(),
			AttributePath: p,
		})
	}

	return diags
}

func isEmailAddress(i any, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attr := p[len(p)-1].(cty.GetAttrStep)

	v, ok := i.(string)
	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("expected type of %q to be string", attr.Name),
			AttributePath: p,
		})
	}

	if _, err := mail.ParseAddress(i.(string)); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("expected %s to be a valid e-mail address, got: %s", attr.Name, v),
			AttributePath: p,
		})
	}

	return diags
}

func getNetworkInterface(client *client.Client, machineSystemID string, identifier string) (*entity.NetworkInterface, error) {
	networkInterfaces, err := client.NetworkInterfaces.Get(machineSystemID)
	if err != nil {
		return nil, err
	}

	for _, n := range networkInterfaces {
		if n.MACAddress == identifier || n.Name == identifier || fmt.Sprintf("%v", n.ID) == identifier {
			return &n, nil
		}
	}

	return nil, fmt.Errorf("network interface (%s) was not found on machine (%s)", identifier, machineSystemID)
}

func setTerraformState(d *schema.ResourceData, tfState map[string]any) error {
	if val, ok := tfState["id"]; ok {
		d.SetId(val.(string))
		delete(tfState, "id")
	}

	for k, v := range tfState {
		// NOTE: Ignore R001. We have this method being invoked in multiple places,
		// however key values are actually string literals. Consider this a false positive.
		// https://github.com/bflad/tfproviderlint/tree/main/passes/R001
		//
		//lintignore:R001
		if err := d.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}

// Get the system ID of the relevant entity from resource data that accepts either a `machine` or `device`.
func getMachineOrDeviceSystemID(client *client.Client, d *schema.ResourceData) (string, error) {
	if d.Get("machine") != "" {
		machine, err := getMachine(client, d.Get("machine").(string))
		if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
			return "", err
		}

		return machine.SystemID, nil
	}

	if d.Get("device") != "" {
		device, err := getDevice(client, d.Get("device").(string))
		if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
			return "", err
		}

		return device.SystemID, nil
	}

	return "", fmt.Errorf("either `machine` or `device` is required")
}

// Get the type of the relevant entity from the system ID, by checking if the device or machine exists in MAAS.
// This gets all devices, then all machines, so is not efficient and should be used sparingly.
func getMachineOrDeviceTypeFromSystemID(client *client.Client, systemID string) (string, error) {
	device, err := getDevice(client, systemID)
	if err == nil && device != nil {
		return "device", nil
	}

	if !strings.Contains(err.Error(), fmt.Sprintf("device (%s) was not found", systemID)) {
		return "", fmt.Errorf("error getting device for system ID (%s): %w", systemID, err)
	}

	machine, err := getMachine(client, systemID)
	if err == nil && machine != nil {
		return "machine", nil
	}

	if !strings.Contains(err.Error(), fmt.Sprintf("machine (%s) was not found", systemID)) {
		return "", fmt.Errorf("error getting machine for system ID (%s): %w", systemID, err)
	}

	return "", fmt.Errorf("system ID (%s) was not found as either a device or machine", systemID)
}

func SplitStateIDIntoInts(stateID string, delimeter string) (int, int, error) {
	id1, id2, err := SplitStateID(stateID, delimeter)
	if err != nil {
		return 0, 0, err
	}

	id1Int, err := strconv.Atoi(id1)
	if err != nil {
		return 0, 0, err
	}

	id2Int, err := strconv.Atoi(id2)
	if err != nil {
		return 0, 0, err
	}

	return id1Int, id2Int, nil
}
func SplitStateID(stateID string, delimeter string) (string, string, error) {
	splitID := strings.SplitN(stateID, delimeter, 2)
	if len(splitID) != 2 {
		return "", "", fmt.Errorf("invalid resource ID: %s", stateID)
	}

	return splitID[0], splitID[1], nil
}
