package main

import (
	"github.com/sohlich/ticktock/security"
)

func main() {
	// Configure social
	config := security.SecurityConfig{
		map[string]security.OAuthConfig{
			"twitter": security.OAuthConfig{
				ClientID: "4GyXZFSRWZwC1ABvvMA8UUkgz",
				Secret:   "Yl2l5Ts0og3KkEsHrBFw1NxWQEFzGfGBTOZjfJTgGojK3cWIrX",
			},
			"github": security.OAuthConfig{
				ClientID: "72097818320c37312222",
				Secret:   "b06e37cb3e1f5011b45fda0871e75bfdcc393ca1",
			},
		},
		"secret",
	}

	security.Configure(config)
}
