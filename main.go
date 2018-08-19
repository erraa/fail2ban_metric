package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fail2banlog string = "./fail2ban.log"
	timesBanned int    = 0
)

var (
	timesBannedprom = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "times_banned",
			Help: "Number of times someone has tried to hack the machine",
		},
		[]string{"device"},
	)
)

func scan(s string) (bool, string) {
	if strings.Contains(s, "Ban") {
		return true, s
	} else {
		return false, ""
	}
}

func parseLog() float64 {
	file, err := os.Open(fail2banlog)
	if err != nil {
		panic("couldnt open fail2banlog")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	bans := []string{}
	for scanner.Scan() {
		result, line := scan(scanner.Text())
		if result {
			bans = append(bans, line)
		}
	}

	if timesBanned == 0 {
		timesBanned = len(bans)
		return float64(len(bans))
	} else if len(bans) == timesBanned {
		return float64(0)
	} else {
		return_value := len(bans) - timesBanned
		timesBanned = len(bans)
		return float64(return_value)
	}
}

func init() {
	prometheus.MustRegister(timesBannedprom)
}

func main() {
	go func() {
		for {
			timesBannedprom.With(prometheus.Labels{"device": "firewall"}).Add(parseLog())
			time.Sleep(time.Second * 2)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":6060", nil))
}
