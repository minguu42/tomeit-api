package tomeit

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
)

type User struct {
	id        int64
	digestUID string
}

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client, err := firebaseApp.Auth(context.Background())
		if err != nil {
			_ = render.Render(w, r, authenticateErr(err))
		}

		idToken := r.Header.Get("Authorization")

		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			_ = render.Render(w, r, authenticateErr(err))
		}

		var user User

		user, err = getUserByDigestUID(hash(token.UID))
		if err != nil {
			user, err = createUser(hash(token.UID))
			if err != nil {
				_ = render.Render(w, r, unexpectedErr(err))
			}
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
