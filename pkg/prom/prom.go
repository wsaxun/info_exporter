package prom

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"info_exporter/configs"
	"info_exporter/pkg/tools"
)

func Run() {
	// 启动一个 HTTP 服务器
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(configs.LISTENSERSER, nil); err != nil {
			tools.LogFatal(map[string]interface{}{"err": err})
		}
	}()
}
