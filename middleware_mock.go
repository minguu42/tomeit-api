package tomeit

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

func mockUserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("Authorization")
		if idToken == "" {
			_ = render.Render(w, r, authenticateErr(errors.New("authorization header is empty")))
		}

		ctx := r.Context()

		user, err := getUserByDigestUID(hash("digestUID"))
		if err != nil {
			user, err = createUser(hash("digestUID"))
			if err != nil {
				_ = render.Render(w, r, unexpectedErr(err))
			}
		}

		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
