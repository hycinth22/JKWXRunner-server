package main

import (
	"../model"
	"log"
	"os"
	"time"
)

var Debug bool

func init() {
	_, err := os.Stat("debug")
	Debug = !os.IsNotExist(err)
	if Debug{
		log.Println("DEBUG MODE")
	}
}

func RunOnce() {
	accounts, err := model.GetAllAccountsTodayNotRun()
	if len(accounts) == 0 && Debug {
		accounts, err = model.GetAllAccounts()
	}
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, account := range accounts {
		result := RunForAccount(account)
		account.LastTime, account.LastStatus, account.LastDistance = result.lastTime, result.status, result.lastDistance
		model.UpdateAccount(account)
	}
}

func RunAsDaemon() {
	lastRun := time.Now()
	for {
		now := time.Now()
		for lastRun.Day() == now.Day() {
		}
		lastRun = now
		RunOnce()
	}
}

func main() {
	RunOnce()
}
