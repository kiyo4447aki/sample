package notifications

import (
	"context"
	"errors"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)


type FCM struct {
	Client *messaging.Client
}

func NewFCM(ctx context.Context, credPath string)(*FCM, error){
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credPath))
	if err != nil {
		return nil, errors.New("Firebase Appの初期化に失敗しました" + err.Error())
	}

	mc, err := app.Messaging(ctx)
	if err != nil {
		return nil, errors.New("FCMの初期化に失敗しました" + err.Error())
	}

	return &FCM{Client: mc}, nil
}


