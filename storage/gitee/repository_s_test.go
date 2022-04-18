// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package gitee

import (
	"reflect"
	gitee_model "repostats/model/gitee"
	"repostats/network"
	"repostats/utils"
	"testing"
)

func TestBulkSaveRepos(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	found1, err := network.GetGiteeUserRepos("barat")
	utils.ExitOnError(err)

	found2, err := network.GetGiteeOrgRepos("openharmony")
	utils.ExitOnError(err)

	type args struct {
		repos []gitee_model.Repository
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase barat", args: args{found1}, wantErr: false},
		{name: "TestCase openharmony", args: args{found2}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BulkSaveRepos(tt.args.repos); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveRepos() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindRepos(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{name: "TestCase FindRepos()", want: 399, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindRepos()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindRepos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("FindRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}
