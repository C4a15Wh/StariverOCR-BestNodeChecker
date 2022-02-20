package main

import (
	"time"

	"main.go/common"
	"main.go/model"
)

var Config model.Config

func main() {
	var IsMatch = false
	var RecordLineID string
	var RecordID string
	var BestRecord = ""

	common.Logger(0, "正在读取Config...")
	Config, err := common.ReadConfig()
	if err != nil {
		common.Logger(2, "致命错误！读取日志失败！")
		common.Logger(2, err.Error())
		return
	}

	HandleDomain := Config.HandleDomain.SubDomain + "." + Config.HandleDomain.RootDomain
	ResolveDomain := Config.ResolveDomain.SubDomain + "." + Config.ResolveDomain.RootDomain

	common.Logger(0, "控制域名: "+HandleDomain)
	common.Logger(0, "解析域名: "+ResolveDomain)

	for {
		LocalRecord, err := common.QueryDomainLocalRecord(ResolveDomain) // 查询解析域名
		if err != nil {
			common.Logger(1, "在查询DNS信息的时候遇到了以下问题："+err.Error())
			continue
		}

		common.Logger(0, "当前最佳纪录: "+LocalRecord[0])

		if LocalRecord[0] != BestRecord {
			common.Logger(0, "原纪录: "+BestRecord)
			BestRecord = LocalRecord[0]
		} else {
			common.Logger(0, "记录未更改。")
			time.Sleep(time.Duration(4) * time.Second)
			continue
		}

		RootDomainInfo, err := common.QueryDomainInfo(Config.HandleDomain.RootDomain)
		if err != nil {
			common.Logger(1, "在查询域名信息的时候遇到了以下问题："+err.Error())
			continue
		}

		if RootDomainInfo.Status.Code != "1" {
			common.Logger(1, "在查询域名信息的时候遇到了以下问题："+RootDomainInfo.Status.Message)
			continue
		}

		for _, key := range RootDomainInfo.Records {
			if key.Name == Config.HandleDomain.SubDomain {
				IsMatch = true
				RecordID = key.ID
				RecordLineID = key.LineID
				break
			}
		}

		if !IsMatch {
			common.Logger(1, "在查询域名信息的时候遇到了以下问题：该记录不存在")
			continue
		}

		common.Logger(0, "RecordID: "+RecordID)
		common.Logger(0, "LineID: "+RecordLineID)

		err = common.ChangeDomainRecord(Config.HandleDomain.RootDomain, RecordLineID, RecordID, BestRecord, Config.HandleDomain.SubDomain)
		if err != nil {
			common.Logger(1, "更改记录时遇到了以下问题："+err.Error())
			continue
		}

		common.Logger(0, "修改已完成。")
		time.Sleep(time.Duration(4) * time.Second)
	}
}
