package server

import (
	"github.com/NBCUniversal/gvs-ms-common/logx"
)

type Logger interface {
	Log(params ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Log(params ...interface{}) {
	logx.Infoln(params...)
}
