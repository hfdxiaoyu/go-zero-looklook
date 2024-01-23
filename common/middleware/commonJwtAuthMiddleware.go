package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/handler"
)

// CommonJwtAuthMiddleware : with jwt on the verification, no jwt on the verification
// 是否启用jwt验证
type CommonJwtAuthMiddleware struct {
	secret string
}

// 函数“NewCommonJwtAuthMiddleware”返回带有提供的秘密的“CommonJwtAuthMiddleware”结构的新实例。
func NewCommonJwtAuthMiddleware(secret string) *CommonJwtAuthMiddleware {
	return &CommonJwtAuthMiddleware{
		secret: secret,
	}
}

// “Handle”函数是“CommonJwtAuthMiddleware”结构的一个方法。它采用“http.HandlerFunc”作为参数并返回“http.HandlerFunc”。
func (m *CommonJwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get("Authorization")) > 0 {
			//has jwt Authorization
			authHandler := handler.Authorize(m.secret)
			authHandler(next).ServeHTTP(w, r)
			return
		} else {
			//no jwt Authorization
			next(w, r)
		}
	}
}
