package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"go-api-mono/internal/pkg/errors"
	"go-api-mono/internal/pkg/logger"

	"go.uber.org/zap"
)

// Recovery 恢复中间件
func Recovery(log *logger.Logger) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID, _ := r.Context().Value("request_id").(string)
					log.Error("Panic recovered",
						zap.String("error", fmt.Sprint(err)),
						zap.String("request_id", requestID),
						zap.String("stack", string(debug.Stack())),
					)
					http.Error(w, errors.ErrInternal.Error(), http.StatusInternalServerError)
				}
			}()
			next(w, r)
		}
	}
}
