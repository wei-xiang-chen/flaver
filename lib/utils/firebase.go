package utils

import (
	"context"
	"errors"
	"flaver/globals"
	"net/mail"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	firebaseAuthClient     *auth.Client
	firebaseAuthClientMux  sync.Mutex
	firebaseAuthClientOnce sync.Once
)

type IFirebaseAuthUtil interface {
	ParseIdToken(idToken string) (*IdTokenInfo, error)
}

type IdTokenInfo struct {
	Email    string
	TokenUID string
}

func newFirebaseAuthClient() *auth.Client {
	var client *auth.Client
	opt := option.WithCredentialsFile(globals.GetConfig().Firebase.GetCredentialFilePath())
	if app, err := firebase.NewApp(context.Background(), nil, opt); err != nil {
		globals.GetLogger().Errorf("[FirebaseAccessServicesMultipleApp] error: %v", err)
	} else if client, err = app.Auth(context.Background()); err != nil {
		globals.GetLogger().Errorf("[FirebaseAccessServicesMultipleApp] error getting Auth clien : %v", err)
	}
	return client
}

func getFirebaseAuthClient() *auth.Client {
	firebaseAuthClientMux.Lock()
	defer firebaseAuthClientMux.Unlock()
	firebaseAuthClientOnce.Do(func() {
		if firebaseAuthClient == nil {
			firebaseAuthClient = newFirebaseAuthClient()
		}
	})

	return firebaseAuthClient
}

type FirebaseAuthUtil struct {
}

func (this *FirebaseAuthUtil) ParseIdToken(idToken string) (*IdTokenInfo, error) {
	result := &IdTokenInfo{}

	if token, err := getFirebaseAuthClient().VerifyIDTokenAndCheckRevoked(
		context.Background(),
		idToken,
	); err != nil {
		return nil, err
	} else {
		result.Email = this.getEmailFromToken(token)
		result.TokenUID = token.UID
	}

	if len(result.Email) == 0 {
		return nil, errors.New("invalid email format")
	} else {
		return result, nil
	}
}

func (this *FirebaseAuthUtil) getEmailFromToken(token *auth.Token) string {
	if token == nil {
		return ""
	} else if emailSlice, ok := token.Firebase.Identities["email"]; !ok {
		return ""
	} else if emails, ok := emailSlice.([]interface{}); !ok {
		return ""
	} else if email, ok := emails[0].(string); !ok {
		return ""
	} else if _, err := mail.ParseAddress(email); err != nil {
		return ""
	} else {
		return email
	}
}
