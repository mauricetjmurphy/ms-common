package xray

import (
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/mauricetjmurphy/ms-common/logx"
)

func New(opts ...Option) error {
	cf := &config{
		daemonAddr: "127.0.0.1:2000",
		version:    "1.2.3",
	}
	for _, o := range opts {
		o(cf)
	}
	if !cf.enable {
		return nil
	}
	xray.SetLogger(LogrusXray{})
	return xray.Configure(xray.Config{
		DaemonAddr:     "127.0.0.1:2000",
		ServiceVersion: "1.2.3",
	})
}

// LogrusXray presents the customized implementation xraylog.Logger
type LogrusXray struct {
}

func (LogrusXray) Log(level xraylog.LogLevel, msg fmt.Stringer) {
	if msg == nil {
		return
	}
	logs := msg.String()
	switch level {
	case xraylog.LogLevelInfo:
		logx.Infof(logs)
	case xraylog.LogLevelDebug:
		logx.Debugf(logs)
	case xraylog.LogLevelWarn:
		logx.Warnf(logs)
	case xraylog.LogLevelError:
		logx.Errorf(logs)
	default:
		logx.Errorf(logs)
	}
}
