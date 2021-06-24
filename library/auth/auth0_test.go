package auth

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load("/Users/sumitasok/go/src/github.com/asteriaaerospace/back-office/.env")
	if err != nil {
		log.Print(err.Error())
	}
}

func TestNewAuthenticatorExchange(t *testing.T) {
	t.Log(os.Getenv("AUTH0_ISSUER"))
	tests := []struct {
		name    string
		want    error
		wantErr bool
	}{
		{name: "testExchange", want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authenticator, err := NewAuthenticator()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAuthenticator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, err2 := authenticator.Config.Exchange(context.TODO(), "yjRQsK8gOopMKLEB")
			if !reflect.DeepEqual(err2, tt.want) {
				t.Errorf("NewAuthenticator().Config.Exchange got = %v, want %v", err2, tt.want)
			}
		})
	}
}
