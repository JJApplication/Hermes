//go:build linux || darwin

/*
Create: 2022/8/11
Project: Hermes
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"sort"

	"github.com/JJApplication/fushin/server/uds"
	"github.com/JJApplication/fushin/utils/env"
)

// hermes核心:
// 提供高可用的邮件通知服务
// 基于uds服务间通信, 支持定时邮件 邮件群发 特殊邮件
// HermesCore 定义结构体包含服务基本的信息

type HermesCore struct {
	AppName   string         // 服务的名称用于uds通信作为服务标识id
	Mail      *mail          // 邮件配置
	EnvLoader *env.EnvLoader // 环境变量加载器
	UdsServer *uds.UDSServer // uds服务器
}

type mail struct {
	UserMail   string
	UserPasswd string
	NickName   string
	Host       string
	Port       int
}

// 邮件配置无法从hermesCore中拿到 再次初始化
var mailConfig mail

// Init 初始化全部数据
func (h *HermesCore) Init() {
	// init mail
	h.Mail = &mail{
		UserMail:   h.EnvLoader.Get(User).Raw(),
		UserPasswd: h.EnvLoader.Get(Pass).Raw(),
		NickName:   h.EnvLoader.Get(Nickname).Raw(),
		Host:       h.EnvLoader.Get(Host).Raw(),
		Port:       h.EnvLoader.Get(Port).Int(),
	}

	mailConfig = *h.Mail
	// init uds
	h.UdsServer.Name = h.EnvLoader.Get(UnixAddress).Raw()
	h.UdsServer.Option.MaxSize = MaxRecvSize
	h.UdsServer.Option.AutoCheck = false
	h.UdsServer.Option.AutoRecover = true
	h.UdsServer.Option.LogTrace = true

	var v []string
	for k, _ := range events {
		v = append(v, k)
	}
	sort.SliceStable(v, func(i, j int) bool {
		return v[i] < v[j]
	})

	for _, o := range v {
		h.UdsServer.AddFunc(o, events[o])
		logger.InfoF("event [%s] is add to HermesCore", o)
	}
}

// Run 启动内部的uds服务器
func (h *HermesCore) Run() error {
	logger.InfoF("%s start to run @ [%s]", Hermes, h.EnvLoader.Get(UnixAddress).Raw())
	return h.UdsServer.Listen()
}
