package qlang

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"log/slog"

	plog "github.com/phuslu/log"
	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v2/qcoll"
	"github.com/qiangyt/go-comm/v2/qerr"
	"github.com/qiangyt/go-event"
	eventloggers "github.com/qiangyt/go-event/loggers/phuslu"
	"go.uber.org/atomic"
	"gopkg.in/natefinch/lumberjack.v2"
)

func LogMapper(logger *plog.Entry, key string, mapper qcoll.ToMap) *plog.Entry {
	return LogMap(logger, key, mapper.ToMap())
}

func LogMap(logger *plog.Entry, key string, m map[string]any) *plog.Entry {
	dict := plog.NewContext(nil)
	dict.Fields(m)

	return logger.Dict(key, dict.Value())
}

type LoggerT struct {
	plog.Logger
	parent           Logger
	lumberjackLogger *lumberjack.Logger
}

type (
	Logger     = *LoggerT
	LogEntry   = *plog.Entry
	LogContext = LogEntry
)

var TraceId atomic.Int64

func NewLogContext(generateNewTraceId bool) LogContext {
	r := plog.NewContext(nil)
	if generateNewTraceId {
		r.Int64("traceId", TraceId.Add(1))
	}
	return r
}

// see lumberjack.Logger
type LoggerConfigT struct {
	MaxSize    int  `json:"max_size" yaml:"maxsize"`
	MaxAge     int  `json:"max_age" yaml:"maxage"`
	MaxBackups int  `json:"max_backups" yaml:"maxbackups"`
	LocalTime  bool `json:"local_time" yaml:"localtime"`
	Compress   bool `json:"compress" yaml:"compress"`
}

type LoggerConfig = *LoggerConfigT

func (me Logger) NewSubLogger(lctx LogContext) Logger {
	r := *me
	r.parent = me
	if lctx != nil {
		r.Context = lctx.Value()
	}
	return &r
}

func (me Logger) Parent() Logger {
	return me.parent
}

func (me Logger) Close() {
	if me.lumberjackLogger != nil {
		me.lumberjackLogger.Close()
	}
}

func (me Logger) Error(err any) LogEntry {
	r := me.Logger.Error()
	eventloggers.PhusluMarshalAnyError(r, err)
	return r
}

func NewLoggerP(console io.Writer, config LoggerConfig, fileName string) Logger {
	r, err := NewLogger(console, config, fileName)
	if err != nil {
		panic(qerr.NewSystemError("create logger", err))
	}
	return r
}

// / verbose: log to console if true
func NewLogger(console io.Writer, config LoggerConfig, fileName string) (Logger, error) {
	logD := filepath.Dir(fileName)
	if err := os.MkdirAll(logD, os.ModePerm); err != nil {
		return nil, errors.Wrapf(err, "create directory: %s", logD)
	}

	// we use lumberjack instead of phuslu/log/FileWritter as said by phuslu/log that:
	// 	"FileWriter creates a symlink to the current logging file, it requires administrator privileges on Windows."
	//  'administrator privileges' is not acceptable for our scenario
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		LocalTime:  config.LocalTime,
		Compress:   config.Compress,
	}
	fileW := &plog.IOWriter{
		Writer: lumberjackLogger,
	}

	writers := plog.MultiEntryWriter{fileW}

	if console != nil {
		consoleW := &plog.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    false,
			EndWithMessage: true,
			Writer:         console,
		}

		writers = append(writers, consoleW)
	}

	return &LoggerT{
		Logger: plog.Logger{
			Level:  plog.InfoLevel,
			Caller: 3,
			Writer: &writers,
		},
		parent:           nil,
		lumberjackLogger: lumberjackLogger,
	}, nil
}

func NewDiscardLogger() Logger {
	return &LoggerT{
		Logger: plog.Logger{
			Level: plog.FatalLevel,
			Writer: plog.IOWriter{
				Writer: io.Discard, // log is off by default
			},
		},
		parent: nil,
	}
}

// IsDiscardLogger 检查是否为 discard logger
func IsDiscardLogger(logger Logger) bool {
	if logger == nil {
		return false
	}
	w := logger.Logger.Writer
	if iw, ok := w.(plog.IOWriter); ok {
		return iw.Writer == io.Discard
	}
	return false
}

type EventLoggerT struct {
	target Logger
}

type EventLogger = *EventLoggerT

