package tomeit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func InitFirebaseApp() {
	var err error
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	firebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("initialize firebase firebaseApp failed: %v\n", err)
	}
}

func hash(token string) string {
	digestTokenByte := sha256.Sum256([]byte(token))
	digestToken := hex.EncodeToString(digestTokenByte[:])
	return digestToken
}
