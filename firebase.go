package tomeit

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type firebaseAppInterface interface {
	verifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type FirebaseApp struct {
	*firebase.App
}

func InitFirebaseApp() *FirebaseApp {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalln("Init firebase app failed:", err)
	}
	return &FirebaseApp{app}
}

func (app *FirebaseApp) verifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verifyIDToken failed: %w", err)
	}

	return token, nil
}

type firebaseAppMock struct{}

func (app *firebaseAppMock) verifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token := auth.Token{
		UID: "someUserUID",
	}
	return &token, nil
}
