package modules

import (
	"fmt"
	"os/exec"
	"os/user"

	"github.com/wx13/genesis"
)

type User struct {
	Name, Passwd string
}

func (u User) ID() string {
	return fmt.Sprintf("User: %s *****", u.Name)
}

func (u User) Files() []string {
	return []string{}
}

func (u User) Status() (genesis.Status, string, error) {
	_, err := user.Lookup(u.Name)
	if err != nil {
		return genesis.StatusFail, "User does not exist.", err
	}
	return genesis.StatusPass, "User exists.", nil
}

func (u User) Remove() (string, error) {
	usrStr := fmt.Sprintf("'%s'", u.Name)
	cmd := exec.Command("userdel", usrStr)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (u User) Install() (string, error) {

	// Create the user.
	usrStr := fmt.Sprintf("'%s'", u.Name)
	cmd := exec.Command("useradd", "-m", usrStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	// Set the password.
	cmd = exec.Command("chpasswd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "Could not create command pipe.", err
	}
	err = cmd.Start()
	if err != nil {
		return "Could not run password setter.", err
	}
	str := fmt.Sprintf(`"%s:%s"`, u.Name, u.Passwd)
	_, err = stdin.Write([]byte(str))
	if err != nil {
		return "Could not write to stdin.", err
	}
	stdin.Close()
	err = cmd.Wait()
	if err != nil {
		return "Could not set user password.", err
	}
	return "Created user.", nil

}
