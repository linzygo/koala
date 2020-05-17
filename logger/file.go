package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileWriter 日志写到文件
type FileWriter struct {
	curDate    string
	infoFile   *os.File
	accessFile *os.File
	warnFile   *os.File
}

// NewFileWriter 创建日志文件器
// 返回值
//   Writer: 日志输出器
func NewFileWriter() Writer {
	log := &FileWriter{}
	return log
}

// Write 实现接口Writer
func (writer *FileWriter) Write(data *LogData) {
	time, err := time.Parse("2006-01-02 15:04:05.000", data.timeStr)
	if err != nil {
		return
	}
	// 新的日期，创建新的日志文件
	logdate := time.Format("2006-01-02")
	if writer.curDate != logdate {
		writer.Close()
		writer.createDir(logdate)
		writer.curDate = logdate
	}

	logSuffix := ".log"
	file := &writer.infoFile
	if data.loglevel > LogLevelAccess {
		file = &writer.warnFile
		logSuffix = ".wf.log"
	} else if data.loglevel > LogLevelInfo {
		file = &writer.accessFile
		logSuffix = ".af.log"
	}

	// 日志文件还未创建, 先建日志文件
	if (*file) == nil {
		logtime := time.Format("15-04-05")
		*file, err = writer.createFile(logdate, logtime, logSuffix)
		if err != nil {
			return
		}
	}

	// 日志文件超过配置大小，则创建新的日志文件
	if fileInfo, err := (*file).Stat(); err == nil && fileInfo.Size() > config.FileSize {
		(*file).Close()
		logtime := time.Format("15-04-05")
		*file, err = writer.createFile(logdate, logtime, logSuffix)
		if err != nil {
			return
		}
	}

	// 日志写到文件
	(*file).WriteString(data.String())
	(*file).Sync()
}

// Close 实现接口Writer
func (writer *FileWriter) Close() {
	if writer.infoFile != nil {
		writer.infoFile.Close()
		writer.infoFile = nil
	}
	if writer.accessFile != nil {
		writer.accessFile.Close()
		writer.accessFile = nil
	}
	if writer.warnFile != nil {
		writer.warnFile.Close()
		writer.warnFile = nil
	}
}

func (writer *FileWriter) createDir(name string) (err error) {
	rootDir := filepath.FromSlash(config.LogDir)
	logDir := filepath.Join(rootDir, name)
	// 确保日志目录存在
	if err = os.MkdirAll(logDir, os.ModeDir); err != nil && !os.IsExist(err) {
		err = fmt.Errorf("init创建日志目录[%s]失败, err=%v", logDir, err)
	}

	return
}

func (writer *FileWriter) createFile(logdate, logtime, ext string) (file *os.File, err error) {
	filePath := filepath.Join(config.LogDir, logdate, config.ModuleName+"-"+logtime+ext)
	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		err = fmt.Errorf("createFile, 创建日志文件[%s]失败, %v", filePath, err)
	}
	return
}
