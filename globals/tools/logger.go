package tools

import (
	"flaver/lib/constants"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerSetting struct {
	logger *zap.SugaredLogger
	level  zapcore.Level
}

func (this LoggerSetting) getLoggerLevel() zapcore.Level {
	zapLevel := GetConfig().GetZap().GetLevel()
	switch zapLevel {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func (this LoggerSetting) getEncodeLevel() func(zapcore.Level, zapcore.PrimitiveArrayEncoder) {
	encodeLevel := GetConfig().GetZap().GetEncodeLevel()
	switch encodeLevel {
	case "LowercaseLevelEncoder":
		return zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder":
		return zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder":
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder":
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func (this LoggerSetting) getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  GetConfig().GetZap().GetStacktraceKey(),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    this.getEncodeLevel(),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(GetConfig().GetZap().GetPrefix() + "2006/01/02 - 15:04:05.000"))
		},
	}
}

func (this LoggerSetting) getWriteSyncer() (zapcore.WriteSyncer, error) {
	var syncer zapcore.WriteSyncer
	if fileWriter, err := zaprotatelogs.New(
		path.Join(GetConfig().GetZap().GetDirector(), "%Y-%m-%d.log"),
		zaprotatelogs.WithLinkName(GetConfig().GetZap().GetLinkName()),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	); err != nil {
		return nil, err
	} else {
		syncer = zapcore.AddSync(fileWriter)
	}
	if !GetConfig().GetZap().GetLogInConsole() {
		return syncer, nil
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), syncer), nil
}

func (this LoggerSetting) getEncoderCore() zapcore.Core {
	if writer, err := this.getWriteSyncer(); err != nil {
		log.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	} else {
		config := this.getEncoderConfig()
		encoder := zapcore.NewConsoleEncoder(config)
		return zapcore.NewCore(encoder, writer, this.level)
	}
}

func (this LoggerSetting) getLogger() *zap.Logger {
	encoderCore := this.getEncoderCore()
	var logger *zap.Logger
	if this.level == zap.DebugLevel || this.level == zap.ErrorLevel {
		logger = zap.New(encoderCore, zap.AddStacktrace(this.level))
	} else {
		logger = zap.New(encoderCore)
	}
	if GetConfig().GetZap().GetShowLine() {
		return logger.WithOptions(zap.AddCaller())
	}
	return logger
}

func (this *LoggerSetting) SetupLogger() {
	zaphook := func(entry zapcore.Entry) error {
		if entry.Level < zap.ErrorLevel {
			return nil
		}
		runTimeEnv := GetViper().GetString("server.env")
		if runTimeEnv == constants.EnvLocal {
			return nil
		}
		return nil
	}
	this.logger = this.getLogger().
		WithOptions(zap.Hooks(zaphook)).Sugar()
}

func (this LoggerSetting) SetupFolder() error {
	directory := GetConfig().GetZap().GetDirector()
	if folder, err := os.Stat(directory); err == nil {
		if !folder.IsDir() {
			return fmt.Errorf("%s is exists and it's not a folder", directory)
		}
		return nil
	} else if os.IsNotExist(err) {
		return os.MkdirAll(directory, os.ModePerm)
	} else {
		return err
	}
}

func (this *LoggerSetting) SetupLevel() {
	this.level = this.getLoggerLevel()
}

func (this *LoggerSetting) Setup() error {
	if err := this.SetupFolder(); err != nil {
		return err
	}
	this.SetupLevel()
	this.SetupLogger()
	return nil
}

func (this LoggerSetting) GetLogger() *zap.SugaredLogger {
	return this.logger
}

var (
	logger    *zap.SugaredLogger
	loggerMux sync.Mutex
)

func GetLogger() *zap.SugaredLogger {
	if logger != nil {
		return logger
	}
	loggerMux.Lock()
	defer loggerMux.Unlock()
	if logger != nil {
		return logger
	}
	output := LoggerSetting{}
	if err := output.Setup(); err != nil {
		log.Fatalln(err)
	}
	logger = output.GetLogger()
	return logger
}
