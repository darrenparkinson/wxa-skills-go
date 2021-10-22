// Copyright 2021 Darren Parkinson

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wxaskillsservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Err implements the error interface so we can have constant errors.
type Err string

func (e Err) Error() string {
	return string(e)
}

// Skill represents a skill service skill
// TODO: Make helpers to access pointer values to avoid panics
type Skill struct {
	SkillID                  *string  `json:"skill_id,omitempty"`
	DeveloperID              *string  `json:"developer_id,omitempty"`
	URL                      *string  `json:"url,omitempty"`
	Name                     *string  `json:"name,omitempty"`
	ContactEmail             *string  `json:"contact_email,omitempty"`
	Public                   *bool    `json:"public,omitempty"`
	Deleted                  *bool    `json:"deleted,omitempty"`
	CreatedAt                *string  `json:"created_at,omitempty"`
	DeletedAt                *string  `json:"deleted_at,omitempty"`
	ModifiedAt               *string  `json:"modified_at,omitempty"`
	LastActiveAt             *string  `json:"last_active_at,omitempty"`
	SuggestedInvocationNames []string `json:"suggested_invocation_names,omitempty"`
	Languages                []string `json:"languages,omitempty"`
	HomePage                 *string  `json:"home_page,omitempty"`
	Description              *string  `json:"description,omitempty"`
	Secret                   *string  `json:"secret,omitempty"`
	PublicKey                *string  `json:"public_key,omitempty"`
}

// Error Constants
const (
	ErrBadRequest    = Err("api: bad request")
	ErrUnauthorized  = Err("api: unauthorized request")
	ErrForbidden     = Err("api: forbidden")
	ErrNotFound      = Err("api: resource not found")
	ErrInternalError = Err("api: internal error")
	ErrUnknown       = Err("api: unexpected error occurred")
)

// Client is the main client for interacting with the library.  It can be created using NewClient
type Client struct {
	// BaseURL for API.  Set using NewClient or you can set directly.
	BaseURL string

	// DeveloperID is the base64 decoded developer ID
	DeveloperID string

	// Token is the personal access token for interacting with the skills service
	Token string

	//HTTP Client to use for making requests, allowing the user to supply their own if required.
	HTTPClient *http.Client
}

// NewClient is a helper function that returns an new api client given a token and developer ID.
// Optionally you can provide your own http client or use nil to use the default.  This is done to
// ensure you're aware of the decision you're making to not provide your own http client.
func NewClient(developerID, token string, client *http.Client) (*Client, error) {
	if token == "" {
		return nil, errors.New("token required")
	}
	if developerID == "" {
		return nil, errors.New("token required")
	}
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	c := &Client{
		BaseURL:    fmt.Sprintf("https://assistant.us-east-2.intelligence.webex.com/skills/api/developers/%s", developerID),
		HTTPClient: client,
		Token:      token,
	}
	return c, nil
}

// ListSkills will list all skills
func (c *Client) ListSkills(ctx context.Context) ([]Skill, error) {
	s := []Skill{}
	url := fmt.Sprintf("%s/skills", c.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return s, err
	}
	if err := c.makeRequest(ctx, req, &s); err != nil {
		return s, err
	}
	return s, nil
}

// CreateSkill will create a new skill
func (c *Client) CreateSkill(ctx context.Context, skill Skill) (*Skill, error) {
	s := &Skill{}
	url := fmt.Sprintf("%s/skills", c.BaseURL)
	payload, err := json.Marshal(skill)
	if err != nil {
		return s, err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return s, err
	}
	if err := c.makeRequest(ctx, req, s); err != nil {
		return s, err
	}
	return s, nil
}

// DeleteSkill is used delete a skill. It is required to pass the ID.
func (c *Client) DeleteSkill(ctx context.Context, id string, hardDelete bool) error {
	url := fmt.Sprintf("%s/skills/%s?HARD_DELETE=%t", c.BaseURL, id, hardDelete)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	if err := c.makeRequest(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

// makeRequest provides a single function to add common items to the request.
func (c *Client) makeRequest(ctx context.Context, req *http.Request, v interface{}) error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rc := req.WithContext(ctx)
	res, err := c.HTTPClient.Do(rc)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var apiError error
		switch res.StatusCode {
		case 400:
			apiError = ErrBadRequest
		case 401:
			apiError = ErrUnauthorized
		case 403:
			apiError = ErrForbidden
		case 404:
			apiError = ErrNotFound
		case 500:
			apiError = ErrInternalError
		default:
			apiError = ErrUnknown
		}
		return apiError
	}
	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }
