package v1beta1

import "io/ioutil"

func (ssh *SSHKey) Write(file string) {
	ioutil.WriteFile(file, ssh.PrivateKey, 0400)
	ioutil.WriteFile(file+".pub", ssh.PublicKey, 0400)
}
