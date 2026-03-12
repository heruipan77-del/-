package devicemanager

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient(cfg Config) (*Client, error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.Insecure}}
	c := &Client{
		baseURL: fmt.Sprintf("https://%s:%d/deviceManager/rest/%s", cfg.IP, cfg.Port, cfg.DeviceID),
		http: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
		},
	}
	if err := c.login(cfg.Username, cfg.Password); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) login(username, password string) error {
	body := map[string]string{"username": username, "password": password, "scope": "0"}
	var out struct {
		Data struct {
			Token string `json:"iBaseToken"`
		} `json:"data"`
	}
	if err := c.do(http.MethodPost, "sessions", body, &out); err != nil {
		return err
	}
	if out.Data.Token == "" {
		return fmt.Errorf("empty iBaseToken in login response")
	}
	c.token = out.Data.Token
	return nil
}

func (c *Client) GetStoragePools() ([]StoragePool, error) {
	var out struct {
		Data []StoragePool `json:"data"`
	}
	if err := c.do(http.MethodGet, "storagepool", nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

func (c *Client) GetLUN(lunID string) (*LUN, error) {
	var out struct {
		Data LUN `json:"data"`
	}
	if err := c.do(http.MethodGet, "lun/"+lunID, nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) do(method, path string, payload any, out any) error {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+"/"+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("iBaseToken", c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		return fmt.Errorf("request %s %s failed: status=%d body=%s", method, path, resp.StatusCode, string(msg))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
