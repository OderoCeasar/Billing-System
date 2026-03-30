package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/OderoCeasar/system/config"
)

type Client struct {
	config		*config.MpesaConfig
	httpClient	*http.Client
	baseURL		string
}

func NewClient(cfg *config.MpesaConfig) *Client {
	baseURL := "https://sandbox.safaricom.co.ke"
	if cfg.Environment == "production" {
		baseURL = "https://api.safaricom.co.ke"
	}

	return &Client{
		config: 	cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL: baseURL,
	}
}

// Get access token
func ( c *Client) GetAccessToken() (string, error) {
	url := c.baseURL + "/oauth/v1/generate?grant_type=client_credentials"

	req, err := http.NewRequest("GET", url, nil) 
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(
		c.config.ConsumerKey + ":" + c.config.ConsumerSecret,
	))
	req.Header.Add("Authorization", "Basic " + auth)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	return tokenResp.AccessToken, nil
}


// Initiate STK push
func (c *Client) InitiateSTKPush(request *STKPushRequest) (*STKPushResponse, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := c.baseURL + "/mpesa/stkpush/v1/processrequest"

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send STK push: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var stkResp STKPushResponse
	if err := json.Unmarshal(body, &stkResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if stkResp.ResponseCode != "0" {
		return nil, fmt.Errorf("STK Push failed: %s", stkResp.ResponseDescription)
	}

	return &stkResp, nil
}

// Generate Password
func (c *Client) GeneratePassword(timestamp string) string {
	data := c.config.ShortCode + c.config.PassKey + timestamp
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Generate Timestamp
func (c *Client) GenerateTimestamp() string {
	return time.Now().Format("20060102150405")
}


// Format phoneNumber
func (c *Client) FormatPhoneNumber(phone string) string {
	if len(phone) > 0 && phone[0] == '0' {
		return "254" + phone[1:]
	}

	if len(phone) >= 12 {
		return phone
	}

	return "254" + phone
}
