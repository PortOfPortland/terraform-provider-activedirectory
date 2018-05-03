# Terraform Windows Active Directory Provider

This is the repository for a Terraform Active Directory Provider, which you can use to perform operations against Microsoft Active Directory.

The provider uses the [github.com/gorillalabs/go-powershell/backend](github.com/gorillalabs/go-powershell/backend) package to "shell out" to PowerShell, fire up a WinRM session, and perform the actual script work. I made this decision because the Go WinRM packages I was able to find only supported WinRM in Basic/Unencrypted mode, which is not doable in our environment. Shelling out to PowerShell is admittedly ugly, but it allows the use of domain accounts, HTTPS, etc.

# Using the Provider

### Configuring the Provider

```hcl
# configure the provider
# username + password - used to build a powershell credential
# server - the server we'll create a WinRM session into to perform the AD operations
# usessl - whether or not to use HTTPS for our WinRM session (by default port TCP/5986)
variable "username" {
  type = "string"
}

variable "password" {
  type = "string"
}

provider "activedirectory" {
  server = "mydc.mydomain.com"
  username = "${var.username}"
  password = "${var.password}"
  usessl = true
  default_computer_container = "Computers,OU=Computers,OU=AD,DC=mydomain,DC=com" #optional
}
```

### Example - Moving an Object to an OU/Container

This doesn't perfectly map to the CRUD workflow, so in general, it works like this:

* CREATE - Moves a resource to the correct OU/Container
* READ - Verifies whether/not a resource is in the correct OU/Container
* UPDATE - Calls the create function to move the resource to the correct OU/Container
* DELETE - Moves the resource to either the default computers container (via Get-ADDomain | select computerscont*) or the specified default_computer_container property on the provider

```hcl
#move an object to an ou.container
resource "activedirectory_ouMapping" "test1" {
  object_name = "MYVM1"
  object_class = "Computer"
  target_path = "OU=2016Servers,OU=Computers,OU=AD,DC=mydomain,DC=com"
}
```

### Example - Adding Objects to an AD Group

```
resource "activedirectory_groupMembership" "test1" {
  object_name = "MVM1"
  object_class = "Computer"
  group_name = "A really cool group"
}
```
