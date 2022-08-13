/*
Create: 2022/8/12
Project: Hermes
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"sync"

	stringUtil "github.com/JJApplication/fushin/utils/string"
)

// 异步执行
// 不保证执行结果 记录日志

var syncTaskMap map[string]struct{}
var cronTaskMap map[string]Task
var lock sync.Mutex
var once sync.Once

func init() {
	// 每次启动时执行
	once.Do(func() {
		tasks := loadGob()
		logger.Info("reload cronTasks from ", TaskFile)
		if tasks != nil {
			cronTaskMap = tasks
		} else {
			cronTaskMap = make(map[string]Task, 1)
		}
	})
}

// Task 定时任务
type Task struct {
	UUID        string // 唯一区分任务的属性
	ID          int    // 定时任务的CronID
	Name        string // 任务名称
	Loop        bool   // 是否轮询执行
	LastRun     int64  // 上次执行的时间
	Event       string // 触发的event
	MailType    string
	MailTo      []string // 接收方
	MailCc      []string
	MailBcc     []string
	MailSubject string
	MailMessage string
	MailAttach  []string
}

func runSyncTask(task func()) {
	lock.Lock()
	uuid := stringUtil.UUID()
	lock.Unlock()
	syncTaskMap[uuid] = struct{}{}
	go task()
}

// 执行完毕后回调
// 根据uuid删除task
func deleteCronTask(id string) {
	lock.Lock()
	delete(cronTaskMap, id)
	lock.Unlock()
}
