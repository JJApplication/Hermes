/*
Create: 2022/8/13
Project: Hermes
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"encoding/gob"
	"os"
)

// 任务的同步存储 加载
// 在服务正常down前存储全部的定时任务

func init() {
	gob.Register(Task{})
	gob.Register([]Task{})
	gob.Register(map[string]Task{})
}

const (
	TaskFile = "Hermes.task"
)

// 存储到gob
func toGob() {
	f, err := os.OpenFile(TaskFile, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		logger.ErrorF("open %s error: %s", TaskFile, err.Error())
		return
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(cronTaskMap)
	if err != nil {
		logger.ErrorF("sync %s error: %s", TaskFile, err.Error())
	}
}

// 从gob同步
// 仅在启动时调用
// 不存在文件时不加载
func loadGob() map[string]Task {
	if _, err := os.Stat(TaskFile); os.IsNotExist(err) {
		return nil
	}
	f, err := os.Open(TaskFile)
	if err != nil {
		logger.ErrorF("open %s error: %s", TaskFile, err.Error())
		truncGob()
		return nil
	}
	var t map[string]Task
	dec := gob.NewDecoder(f)
	err = dec.Decode(&t)
	if err != nil {
		logger.ErrorF("load %s error: %s", TaskFile, err.Error())
		truncGob()
		return nil
	}
	return t
}

// 如果无法从.task中同步 即认定该文件损坏
// 删除该文件等待重新同步
func truncGob() {
	if err := os.RemoveAll(TaskFile); err != nil {
		logger.ErrorF("trunc %s error: %s", TaskFile, err.Error())
	}
}
