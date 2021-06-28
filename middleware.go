package tomeit

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
)

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("Authorization")

		ctx := r.Context()

		token, err := verifyIDToken(ctx, firebaseApp, idToken)
		if err != nil {
			_ = render.Render(w, r, errAuthenticate(err))
		}

		var user User

		user, err = getUserByDigestUID(hash(token.UID))
		if err != nil {
			user, err = createUser(hash(token.UID))
			if err != nil {
				_ = render.Render(w, r, errUnexpectCondition(err))
			}
		}

		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
