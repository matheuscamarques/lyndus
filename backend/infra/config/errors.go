package config

import (
	"bitbucket.org/lyndus/backend/infra/logger"
	"go.uber.org/zap"
	"net/http"
	"runtime"
)

type Err struct {
	HttpCode    int
	ResponseErr ResponseErr
}

//ResponseErr responde de erro
type ResponseErr struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // api response status code

	StatusText string `json:"status"`          // bs_user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func ResponsePerErr(w http.ResponseWriter, err error, code int) {
	if code == INTERNALSERVERERROR {
		_, file, line, _ := runtime.Caller(1)
		strErr := zap.String("err", err.Error())
		fieldFile := zap.String("file", file)
		fieldLine := zap.Int("line", line)

		logger.Error("[ERROR]",
			strErr,
			fieldFile,
			fieldLine,
		)
	} else if err != nil {
		strErr := zap.String("err", err.Error())
		logger.Warning(errMap[code].ResponseErr.Message, strErr)
	} else {
		logger.Info(errMap[code].ResponseErr.Message)
	}

	JSONResponse(errMap[code].ResponseErr, errMap[code].HttpCode, w)
	return
}
