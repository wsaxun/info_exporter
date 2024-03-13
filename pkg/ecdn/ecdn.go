package ecdn

import (
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"info_exporter/configs"
	"info_exporter/pkg/tools"
)

type Ecdn struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Ecdn) Fetch(path string, isSleep bool) ([]byte, error) {
	serverUrl, _ := url.JoinPath(configs.DOMAIN, path)

	t := time.Now().Unix()
	tokenStr := configs.GenToken(t)

	params := url.Values{}
	params.Add("timestamp", strconv.FormatInt(t, 10))
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	if isSleep {
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	}
	tools.LogInfo(map[string]interface{}{"api": serverUrl}, "Collect info")

	req, err := http.NewRequest("GET", serverUrl+"?"+params.Encode(), nil)
	if err != nil {
		tools.LogErr(map[string]interface{}{"err": err}, "Error creating request")
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+tokenStr)

	resp, err := client.Do(req)
	if err != nil {
		tools.LogErr(map[string]interface{}{"err": err}, "Error sending request")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			tools.LogErr(map[string]interface{}{"err": err}, "Body close err")
		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}
