package main

import (
	"flag"
	"github.com/inkedawn/JKWXFucker-server/model"
	"log"
	"os"
	"os/signal"
	"syscall"
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

func handleExit(c <-chan os.Signal, cExit chan<- bool) {
	_, ok := <-c
	if ok {
		log.Println("Exiting...")
		log.Println("We need some time to clean.")
		cExit <- true
	}
}

func resumeAccount(acc []*model.Account) {
	for _, acc := range acc {
		if acc.LastStatus == model.StatusReadyRun {
			log.Println("Resume account", acc.Username)
			acc.LastStatus = model.StatusOK
			err := model.UpdateAccount(acc)
			if err != nil {
				log.Println("resume account error, ", err.Error())
			}
		} else {
			log.Println("Skip account", acc.Username)
		}
	}
}

func RunOnce() {
	// 监听终止信号
	exitSigChan := make(chan os.Signal)
	exitChan := make(chan bool)
	signal.Notify(exitSigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go handleExit(exitSigChan, exitChan)

	accounts, err := model.ListAccountsTodayNotRun()
	if len(accounts) == 0 && Debug {
		log.Println("Debug list all")
		accounts, err = model.ListAccounts()
	}
	if err != nil {
		log.Println(err.Error())
		return
	}
accountRun:
	for i, account := range accounts {
		select {
		case <-exitChan:
			resumeAccount(accounts[i:])
			break accountRun
		default:
			log.Println("Run Task For Account ", account.Username)
			result := RunForAccount(account)
			err := saveRunResult(account, result)
			if err == nil {
				log.Println("Account ", account.Username, " task run Successfully.")
			} else {
				log.Println("Account ", account.Username, " task occur an Error: ", err.Error())
			}
		}
		sleepOverChan := make(chan bool)
		go func() {
			SleepALitte(int64(len(accounts)), 6*time.Hour)
			sleepOverChan <- true
		}()
		select {
		case <-sleepOverChan:
		case <-exitChan:
			log.Println("sleep interrupted")
			log.Printf("%d account need to resume \n", len(accounts[i:]))
			resumeAccount(accounts[i:])
			break accountRun
		}
	}

	// 停止监听终止信号，关闭channel通知相关goroutine退出
	signal.Stop(exitSigChan)
	close(exitSigChan)
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
		return
	}
	RunOnce()
}
