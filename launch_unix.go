//go:build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

/*
func getUsersFromPasswd() []string {
	out, err := exec.Command("cat", "/etc/passwd").Output()
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
	return []string{"sudo", "-iu", login}
}

func launch(argv []string) error {
	bin, err := exec.LookPath(argv[0])
	if err != nil {
		return err
	}
	return syscall.Exec(bin, argv, os.Environ())
}
