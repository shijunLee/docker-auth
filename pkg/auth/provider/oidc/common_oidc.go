package oidc

import "C"
import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/shijunLee/docker-auth/pkg/config"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"time"
)

type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	EndPoint     string
	Scopes       []string
}

type OIDC struct {
	config       *OIDCConfig
	oidcProvider *oidc.Provider
	oauth2Config *oauth2.Config
	serverConfig *config.ServerConfig
}

func (c *OIDC) Registry(handler http.HandlerFunc) {

}

// NewOIDC create new oidc client
func NewOIDC(config *OIDCConfig, serverConfig *config.ServerConfig) (*OIDC, error) {
	oidcProvider, err := oidc.NewProvider(context.TODO(), config.EndPoint)

	if err != nil {
		// handle error
		return nil, err
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: oidcProvider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: config.Scopes,
	}
	return &OIDC{config: config, oauth2Config: oauth2Config, oidcProvider: oidcProvider, serverConfig: serverConfig}, nil
}

func (c *OIDC) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	state, err := randString(16)
	if err != nil {
		http.Error(w, "Internal error for get state code", http.StatusInternalServerError)
		return
	}
	c.setCallbackCookie(w, r, "state", state)
	http.Redirect(w, r, c.oauth2Config.AuthCodeURL(state), http.StatusFound)
}

func (c *OIDC) HandleCallBack(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	oauth2Token, err := c.oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := c.oidcProvider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		OAuth2Token *oauth2.Token
		UserInfo    *oidc.UserInfo
	}{oauth2Token, userInfo}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	c.setCallbackCookie(w, r, "id_token", oauth2Token.AccessToken)
}

func (c *OIDC) setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	if c.serverConfig.Domain != "" {
		cookie.Domain = c.serverConfig.Domain
	}
	http.SetCookie(w, cookie)
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
