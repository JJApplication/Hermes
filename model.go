/*
Create: 2022/8/12
Project: Hermes
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

// 邮件发送的报文结构体

type mailInfo struct {
	Type     string   `json:"type"`
	Message  string   `json:"message"`
	IsFile   bool     `json:"isFile"`
	Subject  string   `json:"subject"`
	Attach   []string `json:"attach"`
	To       []string `json:"to"`
	Cc       []string `json:"cc"`
	Bcc      []string `json:"bcc"`
	SyncTask bool     `json:"syncTask"`
	CronJob  string   `json:"cronJob"`
}

type MailInfo mailInfo

// 处理模板数据 当使用模板任务时 message中的信息为模板参数
type tmplArgs map[string]interface{}
