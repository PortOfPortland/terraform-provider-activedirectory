package activedirectory

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/portofportland/goWinRM"

	"strings"
)

func resourcegroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourcegroupMembershipCreate,
		Read:   resourcegroupMembershipRead,
		Delete: resourcegroupMembershipDelete,
		Update: resourcegroupMembershipCreate,

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
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcegroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	object_name := d.Get("object_name").(string)
	object_class := d.Get("object_class").(string)
	group_name := d.Get("group_name").(string)

	var id string = object_name + "_" + object_class + "_" + group_name

	var psCommand string = "$object = Get-ADObject -Filter {(name -eq '" + object_name + "') -AND (ObjectClass -eq '" + object_class + "')}; Add-ADGroupMember -Identity '" + group_name + "' -Members $object.DistinguishedName -Confirm:$false"
	_, err := goWinRM.RunWinRMCommand(client.username, client.password, client.server, psCommand, client.usessl)
	if err != nil {
		//something bad happened
		return err
	}

	d.SetId(id)

	return nil
}

func resourcegroupMembershipRead(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	object_name := d.Get("object_name").(string)
	object_class := d.Get("object_class").(string)
	group_name := d.Get("group_name").(string)

        var psCommand string = "$object = Get-ADGroupMember -Identity '" + group_name + "' | Where-Object {$_.Name -eq '" + object_name + "' -AND $_.objectClass -eq '" + object_class + "'}; if (!$object) { Write-Host 'TERRAFORM_NOT_FOUND' }"
	stdout, err := goWinRM.RunWinRMCommand(client.username, client.password, client.server, psCommand, client.usessl)
	if err != nil {
		//something bad happened
		return err
	}

	if strings.Contains(stdout, "TERRAFORM_NOT_FOUND") {
		//not able to find the record - this is an error but ok
		d.SetId("")
		return nil
	}

	var id string = object_name + "_" + object_class + "_" + group_name
	d.SetId(id)
	return nil
}

func resourcegroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	object_name := d.Get("object_name").(string)
	object_class := d.Get("object_class").(string)
	group_name := d.Get("group_name").(string)

	var psCommand string = "$object = Get-ADObject -Filter {(name -eq '" + object_name + "') -AND (ObjectClass -eq '" + object_class + "')}; Remove-ADGroupMember -Identity '" + group_name + "' -Members $object.DistinguishedName -Confirm:$false"
	_, err := goWinRM.RunWinRMCommand(client.username, client.password, client.server, psCommand, client.usessl)
	if err != nil {
		//something bad happened
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but it is added here for explicitness.
	d.SetId("")

	return nil
}