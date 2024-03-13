package switches

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

type SwitchCollector struct {
	switchDesc *prometheus.Desc
	mu         sync.RWMutex
	ticker     *time.Ticker
	Info
}

func newSwitchCollector() *SwitchCollector {
	return &SwitchCollector{
		switchDesc: prometheus.NewDesc("mfy_switch_ecdn_info",
			"Collect all switch information through ecdn",
			[]string{"id", "ip", "vpnIP", "name", "owner", "status", "isMain", "snmpCommunity", "snmpVersion", "snmpPort"},
			nil,
		),
		ticker: time.NewTicker(configs.ECDNCYCLE),
	}
}

func (s *SwitchCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.switchDesc
}

func (s *SwitchCollector) Collect(ch chan<- prometheus.Metric) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.Info.Items {
		metricValue := float64(StatusMap[item.Status])
		labels := []string{
			strconv.Itoa(item.ID),
			item.IP,
			item.VpnIP,
			item.Name,
			item.Owner,
			item.Status,
			item.IsMain,
			item.SnmpCommunity,
			strconv.Itoa(item.SnmpVersion),
			strconv.Itoa(item.SnmpPort),
		}
		m := prometheus.MustNewConstMetric(s.switchDesc, prometheus.GaugeValue, metricValue, labels...)
		ch <- m
	}
}

func (s *SwitchCollector) CollectInfo() {
	items := Data{}
	body, err := items.Fetch(configs.SWITCH, true)
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

func (s *SwitchCollector) Run() {
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
	Items []Item `json:"list"`
}

type Item struct {
	ID            int    `json:"id"`
	IP            string `json:"ip"`
	VpnIP         string `json:"vpnIP"`
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	Status        string `json:"status"`
	IsMain        string `json:"isMain"`
	SnmpCommunity string `json:"snmpCommunity"`
	SnmpVersion   int    `json:"snmpVersion"`
	SnmpPort      int    `json:"snmpPort"`
}

// StatusMap 未知状态为0
var StatusMap = map[string]int{
	"在线": 1,
	"离线": 2,
}

// Monitor 采集switch信息
func Monitor() {
	// 注册指标
	metric := newSwitchCollector()
	prometheus.MustRegister(metric)
	go metric.Run()
}
