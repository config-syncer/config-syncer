package stride

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenResponse represents the expected response from the auth API.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (tr *TokenResponse) Validate() error {
	if tr.AccessToken == "" {
		return fmt.Errorf("missing access_token in auth response")
	}
	if tr.ExpiresIn == 0 {
		return fmt.Errorf("missing expires_in in auth response")
	}
	if tr.Scope == "" {
		return fmt.Errorf("missing scope in auth response")
	}
	if tr.TokenType == "" {
		return fmt.Errorf("missing token_type in auth response")
	}
	if tr.TokenType != "Bearer" {
		return fmt.Errorf("unexpected token type '%s' in auth response", tr.TokenType)
	}
	return nil
}

func (c *clientImpl) ValidateToken(tokenString string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Expect alg: HS256
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Name != "HS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// The JWT should have been signed with our client secret. Return this to the jwt parser.
		return []byte(c.clientSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}

// GetAccessToken returns a Stride access token.
func (c *clientImpl) GetAccessToken() (string, error) {
	// There's some time until expiry. Return current token.
	if c.tokenExpires.After(time.Now().Add(60 * time.Second)) {
		return c.token, nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// Another goroutine may have updated the token while we acquired the mutex.
	if c.tokenExpires.Before(time.Now().Add(60 * time.Second)) {
		tokenResponse, err := c.getAccessToken()
		if err != nil {
			return "", err
		}

		expiry := time.Duration(tokenResponse.ExpiresIn) * time.Second
		c.tokenExpires = time.Now().Add(expiry)
		c.token = tokenResponse.AccessToken
	}
	return c.token, nil
}

// getAccessToken should not be called except by GetAccessToken.
func (c *clientImpl) getAccessToken() (*TokenResponse, error) {
	url := c.authAPIBaseURL + "/oauth/token"
	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"audience":      apiAudience,
	}
	j, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auth request data: %v", err)
	}
	body := bytes.NewReader(j)
	resp, err := c.httpClient.Post(url, "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code (%d) from auth api", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}
	tr := &TokenResponse{}
	err = json.Unmarshal(b, tr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json to TokenResponse: %v", err)
	}
	if err = tr.Validate(); err != nil {
		return nil, fmt.Errorf("invalid auth token: %v", err)
	}
	return tr, nil
}
