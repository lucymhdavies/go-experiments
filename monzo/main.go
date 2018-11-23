package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/alexsasharegan/dotenv"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Monzo struct {
		ClientID     string
		ClientSecret string
		RedirectURI  string
		AccountID    string
	}
}

var cfg Config

func main() {
	// Load env vars from .env file, if present
	// Ignore errors caused by file not existing
	_ = dotenv.Load()

	if len(os.Getenv("MONZO_ACCOUNT_ID")) == 0 {
		log.Fatalf("MONZO_ACCOUNT_ID not set!")
	}
	cfg.Monzo.AccountID = os.Getenv("MONZO_ACCOUNT_ID")

	if len(os.Getenv("MONZO_CLIENT_ID")) == 0 {
		log.Fatalf("MONZO_CLIENT_ID not set!")
	}
	cfg.Monzo.ClientID = os.Getenv("MONZO_CLIENT_ID")

	if len(os.Getenv("MONZO_CLIENT_SECRET")) == 0 {
		log.Fatalf("MONZO_CLIENT_SECRET not set!")
	}
	cfg.Monzo.ClientSecret = os.Getenv("MONZO_CLIENT_SECRET")

	if len(os.Getenv("MONZO_REDIRECT_URI")) == 0 {
		log.Fatalf("MONZO_REDIRECT_URI not set!")
	}
	cfg.Monzo.RedirectURI = os.Getenv("MONZO_REDIRECT_URI")

	stateToken := uuid.New().String()

	monzoAuthURI := fmt.Sprintf("https://auth.monzo.com/?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		cfg.Monzo.ClientID,
		url.QueryEscape(cfg.Monzo.RedirectURI),
		stateToken,
	)

	e := echo.New()

	log.Infof("Go to http://localhost:8080/auth")

	// auth url, redirect to monzo login page
	e.GET("/auth", func(c echo.Context) error {
		return c.Redirect(302, monzoAuthURI)
	})

	// monzo redirects us back here...
	e.GET("/callback", func(c echo.Context) error {

		code := c.QueryParam("code")
		state := c.QueryParam("state")

		log.Infof("Code: %s", code)
		log.Infof("State: %s", state)

		//one-line post request/response...
		response, err := http.PostForm("https://api.monzo.com/oauth2/token", url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {cfg.Monzo.ClientID},
			"client_secret": {cfg.Monzo.ClientSecret},
			"redirect_uri":  {cfg.Monzo.RedirectURI},
			"code":          {code},
		})

		//okay, moving on...
		if err != nil {
			//handle postform error
			log.Fatalf("Err: %s", err)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			//handle read response error
			log.Fatalf("Err: %s", err)
		}

		fmt.Printf("\n\n%s\n\n", string(body))

		var token MonzoOAuthToken

		err = json.Unmarshal(body, &token)
		if err != nil {
			//handle unmarhsal error
			log.Fatalf("Err: %s", err)
		}

		// TODO: list accounts, iterate over them...

		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.monzo.com/balance?account_id=%s", cfg.Monzo.AccountID), nil)
		if err != nil {
			// handle err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			//handle read response error
			log.Fatalf("Err: %s", err)
		}

		fmt.Printf("\n\n%s\n\n", string(body))

		var balance MonzoBalance

		err = json.Unmarshal(body, &balance)
		if err != nil {
			//handle unmarhsal error
			log.Fatalf("Err: %s", err)
		}

		return c.String(http.StatusOK, fmt.Sprintf("Hello, World!\n\n%s %s", balance.Currency, balance.Balance))
	})

	e.Logger.Fatal(e.Start(":8080"))

}

type MonzoOAuthToken struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
}

type MonzoBalance struct {
	Balance      int    `json:"balance"`
	TotalBalance int    `json:"total_balance"`
	Currency     string `json:"currency"`
	SpendToday   int    `json:"spend_today"`
}
