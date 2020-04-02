package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	appAddr      = "demo-client.local:5556"
	keykloakAddr = "keykloak:8080"
	clientID     = "demo-gallery"
	clientSecret = "5ca6bc50-eec1-4d57-b7f7-486cc852ee99"
	oauthState   = "foobar" //HACK: generate dynamically per request and store somewhere in session
)

var (
	oauth2Config *oauth2.Config
)

type HomePageInfo struct {
	LoginUrl    string
	ProfileUrl  string
	LogoutUrl   string
	AccessToken string
}

func main() {
	http.HandleFunc("/", root())
	http.HandleFunc("/callback", callback())
	http.HandleFunc("/cleartoken", cleartoken())

	log.Fatal(http.ListenAndServe("0.0.0.0:5556", nil))
}

func root() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if oauth2Config == nil {
			oidcProvider, err := oidc.NewProvider(
				context.Background(),
				fmt.Sprintf("http://%s/auth/realms/testrealm", keykloakAddr))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			oauth2Config = &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     oidcProvider.Endpoint(),
				RedirectURL:  fmt.Sprintf("http://%s/callback", appAddr),
				Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
			}
		}


		oauth2AccessTokenCookie, err := r.Cookie("oauth2_access_token")

		var oauth2AccessToken string = ""
		if err == nil {
			oauth2AccessToken = oauth2AccessTokenCookie.Value
		}

		clearTokenUrlEncoded := url.QueryEscape(fmt.Sprintf("http://%s/cleartoken", appAddr))
		logoutUrl := fmt.Sprintf("http://%s/auth/realms/testrealm/protocol/openid-connect/logout?redirect_uri=%s",
			keykloakAddr, clearTokenUrlEncoded)


		pageInfo := &HomePageInfo{
			LoginUrl:    oauth2Config.AuthCodeURL(oauthState),
			ProfileUrl:  fmt.Sprintf("http://%s/auth/realms/testrealm/account", keykloakAddr),
			LogoutUrl:   logoutUrl,
			AccessToken: oauth2AccessToken,
		}

		homePageTemplate := template.New("home.html")
		homePageTemplate, err = homePageTemplate.ParseFiles("web/templates/home.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = homePageTemplate.Execute(w, pageInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func callback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != oauthState {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := oauth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		newDemoCookie := &http.Cookie{
			Name:  "oauth2_access_token",
			Value: oauth2Token.AccessToken,
		}
		http.SetCookie(w, newDemoCookie)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func cleartoken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newDemoCookie := &http.Cookie{
			Name:   "oauth2_access_token",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, newDemoCookie)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
