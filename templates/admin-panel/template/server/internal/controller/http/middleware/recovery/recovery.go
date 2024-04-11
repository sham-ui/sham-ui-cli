package recovery

import (
	"net/http"
	"{{ shortName }}/internal/controller/http/response"
	"runtime"

	"github.com/go-logr/logr"
)

const (
	stackSize = 1024 * 8
)

type recovery struct {
	logger logr.Logger
}

func (rec *recovery) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				response.InternalServerError(rw, r)
				stack := make([]byte, stackSize)
				stack = stack[:runtime.Stack(stack, false)]
				rec.logger.Error(
					newHttpHandlerPanicError(err, r),
					"panic recovered",
					"stack", string(stack),
				)
			}
		}()
		next.ServeHTTP(rw, r)
	})
}

func New(logger logr.Logger) *recovery {
	return &recovery{logger: logger}
}
