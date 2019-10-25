package main

import "github.com/inkedawn/JKWXRunner-server/debug"

func init() {
	debug.SetupSigUsr1Trap()
}
