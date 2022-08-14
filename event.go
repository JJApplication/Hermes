/*
Create: 2022/8/12
Project: Hermes
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"context"
	"time"

	"github.com/JJApplication/fushin/client/smtp"
	"github.com/JJApplication/fushin/server/uds"
	"github.com/JJApplication/fushin/utils/json"
)

var events map[string]uds.Func

const (
	MailSendTimeout = 10
)

func init() {
	events = make(map[string]uds.Func, 20)
	events["ping"] = eventPing()
	events["send"] = eventSend()
	events["sendSync"] = eventSendSync()
	events["sendCron"] = eventSendCron()
	events["sendSchedule"] = eventSendSchedule()
	events["sendMgek"] = eventSendMgek()
	events["sendAlarm"] = eventSendAlarm()
	events["sendAlarmHtml"] = eventSendAlarmHtml()
	events["sendHomeSub"] = eventSendHomeSub()
	events["sendBlogSub"] = eventSendBlogSub()
	events["tasks"] = eventTasks()
	events["cancelTask"] = eventCancelTask()
}

func eventPing() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		if err := c.Response(uds.Res{
			Error: "",
			Data:  "pong",
			From:  Hermes,
			To:    []string{req.From},
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

func eventSend() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		var reqBody mailInfo
		err := json.Json.UnmarshalFromString(req.Data, &reqBody)
		if err != nil {
			logger.ErrorF("event [%s] parse request error: %s", c.Operation(), err.Error())
			return
		}
		m := newSmtp(reqBody)
		if reqBody.Type == smtp.Html {
			err = m.SendHtml(reqBody.Subject, reqBody.Message, reqBody.IsFile, reqBody.Attach)
		} else {
			err = m.Send(reqBody.Subject, reqBody.Message, reqBody.Attach)
		}
		if err != nil {
			logger.ErrorF("send mail to %v error: %s", reqBody.To, err.Error())
		}
		if err := c.Response(uds.Res{
			Error: convertErr(err),
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 异步的发送 直接返回结果
func eventSendSync() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		var reqBody mailInfo
		err := json.Json.UnmarshalFromString(req.Data, &reqBody)
		if err != nil {
			logger.ErrorF("event [%s] parse request error: %s", c.Operation(), err.Error())
			return
		}
		runSyncTask(func() {
			m := newSmtp(reqBody)
			if reqBody.Type == smtp.Html {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*MailSendTimeout)
				defer cancel()
				err := m.SendContext(ctx, reqBody.Subject, reqBody.Message, reqBody.Attach)
				if err != nil {
					logger.ErrorF("send async mail task error: %s", err.Error())
				}
				select {
				case <-ctx.Done():
					logger.WarnF("send async mail task finished: %v", ctx.Err())
				default:
					logger.Info("send async mail task finished")
				}
			} else {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*MailSendTimeout)
				defer cancel()
				err := m.SendHtmlContext(ctx, reqBody.Subject, reqBody.Message, reqBody.IsFile, reqBody.Attach)
				if err != nil {
					logger.ErrorF("send async mail task error: %s", err.Error())
				}
				select {
				case <-ctx.Done():
					logger.WarnF("send async mail task finished: %v", ctx.Err())
				default:
					logger.Info("send async mail task finished")
				}
			}
		})
		if err = c.Response(uds.Res{
			Error: convertErr(err),
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 定时发送
// todo
func eventSendCron() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		if err := c.Response(uds.Res{
			Error: "",
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 计划发送
// todo
func eventSendSchedule() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		if err := c.Response(uds.Res{
			Error: "",
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 特殊邮件 定制的mgek通知邮件 使用本地模板生成

// mgek订阅
// message默认抛弃
func eventSendMgek() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		var reqBody mailInfo
		err := json.Json.UnmarshalFromString(req.Data, &reqBody)
		if err != nil {
			logger.ErrorF("event [%s] parse request error: %s", c.Operation(), err.Error())
			return
		}
		m := newSmtp(reqBody)
		body := renderTemplate(reqBody.Message, TmplMgek)
		err = m.SendHtml(reqBody.Subject, body, false, nil)
		if err := c.Response(uds.Res{
			Error: convertErr(err),
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 告警信息
// message为要告警的信息
func eventSendAlarm() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		var reqBody mailInfo
		err := json.Json.UnmarshalFromString(req.Data, &reqBody)
		if err != nil {
			logger.ErrorF("event [%s] parse request error: %s", c.Operation(), err.Error())
			return
		}
		m := newSmtp(reqBody)
		body := renderTemplate(reqBody.Message, TmplAlarm)
		err = m.SendHtml(reqBody.Subject, body, false, nil)
		if err := c.Response(uds.Res{
			Error: convertErr(err),
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// html格式的alarm
func eventSendAlarmHtml() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		var reqBody mailInfo
		err := json.Json.UnmarshalFromString(req.Data, &reqBody)
		if err != nil {
			logger.ErrorF("event [%s] parse request error: %s", c.Operation(), err.Error())
			return
		}
		m := newSmtp(reqBody)
		body := renderTemplate(reqBody.Message, TmplAlarmHtml)
		err = m.SendHtml(reqBody.Subject, body, false, nil)
		if err := c.Response(uds.Res{
			Error: convertErr(err),
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 博客订阅
// message 为
// []struct {
//	Url   string `json:"url"`
//	Title string `json:"title"`
//}
func eventSendBlogSub() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		var reqBody mailInfo
		err := json.Json.UnmarshalFromString(req.Data, &reqBody)
		if err != nil {
			logger.ErrorF("event [%s] parse request error: %s", c.Operation(), err.Error())
			return
		}
		m := newSmtp(reqBody)
		body := renderTemplate(reqBody.Message, TmplBlog)
		err = m.SendHtml(reqBody.Subject, body, false, nil)
		if err := c.Response(uds.Res{
			Error: convertErr(err),
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 主页订阅
// message为订阅者邮件地址
func eventSendHomeSub() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		var reqBody mailInfo
		err := json.Json.UnmarshalFromString(req.Data, &reqBody)
		if err != nil {
			logger.ErrorF("event [%s] parse request error: %s", c.Operation(), err.Error())
			return
		}
		m := newSmtp(reqBody)
		body := renderTemplate(reqBody.Message, TmplHome)
		err = m.SendHtml(reqBody.Subject, body, false, nil)
		if err := c.Response(uds.Res{
			Error: convertErr(err),
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 返回所有的后台定时任务
// sync任务会过滤掉 因为不保证异步的任务成功 异步任务都自带10s的timeout
func eventTasks() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		if err := c.Response(uds.Res{
			Error: "",
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 取消任务
func eventCancelTask() uds.Func {
	return func(c *uds.UDSContext, req uds.Req) {
		if err := c.Response(uds.Res{
			Error: "",
			Data:  "",
			From:  Hermes,
			To:    nil,
		}); err != nil {
			logger.ErrorF("event [%s] response error: %s", c.Operation(), err.Error())
		}
		logger.InfoF("event [%s] response ok", c.Operation())
	}
}

// 创建smtp客户端
func newSmtp(reqBody mailInfo) smtp.SmtpClient {
	return smtp.SmtpClient{
		Sender:     mailConfig.UserMail,
		NickSender: mailConfig.NickName,
		PassWord:   mailConfig.UserPasswd,
		SmtpHost:   mailConfig.Host,
		SmtpPort:   mailConfig.Port,
		To:         reqBody.To,
		Cc:         reqBody.Cc,
		Bcc:        reqBody.Bcc,
	}
}

func convertErr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
