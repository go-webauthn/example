package logging

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/go-webauthn/example/internal/configuration"
)

func Configure(config *configuration.Log) (logger *zap.Logger, err error) {
	if config == nil {
		config = &configuration.Log{}
	}

	var (
		cores   []zapcore.Core
		encoder zapcore.Encoder
		level   zapcore.Level
	)

	if config.File.Path != "" {
		file, err := os.OpenFile(config.File.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		fileWS := zapcore.Lock(file)

		level = config.Level
		if config.File.Level != nil {
			level = *config.File.Level
		}

		levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= level
		})

		switch config.File.Encoding {
		case "json", "":
			encoderConfig := zap.NewProductionEncoderConfig()
			encoderConfig.StacktraceKey = "stacktrace"
			encoderConfig.CallerKey = "caller"
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		case "console":
			encoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
		default:
			return nil, fmt.Errorf("unknown file encoding type %s", config.File.Encoding)
		}

		cores = append(cores, zapcore.NewCore(encoder, fileWS, levelEnabler))
	}

	if !config.Console.Disable {
		priorityHigh := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})

		priorityLow := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel && lvl >= level
		})

		level = config.Level
		if config.Console.Level != nil {
			level = *config.Console.Level
		}

		console := zapcore.Lock(os.Stdout)
		consoleErr := zapcore.Lock(os.Stderr)

		switch config.Console.Encoding {
		case "json":
			encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		case "console", "":
			encoderConfig := zap.NewDevelopmentEncoderConfig()
			encoderConfig.StacktraceKey = "stacktrace"
			encoderConfig.CallerKey = "caller"
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		default:
			return nil, fmt.Errorf("unknown console encoding type %s", config.File.Encoding)
		}

		cores = append(cores, zapcore.NewCore(encoder, consoleErr, priorityHigh))
		cores = append(cores, zapcore.NewCore(encoder, console, priorityLow))
	}

	core := zapcore.NewTee(cores...)

	stacktraceEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	logger = zap.New(core, zap.AddStacktrace(stacktraceEnabler))

	logger.Debug("configured logger", zap.String("configured-level", config.Level.String()))

	zap.ReplaceGlobals(logger)

	return logger, nil
}
