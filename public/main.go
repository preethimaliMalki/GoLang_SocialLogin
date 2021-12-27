package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

var (
	oauthConFa = &oauth2.Config{
		ClientID:     "225077211913752",
		ClientSecret: "b3e33defed81f511e81f4a9b3e64981c",
		RedirectURL:  "http://www.wixis360.com/",
		Endpoint:     facebook.Endpoint,
		Scopes:       []string{"public_profile"},
	}
	oauthStateString = "thisshouldberandom"

	oauthConGo = &oauth2.Config{
		ClientID:     "391803363826-v06sb3njb224k6dq3ale44v08nckgab8.apps.googleusercontent.com",
		ClientSecret: "PRHm3LznPFCthledv923yOOh",
		RedirectURL:  "https://www.wixis360.com",
		Scopes:       []string{"email https://mail.google.com"},
		Endpoint:     google.Endpoint,
	}
)
var html = template.Must(template.ParseGlob("web/*"))

func handleMain(w http.ResponseWriter, r *http.Request) {

	html.ExecuteTemplate(w, "index.html", r)

}

func handleFacebookLogin(ww http.ResponseWriter, rr *http.Request) {
	Url, err := url.Parse(oauthConFa.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConFa.ClientID)
	parameters.Add("scope", strings.Join(oauthConFa.Scopes, ""))
	parameters.Add("redirect_uri", oauthConFa.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(ww, rr, url, http.StatusTemporaryRedirect)
}
func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	Url, err := url.Parse(oauthConGo.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConGo.ClientID)
	parameters.Add("scope", strings.Join(oauthConGo.Scopes, ""))
	parameters.Add("redirect_uri", oauthConGo.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleFacebookCallback(ww http.ResponseWriter, rr *http.Request) {
	state := rr.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(ww, rr, "/", http.StatusTemporaryRedirect)
		return
	}
	code := rr.FormValue("code")

	token, err := oauthConFa.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConFa.Exchange() failed with '%s'\n", err)
		http.Redirect(ww, rr, "/", http.StatusTemporaryRedirect)
		return
	}
	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(ww, rr, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(ww, rr, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("parseResponseBody: %s\n", string(response))

	http.Redirect(ww, rr, "/", http.StatusTemporaryRedirect)

}
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")

	token, err := oauthConGo.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConGo.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	resp, err := http.Get("https://graph.google.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("parseResponseBody: %s\n", string(response))

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login2", handleFacebookLogin)
	http.HandleFunc("/login3", handleGoogleLogin)
	http.HandleFunc("/oauth2callback2", handleFacebookCallback)
	http.HandleFunc("/oauth2callback3", handleGoogleCallback)
	fmt.Print("Started running on http://localhost:9090\n")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
