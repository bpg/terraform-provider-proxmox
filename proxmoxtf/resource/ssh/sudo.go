package ssh

import (
	"fmt"
)

const (
	TrySudo = `try_sudo(){ if [ $(sudo -n pvesm apiinfo 2>&1 | grep "APIVER" | wc -l) -gt 0 ]; then sudo $1; else $1; fi }`
)

func NewErrSSHUserNoPermission(username string) error {
	return fmt.Errorf("the SSH user '%s' does not have required permissions. "+
		"Make sure 'sudo' is installed and the user is configured in sudoers file. "+
		"Refer to the documentation for more details", username)
}
