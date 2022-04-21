package tech

import (
	"strconv"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger

type AppMess struct {
	Code    int
	Message string
}

func String2Time(ts string) (*time.Time, error) {
	lt, err := time.Parse("2006-01-02 15:04:05", ts)
	if err != nil {
		return nil, err
	}
	return &lt, nil
}

func IntString2Time(ts string) (*time.Time, error) {
	iTime, _ := strconv.ParseInt(ts, 10, 64)
	t := time.Unix(iTime, 0).UTC()
	return &t, nil
}

func IntString2yymmdd_hhmmss(ts string) string {
	t, _ := IntString2Time(ts)
	r := t.Format("2006-01-02 15:04:05")
	return r
}

func init() {
	var err error
	Log, err = zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func LogInfo(message string) {
	Log.Info(message)
}

func LogWarn(message string) {
	Log.Warn(message)
}
