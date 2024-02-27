package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Status struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		NodeInfo struct {
			ProtocolVersion struct {
				P2p   string `json:"p2p"`
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"protocol_version"`
			Id          string            `json:"id"`
			ListenAddr  string            `json:"listen_addr"`
			Network     string            `json:"network"`
			Version     string            `json:"version"`
			Channels    string            `json:"channels"`
			Moniker     string            `json:"moniker"`
			Other       map[string]string `json:"other"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHash     string `json:"latest_block_hash"`
			LatestAppHash       string `json:"latest_app_hash"`
			LatestBlockHeight   string `json:"latest_block_height"`
			LatestBlockTime     string `json:"latest_block_time"`
			EarliestBlockHash   string `json:"earliest_block_hash"`
			EarliestAppHash     string `json:"earliest_app_hash"`
			EarliestBlockHeight string `json:"earliest_block_height"`
			EarliestBlockTime   string `json:"earliest_block_time"`
			CatchingUp          bool   `json:"catching_up"`
		} `json:"sync_info"`
		ValidatorInfo struct {
			Address     string            `json:"address"`
			PubKey      map[string]string `json:"pub_key"`
			VotingPower string            `json:"voting_power"`
		} `json:"validator_info"`
	} `json:"result"`
}

type Exporter struct {
	URI                    string
	mutex                  sync.RWMutex
	latestBlockHeightGauge prometheus.Gauge
	latestBlockTimeGauge   prometheus.Gauge
	earliestBlockHeightGauge prometheus.Gauge
	earliestBlockTimeGauge   prometheus.Gauge
	catchingUpGauge        prometheus.Gauge
	votingPowerGauge       prometheus.Gauge
	networkInfoGauge       prometheus.Gauge
}

func NewExporter(uri string) *Exporter {
	return &Exporter{
		URI: uri,
		latestBlockHeightGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "latest_block_height",
			Help: "The latest block height",
		}),
		latestBlockTimeGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "latest_block_time",
			Help: "The latest block time",
		}),
		earliestBlockHeightGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "earliest_block_height",
			Help: "The earliest block height",
		}),
		earliestBlockTimeGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "earliest_block_time",
			Help: "The earliest block time",
		}),
		catchingUpGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "catching_up",
			Help: "Whether the node is catching up",
		}),
		votingPowerGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "voting_power",
			Help: "The validator's voting power",
		}),
		networkInfoGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "network_info",
			Help: "Information about the network",
		}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.latestBlockHeightGauge.Describe(ch)
	e.latestBlockTimeGauge.Describe(ch)
	e.earliestBlockHeightGauge.Describe(ch)
	e.earliestBlockTimeGauge.Describe(ch)
	e.catchingUpGauge.Describe(ch)
	e.votingPowerGauge.Describe(ch)
	e.networkInfoGauge.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	data, err := ioutil.ReadFile(e.URI)
	if err != nil {
		return
	}

	var status Status
	if err := json.Unmarshal(data, &status); err != nil {
		return
	}

	latestBlockHeight, err := strconv.ParseFloat(status.Result.SyncInfo.LatestBlockHeight, 64)
	if err != nil {
		return
	}
	e.latestBlockHeightGauge.Set(latestBlockHeight)

	earliestBlockHeight, err := strconv.ParseFloat(status.Result.SyncInfo.EarliestBlockHeight, 64)
	if err != nil {
		return
	}
	e.earliestBlockHeightGauge.Set(earliestBlockHeight)

	catchingUp := 0.0
	if status.Result.SyncInfo.CatchingUp {
		catchingUp = 1.0
	}
	e.catchingUpGauge.Set(catchingUp)

	votingPower, err := strconv.ParseFloat(status.Result.ValidatorInfo.VotingPower, 64)
	if err != nil {
		return
	}
	e.votingPowerGauge.Set(votingPower)

	e.latestBlockHeightGauge.Collect(ch)
	e.latestBlockTimeGauge.Collect(ch)
	e.earliestBlockHeightGauge.Collect(ch)
	e.earliestBlockTimeGauge.Collect(ch)
	e.catchingUpGauge.Collect(ch)
	e.votingPowerGauge.Collect(ch)
	e.networkInfoGauge.Collect(ch)
}

func main() {
	exporter := NewExporter("status.txt")
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
