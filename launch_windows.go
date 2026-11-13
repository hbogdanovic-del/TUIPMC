//go:build windows

package main

import (
	"errors"
	"os"
	"os/exec"
)

/*
func getUsersFormPasswd() []string {
	//	out, err := exec.Command("cat", "/etc/passwd").Output()
	out, err := os.ReadFile("passwd")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	l_out := strings.Split(string(out), "\n")
	l_out = l_out[:len(l_out)-1]
	users := []string{}
	for _, line := range l_out {
		user := strings.Split(line, ":")[0]
		users = append(users, user)
	}
	return users
}
*/

func buildArgv(login string) []string {
	return []string{"cmd", "/c", "echo", "Compte choisi", login}
}

func launch(argv []string) error {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		os.Exit(ee.ExitCode())
	}
	return err
}
