package logger

import (
    "go.uber.org/zap"
)

var log *zap.SugaredLogger

func Init(){
    l, _ := zap.NewDevelopment()
    log = l.Sugar()
}

func GetLogger() *zap.SugaredLogger{
    if log == nil{
        Init()
    }

    return log
}

func Close(){
    if log != nil{
        log.Sync()
    }
}
