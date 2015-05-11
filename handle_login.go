package main

import (
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	goauth2 "google.golang.org/api/oauth2/v2"
)

const (
	// random string for oauth2 API calls to protect against CSRF
	oauthSecretPrefix = "34132083213-"
)

var (
	googleEndpoint = oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://accounts.google.com/o/oauth2/token",
	}

	oauthGoogleConf = &oauth2.Config{
		ClientID:     "886450315285-tdimk2d8d0ap8oj693tt416b8s3rrl40.apps.googleusercontent.com",
		ClientSecret: "oQ89rf_D0zfm3BphaKREzCwQ",
		Scopes:       []string{goauth2.UserinfoProfileScope, goauth2.UserinfoEmailScope},
		Endpoint:     googleEndpoint,
	}
)

func getMyHost(r *http.Request) string {
	// on production we force https, but it's done on nginx level, so we have to
	// hardcode the scheme
	scheme := "https"
	if options.IsLocal {
		scheme = "http"
	}
	res := scheme + "://" + r.Host
	return res
}

// googleoauth2cb?redir={redirect}
func handleOauthGoogleCallback(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	LogInfof("url: %s\n", r.URL)
	state := r.FormValue("state")
	if !strings.HasPrefix(state, oauthSecretPrefix) {
		LogErrorf("invalid oauth state, expected '%s'*, got '%s'\n", oauthSecretPrefix, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	redir := state[len(oauthSecretPrefix):]
	if redir == "" {
		LogErrorf("Missing 'redir' arg for /googleoauth2cb\n")
		redir = "/"
	}

	code := r.FormValue("code")
	token, err := oauthGoogleConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		LogErrorf("oauthGoogleConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthGoogleConf.Client(oauth2.NoContext, token)

	service, err := goauth2.New(oauthClient)
	if err != nil {
		LogErrorf("goauth2.New() failed with %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	call := service.Userinfo.Get()
	userInfo, err := call.Do()
	if err != nil {
		LogErrorf("call.Do() failed with %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//LogInfof("Logged in as Google user: %s\n", userInfo.Email)
	fullName := userInfo.Name

	// also might be useful:
	// Picture
	dbUser, err := dbGetOrCreateUser(userInfo.Email, fullName)
	if err != nil {
		LogErrorf("dbGetOrCreateUser('%s', '%s') failed with '%s'\n", userInfo.Email, fullName, err)
		// TODO: show error to the user
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	LogInfof("created user %d with email '%s'\n", dbUser.ID, dbUser.Email)
	ctx.Cookie.UserID = dbUser.ID
	ctx.Cookie.IsLoggedIn = true
	setCookie(w, ctx.Cookie)
	// Maybe: dbUserSetGoogleOauth(user, tokenCredJson)
	http.Redirect(w, r, redir, http.StatusTemporaryRedirect)
}

// /logingoogle?redir=${redirect}
func handleLoginGoogle(w http.ResponseWriter, r *http.Request) {
	redir := strings.TrimSpace(r.FormValue("redir"))
	if redir == "" {
		httpErrorf(w, "Missing 'redir' arg for /logingoogle")
		return
	}
	cb := getMyHost(r) + "/googleoauth2cb"
	oauthCopy := oauthGoogleConf
	oauthCopy.RedirectURL = cb
	// oauth2 package has a way to add additional args to url (SetAuthURLParam)
	// but google doesn't seem to send them back to callback url, so I encode
	// redir inside secret
	uri := oauthCopy.AuthCodeURL(oauthSecretPrefix+redir, oauth2.AccessTypeOnline)
	http.Redirect(w, r, uri, http.StatusTemporaryRedirect)
}

// url: GET /logout?redir=${redirect}
func handleLogout(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	removeCurrentUserConnectionInfo(ctx.Cookie.UserID)
	redir := strings.TrimSpace(r.FormValue("redir"))
	if redir == "" {
		LogErrorf("Missing 'redir' arg for /logout\n")
		redir = "/"
	}
	LogInfof("redir: '%s'\n", redir)
	ctx.Cookie.IsLoggedIn = false
	setCookie(w, ctx.Cookie)
	http.Redirect(w, r, redir, http.StatusTemporaryRedirect)
}
