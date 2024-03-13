package host

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"info_exporter/configs"
	"info_exporter/pkg/ecdn"
	"info_exporter/pkg/tools"
	"strconv"
	"sync"
	"time"
)

type HostsCollector struct {
	hostDesc *prometheus.Desc
	mu       sync.RWMutex
	ticker   *time.Ticker
	Info
}

func newHostsCollector() *HostsCollector {
	return &HostsCollector{
		hostDesc: prometheus.NewDesc("mfy_hosts_ecdn_info",
			"Collect all host information through ecdn",
			[]string{"id", "hostname", "sn", "parent", "ip", "isp", "location", "business", "owner", "status", "rate", "cactiNotes", "bwPlan", "ipmi"},
			nil,
		),
		ticker: time.NewTicker(configs.ECDNCYCLE),
	}
}

func (s *HostsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.hostDesc
}

func (s *HostsCollector) Collect(ch chan<- prometheus.Metric) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.Info.Items {
		metricValue := float64(StatusMap[item.Status])
		labels := []string{
			item.ID,
			item.Hostname,
			item.Sn,
			item.Parent,
			item.IP,
			item.Isp,
			item.Location,
			item.Business,
			item.Owner,
			item.Status,
			item.Rate,
			item.CactiNotes,
			strconv.FormatInt(item.BwPlan, 10),
			item.Ipmi,
		}
		m := prometheus.MustNewConstMetric(s.hostDesc, prometheus.GaugeValue, metricValue, labels...)
		ch <- m
	}
}

func (s *HostsCollector) CollectInfo() {
	items := Data{}
	body, err := items.Fetch(configs.HOSTPATH, true)
	if err != nil {
		return
	}

	if err := json.Unmarshal(body, &items); err != nil {
		tools.LogErr(map[string]interface{}{"err": err}, "Unmarshal error")
		return
	}
	if items.Code != 0 {
		tools.LogErr(map[string]interface{}{"err": items.Message}, "Collect err")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Info = items.Data
}

func (s *HostsCollector) Run() {
	s.CollectInfo()
	for range s.ticker.C {
		s.CollectInfo()
	}
}

type Data struct {
	ecdn.Ecdn
	Data Info `json:"data"`
}

type Info struct {
	Items []Item `json:"data"`
}

type Item struct {
	ID         string `json:"id"`
	Hostname   string `json:"hostname"`
	Sn         string `json:"sn"`
	Parent     string `json:"parent"`
	IP         string `json:"ip"`
	Isp        string `json:"isp"`
	Location   string `json:"location"`
	Business   string `json:"business"`
	Owner      string `json:"owner"`
	Status     string `json:"status"`
	Rate       string `json:"rate"`
	CactiNotes string `json:"cactiNotes"`
	BwPlan     int64  `json:"bwPlan"`
	Ipmi       string `json:"ipmi"`
}

// StatusMap 未知状态为0
var StatusMap = map[string]int{
	"在线":  1,
	"离线":  2,
	"已审核": 3,
}

// Monitor 采集hosts信息
func Monitor() {
	// 注册指标
	metric := newHostsCollector()
	prometheus.MustRegister(metric)
	go metric.Run()
}
