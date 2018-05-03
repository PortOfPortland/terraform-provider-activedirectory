package activedirectory

import (
	"github.com/hashicorp/terraform/helper/schema"

	ps "github.com/gorillalabs/go-powershell"
	"github.com/gorillalabs/go-powershell/backend"

	//"errors"
	//"strings"
)

func resourceOUMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceOUMappingCreate,
		Read:   resourceOUMappingRead,
		Delete: resourceOUMappingDelete,

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
				ForceNew: true,
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

	var psCommand string = "Get-ADObject -Filter {(name -eq '" + object_name + "') -AND (ObjectClass -eq '" + object_class + "')} | Move-ADObject -TargetPath '" + target_path + "' -Confirm:$false"
	_, err := runWinRMCommand(client.username, client.password, client.server, psCommand, client.usessl)
	if err != nil {
		//something bad happened
		return err
	}

	d.SetId(id)

	return nil
}

func resourceOUMappingRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceOUMappingDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func runWinRMCommand(username string, password string, server string, command string, usessl string) (string, error) {
	// choose a backend
	back := &backend.Local{}

	// start a local powershell process
	shell, err := ps.New(back)
	if err != nil {
		//something bad happened - return an error
		return "", err
	}
	defer shell.Exit()

	// ... and interact with it
	var winRMPre string = "$SecurePassword = '" + password + "' | ConvertTo-SecureString -AsPlainText -Force; $cred = New-Object System.Management.Automation.PSCredential -ArgumentList '" + username + "', $SecurePassword; $s = New-PSSession -ComputerName " + server + " -Credential $cred"
        var winRMPost string = "; Invoke-Command -Session $s -Scriptblock { " + command + " }; Remove-PSSession $s"

	// use SSL if requested
	var winRMCommand string
	if (usessl == "1") {
		winRMCommand = winRMPre + " -UseSSL" + winRMPost
	} else {
		winRMCommand = winRMPre + winRMPost
	}
	stdout, _, err := shell.Execute(winRMCommand)
	
	return stdout, err
}