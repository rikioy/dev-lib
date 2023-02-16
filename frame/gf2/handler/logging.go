package handler

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type DefaultLog struct {
	HttpStatusCode int    `json:"http_status_code"`
	Errno          int    `json:"errno"`
	Errmsg         string `json:"errmsg"`
	ErrDetail      string `json:"err_detail"`
	Host           string `json:"host"`
	ClientIP       string `json:"client_ip"`
}

type JsonOutputsForLogger struct {
	Time    string      `json:"time"`
	Level   string      `json:"level"`
	TraceID string      `json:"trace_id"`
	Content interface{} `json:"content"`
}

var LoggingJsonHandler glog.Handler = func(ctx context.Context, in *glog.HandlerInput) {
	var jsonForLogger JsonOutputsForLogger
	if in.Level == glog.LEVEL_INFO {
		content := &DefaultLog{}
		json.Unmarshal(gconv.Bytes(in.Content), content)
		jsonForLogger = JsonOutputsForLogger{
			Time:    in.TimeFormat,
			Level:   gstr.Trim(in.LevelFormat, "[]"),
			TraceID: in.TraceId,
			Content: content,
		}
	} else {
		jsonForLogger = JsonOutputsForLogger{
			Time:    in.TimeFormat,
			Level:   gstr.Trim(in.LevelFormat, "[]"),
			TraceID: in.TraceId,
			Content: in.Content,
		}
	}
	jsonBytes, err := json.Marshal(jsonForLogger)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	in.Buffer.Write(jsonBytes)
	in.Buffer.WriteString("\n")
	in.Next(ctx)
}
