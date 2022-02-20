package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"main.go/model"
)

var FullLog string
var Config model.Config

func Logger(LogType int, log string) {
	switch LogType {
	case 0: // INFO
		log = "[INFO] " + " " + log
		color.White(log)

	case 1:
		log = "[WARN] " + " " + log
		color.Yellow(log)

	case 2:
		log = "[ERROR] " + log
		color.Red(log)

	default:
		log = "[DEFAULT] " + " " + log
		color.White(log)
	}

	FullLog = FullLog + "\n" + log
	ioutil.WriteFile("./NodeChecker.log", []byte(FullLog), 0777)
}

func JsonifyPostRequest(Endpoint string, Request interface{}) ([]byte, error) {
	ByterizationRequest, err := json.Marshal(Request)
	if err != nil {
		return nil, err
	}

	RequestReader := bytes.NewReader(ByterizationRequest)
	HTTPRequest, err := http.NewRequest("POST", Endpoint, RequestReader)
	if err != nil {
		return nil, err
	}
	HTTPRequest.Header.Set("Content-Type", "application/json")

	Client := http.Client{}
	Response, err := Client.Do(HTTPRequest)
	if err != nil {
		return nil, err
	}

	ByterizationResponse, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		return nil, err
	}

	return ByterizationResponse, nil
}

func FormPostRequest(Endpoint string, Request url.Values) ([]byte, error) {
	Response, err := http.PostForm(Endpoint, Request)
	if err != nil {
		return nil, err
	}

	ByterizationResponse, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		return nil, err
	}

	return ByterizationResponse, nil
}

func ReadConfig() (model.Config, error) {
	WorkPath, _ := os.Getwd()
	ConfigFilePath := path.Join(WorkPath, "config", "config.yaml")
	ByterizationConfig, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		return model.Config{}, err
	}

	err = yaml.Unmarshal(ByterizationConfig, &Config)
	if err != nil {
		return model.Config{}, err
	}

	return Config, nil
}

func DefinePublicParams(params url.Values) url.Values {
	params.Set("login_token", Config.Token)
	params.Set("format", "json")
	params.Set("lang", "cn")
	return params
}

// 通过本机发送DNS请求查询记录
func QueryDomainLocalRecord(domain string) ([]string, error) {
	var DomainDNSResult = make([]string, 0)

	record, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	if len(record) <= 0 {
		err = errors.New("域名无A级记录")
		return nil, err
	}

	for _, ip := range record {
		DomainDNSResult = append(DomainDNSResult, ip.String())
	}

	return DomainDNSResult, nil
}

// 查询该域名在DNSPod的相关信息
func QueryDomainInfo(domain string) (model.DomainInfo, error) {
	var QueryResponse model.DomainInfo

	RequestValue := DefinePublicParams(make(url.Values))

	RequestValue.Set("domain", domain)

	Response, err := FormPostRequest("https://dnsapi.cn/Record.List", RequestValue)
	if err != nil {
		return model.DomainInfo{}, err
	}

	err = json.Unmarshal(Response, &QueryResponse)
	if err != nil {
		return model.DomainInfo{}, err
	}

	return QueryResponse, nil
}

func ChangeDomainRecord(doamin, line_id, record_id, value, subdomain string) error {
	var QueryResponse model.DomainInfo

	RequestValue := DefinePublicParams(make(url.Values))
	RequestValue.Set("domain", doamin)
	RequestValue.Set("record_line_id", line_id)
	RequestValue.Set("record_id", record_id)
	RequestValue.Set("value", value)
	RequestValue.Set("sub_domain", subdomain)
	RequestValue.Set("record_type", "A")

	Response, err := FormPostRequest("https://dnsapi.cn/Record.Modify", RequestValue)
	if err != nil {
		return err
	}

	err = json.Unmarshal(Response, &QueryResponse)
	if err != nil {
		return err
	}

	if QueryResponse.Status.Code != "1" {
		return errors.New(QueryResponse.Status.Message)
	}

	return nil
}
