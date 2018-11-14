package activedirectory

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/portofportland/goPSRemoting"

	"strings"
	"fmt"
)

func resourceOUMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceOUMappingCreate,
		Read:   resourceOUMappingRead,
		Delete: resourceOUMappingDelete,
		Update: resourceOUMappingCreate,

		Schema: map[string]*schema.Schema{
			"object_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object_class": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceOUMappingCreate(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	object_name := d.Get("object_name").(string)
	object_class := d.Get("object_class").(string)
	target_path := d.Get("target_path").(string)

	var id string = object_name + "_" + object_class + "_" + target_path

	var psCommand string = "Get-ADObject -Filter {(name -eq \\\"" + object_name + "\\\") -AND (ObjectClass -eq \\\"" + object_class + "\\\")} | Move-ADObject -TargetPath \\\"" + target_path + "\\\" -Confirm:$false"
	_, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh)
	if err != nil {
		//something bad happened
		return err
	}

	fmt.Println("THERE WAS AN ERROR")
	fmt.Println("THERE IS SOME OUT STUFF")

	d.SetId(id)

	return nil
}

func resourceOUMappingRead(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	object_name := d.Get("object_name").(string)
	object_class := d.Get("object_class").(string)
	target_path := d.Get("target_path").(string)

        var psCommand string = "$object = Get-ADObject -SearchBase \\\"" + target_path + "\\\" -Filter {(name -eq \\\"" + object_name + "\\\") -AND (ObjectClass -eq \\\"" + object_class + "\\\")}; if (!$object) { Write-Host 'TERRAFORM_NOT_FOUND' }"
	stdout, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh)
	if err != nil {
		//something bad happened
		return err
	}

	if strings.Contains(stdout, "TERRAFORM_NOT_FOUND") {
		//not able to find the record - this is an error but ok
		d.SetId("")
		return nil
	}

	var id string = object_name + "_" + object_class + "_" + target_path
	d.Set("address", id)
	return nil
}

func resourceOUMappingDelete(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	object_name := d.Get("object_name").(string)
	object_class := d.Get("object_class").(string)

	var psCommand string
	if (client.default_computer_container == "") {
		//move the computer to the default computer containers for the domain - Get-ADDomain | select computerscont*
		psCommand = "$container = Get-ADDomain | select computerscont*; Get-ADObject -Filter {(name -eq \\\"" + object_name + "\\\") -AND (ObjectClass -eq \\\"" + object_class + "\\\")} | Move-ADObject -TargetPath $container.ComputersContainer -Confirm:$false"
	} else {
		psCommand = "$container = Get-ADObject \\\"" + client.default_computer_container + "\\\"; Get-ADObject -Filter {(name -eq \\\"" + object_name + "\\\") -AND (ObjectClass -eq \\\"" + object_class + "\\\")} | Move-ADObject -TargetPath $container.DistinguishedName -Confirm:$false"
	}

	_, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh)
	if err != nil {
		//something bad happened
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but it is added here for explicitness.
	d.SetId("")

	return nil
}
