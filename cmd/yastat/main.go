package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newonton/yastat/internal/api"
	"github.com/newonton/yastat/internal/config"
	"github.com/newonton/yastat/internal/csvstore"
	"github.com/newonton/yastat/internal/lock"
	"github.com/newonton/yastat/internal/scheduler"
	"github.com/newonton/yastat/internal/timezone"
	"github.com/newonton/yastat/internal/ui"
)

func main() {
	var (
		period time.Duration
		help   bool
	)

	flag.DurationVar(&period, "period", time.Minute, "update period (1m, 1h, 1d)")
	flag.DurationVar(&period, "p", time.Minute, "update period (shorthand)")
	flag.BoolVar(&help, "help", false, "show help")
	flag.BoolVar(&help, "h", false, "show help (shorthand)")

	flag.Usage = usage
	flag.Parse()

	if help {
		usage()
		return
	}

	if flag.NArg() != 1 {
		fmt.Println("error: csv file is required")
		usage()
		os.Exit(1)
	}

	if period < time.Minute {
		fmt.Println("error: minimal period is 1m")
		os.Exit(1)
	}

	csvFile := flag.Arg(0)

	ui.HideCursor()
	defer ui.ShowCursor()

	cfg, err := config.MustLoad()
	if err != nil {
		fmt.Println(err)
		return
	}

	lockFile, err := lock.Acquire()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer lock.Release(lockFile)

	client := api.Client{
		APIKey: cfg.APIKey,
		AppID:  cfg.AppID,
	}

	start := time.Now().In(timezone.MoscowZone)

	last, ok := csvstore.LastTime(csvFile)
	if !ok {
		write(&client, csvFile)
		last = time.Now().In(timezone.MoscowZone)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		next := scheduler.NextRun(last, period)

		for time.Now().Before(next) && run {
			ui.Render(start, last, next)
			time.Sleep(time.Second)

			select {
			case <-stop:
				run = false
			default:
			}
		}

		if !run {
			break
		}

		write(&client, csvFile)
		last = time.Now().In(timezone.MoscowZone)
	}

	fmt.Println("\nReceived shutdown signal, writing final record...")
	write(&client, csvFile)
	fmt.Println("Shutdown complete.")
}

func write(c *api.Client, file string) {
	shows, reward, _ := c.Fetch()
	csvstore.Append(file, time.Now().In(timezone.MoscowZone), reward, shows)
}

func usage() {
	fmt.Println(`yastat — Yandex Ads statistics collector

Usage:
  yastat [flags] <file.csv>

Flags:
  -p, --period   update period (default 1m)
  -h, --help     show this help

Examples:
  yastat -p 5m stats.csv
  yastat --period=1h stats.csv

Config:
  ~/.config/yastat/config.json

Get api_key:
  https://oauth.yandex.ru/authorize?response_type=token&client_id=da092e6d50b443308da7a28e638070b9`)
}
