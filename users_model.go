package tomeit

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type User struct {
	id        int64
	digestUID string
}

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("UserCtx start")
		idToken := r.Header.Get("Authorization")

		ctx := r.Context()

		token, err := verifyIDToken(ctx, firebaseApp, idToken)
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

		ctx = context.WithValue(ctx, "user", user)
		log.Println("userContext end")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
