package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/workshopapps/pictureminer.api/pkg/jwt"
)

type JWTPayload struct {
	UserID    string
	UserName  string
	IsPremium bool
}

type middlewareFunc func(next http.Handler) http.Handler

type ContextKey string

var (
	JWTPayloadCtxKey ContextKey = "jwt-payload"
)

func PassJWTPayloadIntoContext(blockUnauthenticatedRequests bool) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, jwtPayloadRaw, err := retrieveJWTFromHTTPHeader(w, r)
			if err != nil {
				log.Fatal(err)
			}
			jwtPayload := JWTPayload{
				UserID:    interfaceToStr(jwtPayloadRaw["userId"]),
				UserName:  interfaceToStr(jwtPayloadRaw["userName"]),
				IsPremium: interfaceToBool(jwtPayloadRaw["isPremium"]),
			}
			jwtCtx := context.WithValue(r.Context(), JWTPayloadCtxKey, jwtPayload)
			next.ServeHTTP(w, r.WithContext(jwtCtx))
		})
	}
}

func retrieveJWTFromHTTPHeader(w http.ResponseWriter, r *http.Request) (string, map[string]interface{}, error) {
	authorization := r.Header.Get("Authorization")
	if len(authorization) < len("Bearer ")+1 {
		return "", nil, errors.New("no jwt in authorization header")
	}
	jwtTokenString := authorization[len("Bearer "):]
	jwtClaims, err := jwt.Decode([]byte(os.Getenv("JWT_SECRET_KEY")), jwtTokenString)
	if err != nil {
		return "", nil, err
	}
	return jwtTokenString, jwtClaims, nil
}

func GetJWTPayloadFromContext(ctx context.Context) (payload JWTPayload, err error) {
	ctxPayload := ctx.Value(JWTPayloadCtxKey)
	if ctxPayload == nil {
		return JWTPayload{}, errors.New("no jwt payload in context")
	}
	return ctxPayload.(JWTPayload), nil
}

func interfaceToStr(str interface{}) string {
	switch str.(type) {
	case string:
		return str.(string)

	default:
		return ""
	}
}

func interfaceToBool(b interface{}) bool {
	switch b.(type) {
	case bool:
		return b.(bool)

	default:
		return false
	}
}
