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
	"time"
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

type GrafanaDatasource struct {
	ID   int    `json:"id"`
	UID  string `json:"uid"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type GrafanaFolder struct {
	ID        int       `json:"id"`
	UID       string    `json:"uid"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	CreatedBy string    `json:"createdBy"`
	Created   time.Time `json:"created"`
	UpdatedBy string    `json:"updatedBy"`
	Updated   time.Time `json:"updated"`
	Version   int       `json:"version"`
}

type GrafanaPanel struct {
	ID      int `json:"id"`
	GridPos struct {
		H int `json:"h"`
		W int `json:"w"`
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"gridPos"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Targets []struct {
		Datasource   GrafanaDatasource `json:"datasource"`
		Format       string            `json:"format"`
		Group        []interface{}     `json:"group"`
		MetricColumn string            `json:"metricColumn"`
		RawQuery     bool              `json:"rawQuery"`
		RawSQL       string            `json:"rawSql"`
		RefID        string            `json:"refId"`
		Select       [][]struct {
			Params []string `json:"params"`
			Type   string   `json:"type"`
		} `json:"select"`
		Table          string `json:"table"`
		TimeColumn     string `json:"timeColumn"`
		TimeColumnType string `json:"timeColumnType"`
		Where          []struct {
			Name   string        `json:"name"`
			Params []interface{} `json:"params"`
			Type   string        `json:"type"`
		} `json:"where"`
	} `json:"targets"`
}

type GrafanaDashboard struct {
	ID      int            `json:"id"`
	UID     string         `json:"uid"`
	Refresh string         `json:"refresh"`
	Tags    []string       `json:"tags"`
	Title   string         `json:"title"`
	Panels  []GrafanaPanel `json:"panels"`
	RepoUrl string         `json:"repo_url"`
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

// 创建 Grafana 数据源
//
func CreateDatasource(token GrafanaToken, dbConfig utils.DatabaseConfigInfo) error {
	url := fmt.Sprintf("http://%s:%s/api/datasources", token.Host, token.Port)
	str := fmt.Sprintf(`{
		"name": "RepoStats_PG",
		"type": "postgres",
		"url": "%s:%d",
		"access": "proxy",
		"user": "%s",
		"database": "%s",
		"basicAuth": true,
		"basicAuthUser": "%s",
		"readOnly": true,	
		"withCredentials": false,
		"isDefault": true,	
		"secureJsonData": {
			"password": "%s",   		
			"basicAuthPassword": "%s"
		},
		"jsonData": {
				"maxOpenConns": 30,
				"postgresVersion": 906,
				"sslmode": "disable",
				"timeInterval": "30m",
				"tlsAuth": false,
				"tlsAuthWithCACert": false,
				"tlsConfigurationMethod": "file-path",
				"tlsSkipVerify": true
		}
	}`, dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.DbName, dbConfig.User, dbConfig.Password, dbConfig.Password)
	code, rs, err := HttpPost(token.Key, url, nil, str)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return fmt.Errorf("grafana datasource creation failed. status code: %d ", code)
	}

	var rawMap map[string]json.RawMessage
	json.Unmarshal([]byte(rs), &rawMap)
	return utils.WriteRepoStatsFile(GRAFANA_DATASOURCE, rawMap["datasource"])
}

// 从本地文件中获取 Grafana Datasource
//
func RetrieveGrafanaDatasource() (GrafanaDatasource, error) {
	var datasource GrafanaDatasource
	data, err := utils.ReadRepoStatsFile(GRAFANA_DATASOURCE)
	if err != nil {
		return datasource, err
	}
	return datasource, json.Unmarshal(data, &datasource)
}

// 给某个项目创建 Dashboard
//
// rType 可以是 Gitee、Github
// func CreateRepoDashboard(token GrafanaToken, repoID int64, rType string) error {
// 	token, err := RetrieveGrafanaToken()
// 	folder, err := RetrieveRepostatsFolder()
// 	datasource, err := RetrieveGrafanaDatasource()
// 	if err != nil {
// 		return err
// 	}

// 	url := fmt.Sprintf("http://%s:%s/api/dashboards/db", token.Host, token.Port)
// 	issueTypeChart := fmt.Sprintf(grafana_core.PanelRepoIssueTypeChart, datasource.UID, repoID)
// 	issueStateChart := fmt.Sprintf(grafana_core.PanelRepoIssueStateChart, datasource.UID, repoID)
// 	allPanels := strings.Join([]string{issueStateChart, issueTypeChart}, ",")
// 	data := fmt.Sprintf(grafana_core.RepoDashboard, "title", `"gitee","openharmmony/docs"`, allPanels, folder.UID, time.Now())
// 	rcode, rs, err := network.HttpPost(token.Key, url, nil, data)
// 	if err != nil {
// 		return err
// 	}

// 	if rcode != http.StatusOK {
// 		return fmt.Errorf("create dashboard for repo failed. %s", rs)
// 	}

// 	var strMap map[string]string
// 	json.Unmarshal([]byte(rs), &strMap)
// 	rcode, rrs, rerr := network.HttpGet(token.Key, fmt.Sprintf("http://%s/api/dashboards/uid/%s", graCfg.Url, strMap["uid"]), nil, nil)
// 	if rerr != nil {
// 		return rerr
// 	}

// 	if rcode != http.StatusOK {
// 		return fmt.Errorf("grafana repo dashboard creation failed. status code: %d , response : %s", rcode, rrs)
// 	}

// 	var rawMap map[string]json.RawMessage
// 	json.Unmarshal([]byte(rrs), &rawMap)

// 	return utils.WriteRepoStatsFile(fmt.Sprintf("gitee-repo-%d.json", repoID), []byte(rawMap["dashboard"]))

// 	return nil
// }

//创建 RepoStats 使用的 Folder
//
func CreateRepostatsFolder(token GrafanaToken) error {
	url := fmt.Sprintf("http://%s:%s/api/folders", token.Host, token.Port)
	data := fmt.Sprintf(`{"title":"%s"}`, REPOSTATS_FOLDER_NAME)
	code, rs, err := HttpPost(token.Key, url, nil, data)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return fmt.Errorf("grafana folder creation failed. status code: %d ", code)
	}

	return utils.WriteRepoStatsFile(GRAFANA_FOLDER, []byte(rs))
}

// 获取本地存储的 Folder 信息
func RetrieveRepostatsFolder() (GrafanaFolder, error) {
	var folder GrafanaFolder
	data, _ := utils.ReadRepoStatsFile(GRAFANA_FOLDER)
	return folder, json.Unmarshal(data, &folder)
}
