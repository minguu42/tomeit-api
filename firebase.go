package tomeit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var firebaseApp *firebase.App

func InitFirebaseApp() {
	var err error
	firebaseApp, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("initialize firebase firebaseApp failed: %v\n", err)
	}
}

func verifyIDToken(ctx context.Context, app *firebase.App, idToken string) (*auth.Token, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func hash(token string) string {
	digestTokenByte := sha256.Sum256([]byte(token))
	digestToken := hex.EncodeToString(digestTokenByte[:])
	return digestToken
}
