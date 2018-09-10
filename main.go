package main

import (
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	//"./config"
	"github.com/erraa/fail2ban_metric/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fail2banlog string
	timesBanned int = 0
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

func parseCmd() float64 {
	cmd := exec.Command("fail2ban-client", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	outputs := strings.Split(string(output), "\n")
	for _, v := range outputs {
		if strings.Contains(v, ":") {
			numTimes := strings.TrimSpace(strings.Split(v, ":")[1])
			numTimesInt, err := strconv.Atoi(numTimes)
			if err != nil {
				panic(err)
			}
			return float64(numTimesInt)
		}
	}
	panic("Parsing failed")
}

func init() {
	prometheus.MustRegister(timesBannedprom)
}

func main() {
	var c config.Conf
	c.Parse("./config.yaml")
	fail2banlog = c.FailToBanLoc
	go func() {
		for {
			timesBannedprom.With(prometheus.Labels{"device": c.DeviceName}).Add(parseCmd())
			time.Sleep(time.Second * 10)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":6060", nil))
}
