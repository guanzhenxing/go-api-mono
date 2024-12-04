package interceptor

import (
	"fmt"
	"net/http"

	"gin-api-mono/configs"
	"gin-api-mono/internal/code"
	"gin-api-mono/internal/pkg/core"
	"gin-api-mono/internal/pkg/jwtoken"
	"gin-api-mono/internal/proposal"
)

func (i *interceptor) JWTokenAuthVerify(ctx core.Context) (sessionUserInfo proposal.SessionUserInfo, err core.BusinessError) {
	// 具体 Header 参数，可根据实际情况调整
	headerAuthorizationString := ctx.GetHeader("Authorization")
	if headerAuthorizationString == "" {
		err = core.Error(
			http.StatusUnauthorized,
			code.JWTAuthVerifyError,
			"Header 中缺少 Authorization 参数")

		return
	}

	// 验证 JWT 是否合法
	jwtClaims, jwtErr := jwtoken.New(configs.Get().JWT.Secret).Parse(headerAuthorizationString)
	if jwtErr != nil {
		err = core.Error(
			http.StatusUnauthorized,
			code.JWTAuthVerifyError,
			fmt.Sprintf("jwt token 验证失败： %s", jwtErr.Error()))

		return
	}

	sessionUserInfo = jwtClaims.SessionUserInfo

	return
}
