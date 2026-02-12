package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	ping "github.com/prometheus-community/pro-bing"
)

type Result struct {
	Domain string
	Avg    string
	Jitter string
	IP     string
	Loss   string
}

func main() {
	domains := []string{"google.com", "wikipedia.org", "yandex.ru",
		"spotify.com", "youtube.com", "vk.com", "reddit.com",
		"twitch.tv", "pinterest.com", "dzen.ru"}
	results := make(chan Result, len(domains))

	for _, domain := range domains {
		go func(d string) {
			pinger, err := ping.NewPinger(d)
			if err != nil {
				results <- Result{d, "Error", "Error", "Error", "100%"}
				return
			}

			pinger.SetPrivileged(true)
			pinger.Count = 3
			pinger.Timeout = time.Second * 5

			err = pinger.Run()
			stats := pinger.Statistics()

			results <- Result{
				Domain: d,
				Avg:    fmt.Sprintf("%v", stats.AvgRtt),
				Jitter: fmt.Sprintf("%v", stats.StdDevRtt),
				IP:     pinger.IPAddr().String(),
				Loss:   fmt.Sprintf("%v%%", stats.PacketLoss),
			}
		}(domain)
	}

	file, _ := os.Create("ping.csv")
	writer := csv.NewWriter(file)
	writer.Write([]string{"Domain", "Avg_RTT", "Jitter", "IP", "Packet_Loss"})

	for i := 0; i < len(domains); i++ {
		res := <-results
		writer.Write([]string{res.Domain, res.Avg, res.Jitter, res.IP, res.Loss})
	}
	writer.Flush()
}

// привет
