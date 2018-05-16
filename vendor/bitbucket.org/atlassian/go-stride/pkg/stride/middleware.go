package stride

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

type RequestContext struct {
	ConversationID string
	CloudID        string
	UserID         string
}

// RequireTokenMiddleware requires a Bearer token signed with the client secret.
func RequireTokenMiddleware(c Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		auth := r.Header.Get("Authorization")
		if len(auth) >= 7 && auth[:7] == "Bearer " {
			// Requests from Atlassian/Stride backend use the Bearer format.
			token = auth[7:]
		} else if len(auth) >= 4 && auth[:4] == "JWT " {
			// Requests arising from client calls in refapp html examples use JWT format.
			token = auth[4:]
		} else if jwt := r.URL.Query().Get("jwt"); jwt != "" {
			// Requests from Stride client to baked in resouces use a jwt query param.
			token = jwt
		}
		if token != "" {
			if claims, err := c.ValidateToken(token); err == nil {
				rc := getRequestContextFromClaims(claims)
				r = r.WithContext(context.WithValue(r.Context(), contextKeyRequestContext, rc))
				next.ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Requires valid token"))
	})
}

func getRequestContextFromClaims(claims jwt.MapClaims) (rc *RequestContext) {
	defer func() {
		if rec := recover(); rec != nil {
			rc = nil
		}
	}()
	userID := claims["sub"].(string)
	jwtCtx := claims["context"].(map[string]interface{})
	rc = &RequestContext{
		ConversationID: jwtCtx["resourceId"].(string),
		CloudID:        jwtCtx["cloudId"].(string),
		UserID:         userID,
	}
	return rc
}

// GetRequestContext returns the RequestContext for a request.
func GetRequestContext(r *http.Request) *RequestContext {
	if v := r.Context().Value(contextKeyRequestContext); v != nil {
		if rc, ok := v.(*RequestContext); ok {
			return rc
		}
	}
	return nil
}
