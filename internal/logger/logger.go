package logger

import (
    "go.uber.org/zap"
)

var log *zap.SugaredLogger

func Init(level string){
    if log != nil{
        return
    }

    var cfg zap.Config

    if level == "dev"{
        cfg = zap.NewDevelopmentConfig()
    } else {
        cfg = zap.NewProductionConfig()
    }
    logger, err := cfg.Build()

    if err != nil{
        panic(err)
    }

    log = logger.Sugar()
}

func GetLogger() *zap.SugaredLogger{
    if log == nil{
        Init("debug")
    }

    return log
}

func Close(){
    if log != nil{
        log.Sync()
    }
}
