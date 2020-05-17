package logger

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// LogMgr 日志管理器
// 字段
//   writers: 日志输出器
//   datachan: 日志队列
//   wg: 同步, 等待日志打印完
type LogMgr struct {
	writers   []Writer
	datachan  chan *LogData
	wg        sync.WaitGroup
	stopPrint bool
}

func (mgr *LogMgr) run() {
	mgr.stopPrint = false
	for data := range mgr.datachan {
		for _, writer := range mgr.writers {
			writer.Write(data)
		}
		mgr.wg.Done()
	}
}

func (mgr *LogMgr) stop() {
	mgr.stopPrint = false
	mgr.wg.Wait()
	close(mgr.datachan)
	for _, writer := range mgr.writers {
		writer.Close()
	}
}

func (mgr *LogMgr) printLog(ctx context.Context, level LogLevel, callLevel int, format string, args ...interface{}) {
	if config.Level > level {
		return
	}

	if mgr.stopPrint {
		return
	}

	logTime := time.Now().Format("2006-01-02 15:04:05.000")

	callLevel++ // 加上printLog自己
	filename, funcname, line := GetCallStack(callLevel)
	logData := &LogData{
		msg:      fmt.Sprintf(format, args...),
		timeStr:  logTime,
		filename: filename,
		funcname: funcname,
		lineno:   line,
		loglevel: level,
		traceID:  GetTraceID(ctx),
		field:    getField(ctx),
	}

	select {
	case mgr.datachan <- logData:
		mgr.wg.Add(1) // 增加一个日志任务
	default:
	}
}
