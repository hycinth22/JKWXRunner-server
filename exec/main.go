package main

import (
	"../model"
	"flag"
	"log"
	"os"
	"time"
)

var Debug bool

func init() {
	_, err := os.Stat("debug")
	Debug = !os.IsNotExist(err)
	if Debug {
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
		saveRunResult(account, result)
		randSleep(15*time.Second, 360*time.Second)
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
	username := flag.String("username", "", "username for running")
	flag.Parse()
	if *username != "" {
		account, err := model.GetAccountByUsername(*username)
		if err != nil {
			log.Println(err.Error())
			return
		}
		result := RunForAccount(account)
		saveRunResult(account, result)
	}
	RunOnce()
}
