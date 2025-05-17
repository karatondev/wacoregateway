package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/faisolarifin/wacoregateway/model/constant"
	"github.com/faisolarifin/wacoregateway/provider/dailylogger"
	"github.com/faisolarifin/wacoregateway/util"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type LogType int

const (
	AppLog = iota
	MongoLog
)

type ILogger interface {
	Infof(logType LogType, format string, args ...interface{})
	Errorf(logType LogType, format string, args ...interface{})
	Debugf(logType LogType, format string, args ...interface{})
	WithFields(logType LogType, fields logrus.Fields) *logrus.Entry
}

type logrusLogger struct {
	appLog *logrus.Logger
}

type CustomFormatter struct {
	TimestampFormat string
	FieldMap        logrus.FieldMap
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+4)

	// Add timestamp
	if t, ok := data[logrus.FieldKeyTime]; ok {
		data[f.FieldMap[logrus.FieldKeyTime]] = t
	} else {
		data[f.FieldMap[logrus.FieldKeyTime]] = entry.Time.Format(f.TimestampFormat)
	}

	// Add message as JSON
	messageBytes, err := json.Marshal(entry.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	data[f.FieldMap[logrus.FieldKeyMsg]] = json.RawMessage(messageBytes)

	uniqueID := uuid.New().String()

	if reqID, ok := entry.Data[constant.ReqIDLog]; ok {
		entry.Data["uniqueId"] = reqID
	}

	// Add uniqueId and xRequestId fields
	if uniqueId, ok := entry.Data["uniqueId"]; ok {
		data["uniqueId"] = uniqueId
	} else {
		data["uniqueId"] = uniqueID
	}

	if xRequestId, ok := entry.Data[constant.ReqIDLog]; ok {
		data[constant.ReqIDLog] = xRequestId
	} else {
		data[constant.ReqIDLog] = uniqueID
	}

	fields := make(map[string]interface{})

	// Add other fields
	for k, v := range entry.Data {
		if k != "uniqueId" && k != constant.ReqIDLog {
			fields[k] = v
		}
	}

	fields[entry.Level.String()] = entry.Message

	data["message"] = fields

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON: %w", err)
	}
	return append(serialized, '\n'), nil
}

func NewLogger() *logrusLogger {
	// workingDirectory, _ := os.Getwd()
	// logDirectory := path.Join(workingDirectory, "log")
	appErrorLogFile := path.Join(util.Configuration.Logger.Dir, "error", fmt.Sprintf("%s.app.error.log", util.Configuration.Logger.FileName))

	appLog := logrus.New()
	maxAge := util.Configuration.Logger.MaxAge
	maxBackups := util.Configuration.Logger.MaxBackups
	maxSize := util.Configuration.Logger.MaxSize
	compress := util.Configuration.Logger.Compress
	localTime := util.Configuration.Logger.LocalTime

	formatter := &CustomFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	}

	appLog.SetFormatter(formatter)

	// Send logs with level higher than warning to stderr
	appLog.AddHook(&WriterHook{
		Writer: dailylogger.NewDailyRotateLogger(appErrorLogFile, maxSize, maxBackups, maxAge, localTime, compress),
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	return &logrusLogger{appLog: appLog}
}

func (l *logrusLogger) Infof(logType LogType, format string, args ...interface{}) {
	logger := l.checkType(logType)
	logger.Infof(format, args...)
}

func (l *logrusLogger) Errorf(logType LogType, format string, args ...interface{}) {
	logger := l.checkType(logType)
	logger.Errorf(format, args...)
}

func (l *logrusLogger) Debugf(logType LogType, format string, args ...interface{}) {
	logger := l.checkType(logType)
	logger.Debugf(format, args...)
}

func (l *logrusLogger) WithFields(logType LogType, fields logrus.Fields) *logrus.Entry {
	logger := l.checkType(logType)
	return logger.WithFields(fields)
}

func (l *logrusLogger) checkType(logType LogType) *logrus.Logger {
	var logger *logrus.Logger

	logger = l.appLog

	return logger
}

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func InitLogDir() {
	// workingDirectory, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }

	workingDirectory := util.Configuration.Logger.Dir

	logDirectory := path.Join(workingDirectory)
	if _, err := os.Stat(logDirectory); os.IsNotExist(err) {
		if err := util.CreateDirectory(logDirectory); err != nil {
			panic(err)
		}
	}

	errorLogDirectory := path.Join(logDirectory, "error")
	if _, err := os.Stat(errorLogDirectory); os.IsNotExist(err) {
		if err := util.CreateDirectory(errorLogDirectory); err != nil {
			panic(err)
		}
	}

}
