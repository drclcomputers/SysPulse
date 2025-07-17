package sysinfo

import (
	"fmt"
	"os/user"
	"time"

	"github.com/shirou/gopsutil/host"
)

func GetUsers() {
	users, err := host.Users()
	if err != nil {
		panic(err)
	}

	for _, u := range users {
		loginTime := time.Unix(int64(u.Started), 0)
		fmt.Printf("User: %s, Terminal: %s, Host: %s, Login Time: %s\n", u.User, u.Terminal, u.Host, loginTime.Format(time.RFC1123))
	}
}

func GetCurrentUser() string {
	user, err := user.Current()
	if err != nil {
		return "Unknown"
	}

	return user.Username
}
