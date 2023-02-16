package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type DefaultHandlerResponse struct {
	Errno  int         `json:"errno" dc:"Error code"`
	Errmsg string      `json:"errmsg" dc:"Error message"`
	Data   interface{} `json:"data" dc:"Result data for certain request according API definition"`
	St     int64       `json:"st"`
}

type DefaultLog struct {
	HttpStatusCode int    `json:"http_status_code"`
	Errno          int    `json:"errno"`
	Errmsg         string `json:"errmsg"`
	ErrDetail      string `json:"err_detail"`
	Host           string `json:"host"`
	ClientIP       string `json:"client_ip"`
}

func MiddlewareHandlerResponse(r *ghttp.Request) {
	r.Middleware.Next()
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
		ctx  = r.GetCtx()
	)

	if err != nil {
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		msg = err.Error()
	} else {
		if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
			msg = http.StatusText(r.Response.Status)
			switch r.Response.Status {
			case http.StatusNotFound:
				code = gcode.CodeNotFound
			case http.StatusForbidden:
				code = gcode.CodeNotAuthorized
			default:
				code = gcode.CodeUnknown
			}
			err = gerror.NewCode(code, msg)
			r.SetError(err)
		} else {
			code = gcode.CodeOK
		}
	}

	log := &DefaultLog{
		HttpStatusCode: r.Response.Status,
		Errno:          code.Code(),
		Errmsg:         msg,
		ErrDetail:      gconv.String(code.Detail()),
		Host:           r.Host,
		ClientIP:       r.GetClientIp(),
	}
	g.Log().Info(ctx, log)

	r.Response.WriteJson(DefaultHandlerResponse{
		Errno:  code.Code(),
		Errmsg: msg,
		Data:   res,
		St:     gtime.Now().TimestampMilli(),
	})
}
