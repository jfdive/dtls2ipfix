package logging

import (
    "os"
    "io"

    "github.com/sirupsen/logrus"
)

type LogLevel int

const (
    INFO LogLevel = iota
    WARNING
    ERROR
    DEBUG
)

func Init() {
    Log = logrus.New()
    Log.SetFormatter(&logrus.JSONFormatter{})
    LogLib = logrus.New()
    LogLib.SetFormatter(&logrus.JSONFormatter{})
}

var Log *logrus.Logger
var LogLib *logrus.Logger

func SetLevel(lvl LogLevel) {
    switch lvl {
    case INFO:
        Log.SetLevel(logrus.InfoLevel)
        // overwrite debug lvl for kafka
        LogLib.SetLevel(logrus.ErrorLevel)
    case WARNING:
        Log.SetLevel(logrus.WarnLevel)
        LogLib.SetLevel(logrus.WarnLevel)
    case DEBUG:
        Log.SetLevel(logrus.DebugLevel)
        Log.SetReportCaller(true)
        LogLib.SetLevel(logrus.DebugLevel)
        LogLib.SetReportCaller(true)
    case ERROR:
        Log.SetLevel(logrus.ErrorLevel)
        LogLib.SetLevel(logrus.ErrorLevel)
    }
}

func SetOutput(logfile string, stdout bool) error {
    if(logfile == "" && !stdout) {
        Log.SetOutput(io.Discard)
        LogLib.SetOutput(io.Discard)
        return nil
    }
    var fhList []*os.File
    if logfile != "" {
        fh, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
            Log.Errorf("Failed to create log file %s: %s", logfile, err.Error())
            return err
        }
        fhList = append(fhList, fh)
        Log.Infof("Writing to logfile %s", logfile)
    }
    if stdout {
        fhList = append(fhList, os.Stdout)
    }
    if len(fhList) > 1 {
        mw := io.MultiWriter(fhList[0], fhList[1])
        Log.SetOutput(mw)
        LogLib.SetOutput(mw)
    } else {
        Log.SetOutput(fhList[0])
        LogLib.SetOutput(fhList[0])
    }
    return nil
}
