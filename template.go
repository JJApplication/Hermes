/*
Create: 2022/8/12
Project: Hermes
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/JJApplication/fushin/utils/json"
)

// html模板处理

const (
	TmplMgek      = "mgek.tmpl"
	TmplAlarm     = "alarm.tmpl"
	TmplAlarmHtml = "alarm.html.tmpl"
	TmplBlog      = "blog.tmpl"
	TmplHome      = "home.tmpl"
)

const (
	TmplErrRender = "error.tmpl"
)

func renderTemplate(args string, tmpl string) string {
	switch tmpl {
	case TmplMgek:
		return renderMgek(args)
	case TmplBlog:
		return renderBlog(args)
	case TmplAlarm:
		return renderAlarm(args)
	case TmplAlarmHtml:
		return renderAlarmHtml(args)
	case TmplHome:
		return renderHome(args)
	default:
		return errBody()
	}
}

func errBody() string {
	data, err := ioutil.ReadFile(getTmpl(TmplErrRender))
	if err != nil {
		return ""
	}
	return string(data)
}

func getTmpl(tmpl string) string {
	return "tmpl/" + tmpl
}

func renderMgek(args string) string {
	html, err := template.ParseFiles(getTmpl(TmplMgek))
	if err != nil {
		logger.ErrorF("render template %s error: %s", TmplMgek, err.Error())
		return errBody()
	}
	var buf bytes.Buffer
	err = html.Execute(&buf, nil)
	if err != nil {
		logger.ErrorF("execute template %s error: %s", TmplMgek, err.Error())
	}
	return buf.String()
}

// 传入json 数组[{url, title}]
func renderBlog(args string) string {
	var data []struct {
		Url   string `json:"url"`
		Title string `json:"title"`
	}
	err := json.Json.UnmarshalFromString(args, &data)
	if err != nil {
		return errBody()
	}
	html, err := template.ParseFiles(getTmpl(TmplBlog))
	if err != nil {
		logger.ErrorF("render template %s error: %s", TmplBlog, err.Error())
		return errBody()
	}
	var buf bytes.Buffer
	err = html.Execute(&buf, data)
	if err != nil {
		logger.ErrorF("execute template %s error: %s", TmplBlog, err.Error())
	}
	return buf.String()
}

func renderAlarm(args string) string {
	data := struct {
		Data string
	}{
		Data: args,
	}
	html, err := template.ParseFiles(getTmpl(TmplAlarm))
	if err != nil {
		logger.ErrorF("render template %s error: %s", TmplAlarm, err.Error())
		return errBody()
	}
	var buf bytes.Buffer
	err = html.Execute(&buf, data)
	if err != nil {
		logger.ErrorF("execute template %s error: %s", TmplAlarm, err.Error())
	}
	return buf.String()
}

// 默认可以直接渲染转义的html
// 所以对args进行html转换
func renderAlarmHtml(args string) string {
	data := struct {
		Data template.HTML
	}{
		Data: template.HTML(args),
	}
	html, err := template.ParseFiles(getTmpl(TmplAlarmHtml))
	if err != nil {
		logger.ErrorF("render template %s error: %s", TmplAlarm, err.Error())
		return errBody()
	}
	var buf bytes.Buffer
	err = html.Execute(&buf, data)
	if err != nil {
		logger.ErrorF("execute template %s error: %s", TmplAlarm, err.Error())
	}
	return buf.String()
}

func renderHome(args string) string {
	data := struct {
		Address string
	}{
		Address: args,
	}
	html, err := template.ParseFiles(getTmpl(TmplHome))
	if err != nil {
		logger.ErrorF("render template %s error: %s", TmplHome, err.Error())
		return errBody()
	}
	var buf bytes.Buffer
	err = html.Execute(&buf, data)
	if err != nil {
		logger.ErrorF("execute template %s error: %s", TmplHome, err.Error())
	}
	return buf.String()
}
