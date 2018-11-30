package main

import (
	"flag"
	"github.com/inkedawn/JKWXFucker-server/model"
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

// can't guarantee must be not timeout
func SleepALitte(totalAccountCount int64, totalExecTime time.Duration) {
	totalExecTime = time.Duration(0.8 * float64(totalExecTime)) // 20% for delay & other
	single := totalExecTime.Nanoseconds() / totalAccountCount

	var d time.Duration
	if time.Duration(single) > 5*time.Minute {
		d = randSleepDuration(15*time.Second, 5*time.Minute)
	} else {
		d = randSleepDuration(time.Duration(0.8*float64(single)), time.Duration(1.2*float64(single)))
	}

	log.Println("Sleep ", d.String())
	time.Sleep(d)
}

func RunOnce() {
	accounts, err := model.ListAccountsTodayNotRun()
	if len(accounts) == 0 && Debug {
		accounts, err = model.ListAccounts()
	}
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, account := range accounts {
		log.Println("Run Task For Account ", account.Username)
		result := RunForAccount(account)
		err := saveRunResult(account, result)
		if err == nil {
			log.Println("Account ", account.Username, " task run Successfully.")
		} else {
			log.Println("Account ", account.Username, " task occur an Error: ", err.Error())
		}
		SleepALitte(int64(len(accounts)), 6*time.Hour)
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
