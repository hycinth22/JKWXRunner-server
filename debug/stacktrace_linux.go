package debug

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupSigUsr1Trap() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go func() {
		for range c {
			DumpStacks()
		}
	}()
}
