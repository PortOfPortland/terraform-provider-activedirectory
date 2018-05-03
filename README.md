# Terraform Windows Active Directory Provider

This is the repository for a Terraform Active Directory Provider, which you can use to perform operations against Microsoft Active Directory.

The provider uses the [github.com/gorillalabs/go-powershell/backend](github.com/gorillalabs/go-powershell/backend) package to "shell out" to PowerShell, fire up a WinRM session, and perform the actual script work. I made this decision because the Go WinRM packages I was able to find only supported WinRM in Basic/Unencrypted mode, which is not doable in our environment. Shelling out to PowerShell is admittedly ugly, but it allows the use of domain accounts, HTTPS, etc.

# Using the Provider

### Example

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
}

#move an object to an ou
resource "activedirectory_ouMapping" "test1" {
  object_name = "MYVM1"
  object_class = "Computer"
  target_path = "OU=2016Servers,OU=Computers,OU=AD,DC=mydomain,DC=com"
}
```