func NewEventLogger(target Logger) EventLogger {
	return &EventLoggerT{
		target: target,
	}
}

func (me EventLogger) Target() Logger { return me.target }

func (me EventLogger) LogDebug(enm event.LogEnum, hub string, topic string, lsner string) {
	me.target.Debug().Str("hub", hub).Str("topic", topic).Str("listener", lsner).Msg(enm.String())
}

func (me EventLogger) LogInfo(enm event.LogEnum, hub string, topic string, lsner string) {
	me.target.Info().Str("hub", hub).Str("topic", topic).Str("listener", lsner).Msg(enm.String())
}

func (me EventLogger) LogError(enm event.LogEnum, hub string, topic string, lsner string, err any) {
	entry := me.target.Error(err).Str("hub", hub).Str("topic", topic).Str("listener", lsner)
	entry.Msg(enm.String())
}

func (me EventLogger) LogEventDebug(enm event.LogEnum, lsner string, evnt event.Event) {
	me.target.Debug().Object("event", evnt).Str("listener", lsner).Msg(enm.String())
}

func (me EventLogger) LogEventInfo(enm event.LogEnum, lsner string, evnt event.Event) {
	me.target.Info().Object("event", evnt).Str("listener", lsner).Msg(enm.String())
}

func (me EventLogger) LogEventError(enm event.LogEnum, lsner string, evnt event.Event, err any) {
	entry := me.target.Error(err).Object("event", evnt).Str("listener", lsner)
	entry.Msg(enm.String())
}

// SlogLoggerT 将 qlang.Logger 包装为 slog.Logger
type SlogLoggerT struct {
	logger Logger
	attrs  []slog.Attr
	group  string
}

type SlogLogger = *SlogLoggerT

// NewSlogLogger 创建 slog.Logger
func NewSlogLogger(logger Logger) *SlogLoggerT {
	return &SlogLoggerT{logger: logger}
}

func (me SlogLogger) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelInfo
}

func (me SlogLogger) Handle(_ context.Context, r slog.Record) error {
	var entry *plog.Entry
	if r.Level == slog.LevelError {
		entry = me.logger.Logger.Error()
	} else if r.Level == slog.LevelWarn {
		entry = me.logger.Logger.Warn()
	} else if r.Level == slog.LevelDebug {
		entry = me.logger.Logger.Debug()
	} else {
		entry = me.logger.Logger.Info()
	}

	// 添加时间
	entry = entry.Time("time", r.Time)

	// 添加级别
	entry = entry.Str("level", r.Level.String())

	// 添加消息
	entry = entry.Str("msg", r.Message)

	// 添加group
	if me.group != "" {
		entry = entry.Str("group", me.group)
	}

	// 添加保存的属性
	for _, attr := range me.attrs {
		entry = me.addAttr(entry, attr)
	}

	// 添加record的属性
	r.Attrs(func(a slog.Attr) bool {
		entry = me.addAttr(entry, a)
		return true
	})

	entry.Msg("")
	return nil
}

func (me SlogLogger) addAttr(entry *plog.Entry, a slog.Attr) *plog.Entry {
	switch v := a.Value.Any().(type) {
	case string:
		return entry.Str(a.Key, v)
	case int:
		return entry.Int64(a.Key, int64(v))
	case int64:
		return entry.Int64(a.Key, v)
	case bool:
		return entry.Bool(a.Key, v)
	case time.Duration:
		return entry.Str(a.Key, v.String())
	case time.Time:
		return entry.Time(a.Key, v)
	default:
		return entry.Str(a.Key, fmt.Sprintf("%+v", v))
	}
}

func (me SlogLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	newLogger := &SlogLoggerT{
		logger: me.logger,
		attrs:  append([]slog.Attr{}, me.attrs...),
		group:  me.group,
	}
	newLogger.attrs = append(newLogger.attrs, attrs...)
	return newLogger
}

func (me SlogLogger) WithGroup(name string) slog.Handler {
	newLogger := &SlogLoggerT{
		logger: me.logger,
		attrs:  me.attrs,
		group:  name,
	}
	return newLogger
}

// ToSlogLogger 将 qlang.Logger 转换为 *slog.Logger
// 使用 phuslu/log 的 Slog() 方法创建
func ToSlogLogger(logger Logger) *slog.Logger {
	if logger == nil {
		return nil
	}
	return logger.Logger.Slog()
}
