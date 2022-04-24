// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package network

import (
	gitee_model "repostats/model/gitee"
	"repostats/utils"
	"testing"
)

func TestCreateDatasource(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	token, err := RetrieveGrafanaToken()
	utils.ExitOnError(err)

	type args struct {
		token    GrafanaToken
		dbConfig utils.DatabaseConfigInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase Create Grafana Datasource", args: args{token: token, dbConfig: utils.DatabaseConifg}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateDatasource(tt.args.token, tt.args.dbConfig); (err != nil) != tt.wantErr {
				t.Errorf("CreateDatasource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateRepostatsFolder(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	token, err := RetrieveGrafanaToken()
	utils.ExitOnError(err)

	type args struct {
		token GrafanaToken
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase Create Repostats Folder", args: args{token: token}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateGiteeRepostatsFolder(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("CreateRepostatsFolder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateHomeDashboard(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	token, err := RetrieveGrafanaToken()
	utils.ExitOnError(err)

	datasource, err := RetrieveGrafanaDatasource()
	utils.ExitOnError(err)

	folder, err := RetrieveGiteeRepostatsFolder()
	utils.ExitOnError(err)

	type args struct {
		token      GrafanaToken
		folder     GrafanaFolder
		datasource GrafanaDatasource
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase", args: args{token: token, folder: folder, datasource: datasource}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateGiteeHomeDashboard(tt.args.token, tt.args.folder, tt.args.datasource); (err != nil) != tt.wantErr {
				t.Errorf("CreateHomeDashboard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateRepoDashboard(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	token, err := RetrieveGrafanaToken()
	utils.ExitOnError(err)

	datasource, err := RetrieveGrafanaDatasource()
	utils.ExitOnError(err)

	folder, err := RetrieveGiteeRepostatsFolder()
	utils.ExitOnError(err)

	repos, err := GetGiteeUserRepos("barat")
	// repos, err := GetGiteeOrgRepos("openharmony")
	utils.ExitOnError(err)

	type args struct {
		token      GrafanaToken
		folder     GrafanaFolder
		datasource GrafanaDatasource
		repo       gitee_model.Repository
	}

	type MyTest struct {
		name    string
		args    args
		wantErr bool
	}

	var tests []MyTest
	for _, r := range repos {
		tests = append(tests, MyTest{name: "TestCase", args: args{token: token, folder: folder, datasource: datasource, repo: r}, wantErr: false})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateGiteeRepoDashboard(tt.args.token, tt.args.folder, tt.args.datasource, tt.args.repo); (err != nil) != tt.wantErr {
				t.Errorf("CreateRepoDashboard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
