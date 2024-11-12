package shared

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/vk-rv/pvx"
)

type AuthCtxType string

const authHandle AuthCtxType = "auth"

func ExtractHandle(ctx context.Context) string {
	return ctx.Value(authHandle).(string)
}

func SetHandle(ctx context.Context, handle string) context.Context {
	return context.WithValue(ctx, authHandle, handle)
}

type AuthMiddleware struct {
	pv4  *pvx.ProtoV4Local
	symK *pvx.SymKey
}

func NewAuthMiddleware(pv4 *pvx.ProtoV4Local, symK *pvx.SymKey) *AuthMiddleware {
	return &AuthMiddleware{
		pv4:  pv4,
		symK: symK,
	}
}

func (mw *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(r.Header.Get("Authorization"))
		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cc := pvx.RegisteredClaims{}

		err := mw.pv4.
			Decrypt(token, mw.symK).
			ScanClaims(&cc)
		if err != nil {
			slog.Error("failed to read token", slog.Any("error", err))

			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r.WithContext(SetHandle(r.Context(), cc.Subject)))
	})
}
