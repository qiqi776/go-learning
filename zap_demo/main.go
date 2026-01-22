package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.OpenFile("./test.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0744)
	return zapcore.AddSync(file)
}

// func (s *lockedWriteSyncer) Sync() error {
// 	s.Lock()
// 	err := s.ws.Sync()
// 	s.Unlock()
// 	return err
// }

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(
			"Error fetching url...",
			zap.String("url", url),
			zap.Error(err),
		)
	} else {
		logger.Info(
			"Success...",
			zap.String("statusCode", resp.Status),
			zap.String("url", url),
		)
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer logger.Sync()
	simpleHttpGet("www.google.com")
	simpleHttpGet("www.baidu.com")
}