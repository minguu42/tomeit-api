package tomeit

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type Middleware func(http.Handler) http.Handler

type key int

var userKey key

func UserCtx(db dbInterface, firebaseApp firebaseAppInterface) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idToken := r.Header.Get("Authorization")

			ctx := r.Context()

			token, err := firebaseApp.verifyIDToken(ctx, idToken)
			if err != nil {
				log.Println("verifyIDToken failed:", err)
				_ = render.Render(w, r, authenticationError(err))
				return
			}

			var user *User

			user, err = db.getUserByDigestUID(hash(token.UID))
			if user == nil || err != nil {
				user, err = db.createUser(hash(token.UID))
				if err != nil {
					log.Println("db.createUser failed:", err)
					_ = render.Render(w, r, internalServerError(err))
				}
			}

			ctx = context.WithValue(ctx, userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
