// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"repostats/utils"
)

const (
	GRAFANA_API_TOKEN            = "grafana_api_token.json"
	GRAFANA_DATASOURCE           = "grafana_datasource.json"
	GRAFANA_FOLDER               = "grafana_folder.json"
	REPOSTATS_FOLDER_NAME        = "RepoStats"
	REPOSTATS_HOMEDASHBOARD_NAME = "Overview"
	REPOSTATS_HOMEDASHBOARD_FILE = "grafana_home_dashboard.json"
)

type GrafanaToken struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	Host string `json:"host"`
	Port string `json:"port"`
}

// 创建 Grafana API Token
//
func CreateGrafanaApiToken(host, port, user, password string) error {
	url := fmt.Sprintf("http://%s:%s@%s:%s/api/auth/keys", user, password, host, port)
	code, rs, err := HttpPost("", url, nil, map[string]string{
		"name": "repostats_api_token",
		"role": "Admin",
	})

	if err != nil {
		return err
	}

	if code == http.StatusConflict {
		return fmt.Errorf("以存在名 repostats_api 的 token，请在 Grafana 删除该 token 再尝试")
	}

	if code != http.StatusOK {
		return fmt.Errorf("网络请求失败： %d", code)
	}

	var token GrafanaToken
	json.Unmarshal([]byte(rs), &token)
	token.Host = host
	token.Port = port

	data, _ := json.Marshal(token)
	return utils.WriteRepoStatsFile(GRAFANA_API_TOKEN, data)
}

// 从本地文件中获取 Grafana API Token
//
func RetrieveGrafanaToken() (GrafanaToken, error) {
	var grafanaToken GrafanaToken
	data, err := utils.ReadRepoStatsFile(GRAFANA_API_TOKEN)
	if err != nil {
		return grafanaToken, err
	}
	return grafanaToken, json.Unmarshal(data, &grafanaToken)
}
