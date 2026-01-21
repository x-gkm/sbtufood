package twitter

import (
	"context"
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
)

type Auth struct {
	BearerToken    string
	ConsumerKey    string
	ConsumerSecret string
	AccessKey      string
	AccessSecret   string
}

type bearerAuth struct {
	Token string
}

func (ba bearerAuth) Add(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+ba.Token)
}

func New(ctx context.Context, auth Auth) *twitter.Client {
	httpClient := oauth1.NewClient(
		ctx,
		oauth1.NewConfig(
			auth.ConsumerKey,
			auth.ConsumerSecret,
		),
		oauth1.NewToken(
			auth.AccessKey,
			auth.AccessSecret,
		),
	)

	return &twitter.Client{
		Authorizer: bearerAuth{
			Token: auth.BearerToken,
		},
		Client: httpClient,
		Host:   "https://api.twitter.com",
	}
}
