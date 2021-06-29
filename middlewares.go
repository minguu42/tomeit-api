package tomeit

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type Middleware func(http.Handler) http.Handler

func UserCtx(db dbInterface, firebaseApp firebaseAppInterface) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idToken := r.Header.Get("Authorization")

			ctx := r.Context()

			token, err := firebaseApp.verifyIDToken(ctx, idToken)
			if err != nil {
				log.Println("verifyIDToken failed:", err)
				_ = render.Render(w, r, errAuthenticate(err))
				return
			}

			var user *user

			user, err = db.getUserByDigestUID(hash(token.UID))
			if user == nil || err != nil {
				user, err = db.createUser(hash(token.UID))
				if err != nil {
					log.Println("createUser failed:", err)
					_ = render.Render(w, r, errUnexpectedEvent(err))
				}
			}

			ctx = context.WithValue(ctx, "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}