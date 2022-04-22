// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package network

import (
	"log"
	"reflect"
	"repostats/storage"
	"repostats/utils"
	"testing"
)

func testSetup(t *testing.T) {
	utils.InitConfig("../repostats.ini")
	storage.InitDatabaseService()
	log.Println("test start -->")
}

func testTeardown(t *testing.T) {
	log.Println("test done <--")
}

func TestGetRepoCommits(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	//https://gitee.com/barat/ohurlshortener 51 commits so far

	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestCase1", args: args{owner: "barat", repo: "ohurlshortener"}, want: 51, wantErr: false},
		{name: "TestCase1", args: args{owner: "barat111", repo: "ohurlshortener"}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGiteeCommits(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetCommits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrabOrgRepos(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		org string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestOpenharmony", args: args{org: "openharmony"}, want: 394, wantErr: false}, //currently has 394 repos
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGiteeOrgRepos(tt.args.org)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrgRepos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetOrgRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrabUserRepos(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestCase barat", args: args{name: "barat"}, want: 6, wantErr: false}, // I have 6 public repos
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGiteeUserRepos(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserRepos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetUserRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrabIssues(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestCase barat/ohurlshortener", args: args{owner: "barat", repo: "ohurlshortener"}, want: 3, wantErr: false}, //should be 3 at the moment
		// {name: "TestCase openharmony/community", args: args{owner: "openharmony", repo: "community"}, want: 107, wantErr: false}, //should be 107 at the moment
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGiteeIssues(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIssues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetIssues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGiteePullRequests(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestPR barat/ohurlshortener", args: args{owner: "barat", repo: "ohurlshortener"}, want: 1, wantErr: false},
		{name: "TestPR openharmony/vendor_hisilicon", args: args{owner: "openharmony", repo: "vendor_hisilicon"}, want: 448, wantErr: false}, //should be 448 at the moment
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGiteePullRequests(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGiteePullRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetGiteePullRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGiteeUserInfo(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		login string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "TestCase user", args: args{login: "barat"}, want: "巴拉迪维", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGiteeUserInfo(tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGiteeUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Name, tt.want) {
				t.Errorf("GetGiteeUserInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGiteeOrgInfo(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		login string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "TestCase user", args: args{login: "openharmony"}, want: "OpenHarmony", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGiteeOrgInfo(tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGiteeOrgInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Name, tt.want) {
				t.Errorf("GetGiteeOrgInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGiteeStargazers(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestCase barat/ohurlshortener", args: args{owner: "barat", repo: "ohurlshortener"}, want: 127, wantErr: false},                             // 127 star so far
		{name: "TestCase openharmony/communication_dsoftbus", args: args{owner: "openharmony", repo: "communication_dsoftbus"}, want: 153, wantErr: false}, // 153 star so far
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGiteeStargazers(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGiteeStargazers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetGiteeStargazers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGiteeCollaborators(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "TestCase1", args: args{owner: "barat", repo: "ohurlshortener"}, want: 1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGiteeCollaborators(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGiteeCollaborators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetGiteeCollaborators() = %v, want %v", got, tt.want)
			}
		})
	}
}
