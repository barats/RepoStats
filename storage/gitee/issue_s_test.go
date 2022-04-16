// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package gitee

import (
	gitee_model "repostats/model/gitee"
	"repostats/network"
	"repostats/utils"
	"testing"
)

func TestBulkSaveIssues(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	found, err := network.GetGiteeIssues("barat", "ohurlshortener")
	utils.ExitOnError(err)

	found2, err := network.GetGiteeIssues("openharmony", "community")
	utils.ExitOnError(err)

	type args struct {
		iss []gitee_model.Issue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Testcase barat/ohurlshortener", args: args{iss: found}, wantErr: false},
		{name: "Testcase barat/ohurlshortener again", args: args{iss: found}, wantErr: false},
		{name: "Testcase openharmony/community", args: args{iss: found2}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BulkSaveIssues(tt.args.iss); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveIssues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
