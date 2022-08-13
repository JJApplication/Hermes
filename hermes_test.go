/*
Create: 2022/8/12
Project: Hermes
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"testing"

	"github.com/JJApplication/fushin/log"
)

func init() {
	logger = log.Logger{
		Name:   Hermes,
		Option: log.DefaultOption,
		Sync:   true,
	}
	_ = logger.Init()
}

// 模板解析
func TestRenderTmpl(t *testing.T) {
	t.Log(errBody())
	t.Log(renderAlarm("This is a test Alarm"))
	t.Log(renderMgek(""))
	t.Log(renderHome("test@jj.com"))
	t.Log(renderBlog(`[{"url": "http://test.com", "title": "New Test"},{"url": "http://test.com", "title": "New Test2"}]`))
}

// 测试传递html字符串
func TestRenderAlarmHtml(t *testing.T) {
	t.Log(renderAlarmHtml("<h1>This is test Alarm</h1>"))
}

// 任务存储
func TestSyncCronTask(t *testing.T) {
	toGob()
	cronTaskMap = make(map[string]Task, 1)
	cronTaskMap["1"] = Task{
		UUID:        "1",
		Name:        "",
		Loop:        false,
		LastRun:     0,
		Event:       "",
		MailType:    "",
		MailTo:      nil,
		MailCc:      nil,
		MailBcc:     nil,
		MailSubject: "",
		MailMessage: "",
	}
}

// 任务加载
func TestLoadCronTask(t *testing.T) {
	// 初始化
	cronTaskMap = make(map[string]Task, 1)
	cronTaskMap["1"] = Task{
		UUID:        "1",
		Name:        "",
		Loop:        false,
		LastRun:     0,
		Event:       "",
		MailType:    "",
		MailTo:      nil,
		MailCc:      nil,
		MailBcc:     nil,
		MailSubject: "",
		MailMessage: "",
	}
	toGob()
	res := loadGob()
	t.Log(res)
}
