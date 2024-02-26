package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"./gotify_notifier"

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
	URI string
	mu  sync.Mutex
}


func (e *Exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.mu.Lock()
	defer e.mu.Unlock()

	data, err := ioutil.ReadFile(e.URI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var status Status
	if err := json.Unmarshal(data, &status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	latestBlockHeight, err := strconv.ParseFloat(status.Result.SyncInfo.LatestBlockHeight, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	earliestBlockHeight, err := strconv.ParseFloat(status.Result.SyncInfo.EarliestBlockHeight, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	votingPower, err := strconv.ParseFloat(status.Result.ValidatorInfo.VotingPower, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	fmt.Fprintf(w, "latest_block_height %f\n", latestBlockHeight)
	fmt.Fprintf(w, "latest_block_time %s\n", status.Result.SyncInfo.LatestBlockTime)
	fmt.Fprintf(w, "earliest_block_height %f\n", earliestBlockHeight)
	fmt.Fprintf(w, "earliest_block_time %s\n", status.Result.SyncInfo.EarliestBlockTime)
	fmt.Fprintf(w, "catching_up %t\n", status.Result.SyncInfo.CatchingUp)
	fmt.Fprintf(w, "voting_power %f\n", votingPower)
	fmt.Fprintf(w, "network %s\n", status.Result.NodeInfo.Network)
	fmt.Fprintf(w, "moniker %s\n", status.Result.NodeInfo.Moniker)
	fmt.Fprintf(w, "version %s\n", status.Result.NodeInfo.Version)
	fmt.Fprintf(w, "channels %s\n", status.Result.NodeInfo.Channels)
	fmt.Fprintf(w, "p2p_protocol_version %s\n", status.Result.NodeInfo.ProtocolVersion.P2p)
	fmt.Fprintf(w, "block_protocol_version %s\n", status.Result.NodeInfo.ProtocolVersion.Block)
	fmt.Fprintf(w, "app_protocol_version %s\n", status.Result.NodeInfo.ProtocolVersion.App)
	fmt.Fprintf(w, "listen_addr %s\n", status.Result.NodeInfo.ListenAddr)
	fmt.Fprintf(w, "node_id %s\n", status.Result.NodeInfo.Id)
	fmt.Fprintf(w, "validator_address %s\n", status.Result.ValidatorInfo.Address)
}

func main() {
	
	exporter := &Exporter{
		URI: "status.txt", 
	}

	
	http.Handle("/metrics", exporter)


	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
