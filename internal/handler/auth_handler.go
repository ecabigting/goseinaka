package handler

import "golang.org/x/oauth2"

type OAuthHandler struct {
	oauthConfigs map[string]*oauth2.Config
}
