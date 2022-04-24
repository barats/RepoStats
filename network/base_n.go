// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package network

import (
	"fmt"
	"log"
	"repostats/utils"
	"time"

	"github.com/go-resty/resty/v2"
)

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
}

func HttpPost(token string, url string, pathParams map[string]string, body interface{}) (int, string, error) {
	client := NewNetworkClient(token)
	rs, err := client.R().SetPathParams(pathParams).SetBody(body).Post(url)
	return rs.StatusCode(), rs.String(), err
}

func HttpGet(token string, url string, pathParams map[string]string, queryParams map[string]string) (int, string, error) {
	client := NewNetworkClient(token)
	rs, err := client.R().SetPathParams(pathParams).SetQueryParams(queryParams).Get(url)
	return rs.StatusCode(), rs.String(), err
}

func NewNetworkClient(token string) *resty.Client {
	client := resty.New()
	client.SetRetryCount(1).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetAuthToken(token)

	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   fmt.Sprintf("RepoStats/%s build %s", utils.Version, utils.Build),
	})

	client.OnError(func(r *resty.Request, err error) {
		if err != nil {
			log.Printf("Network error: %s", err)
		}
	})

	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		return nil
	})

	return client
}
