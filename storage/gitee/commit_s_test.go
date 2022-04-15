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
	"repostats/storage"
	"repostats/utils"
	"testing"
)

func testSetup(t *testing.T) {
	_, err := utils.InitConfig("../../repostats.ini")
	if err != nil {
		t.Error(err)
		utils.ExitOnError(err)
	}

	_, err = storage.InitDatabaseService()
	if err != nil {
		t.Error(err)
		utils.ExitOnError(err)
	}
}

func testTeardown(t *testing.T) {
	storage.DbClose()
}

func TestBulkSaveCommits(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	found, err := network.GetCommits("openharmony", "community")
	utils.ExitOnError(err)

	type args struct {
		commits []gitee_model.Commit
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{name: "TestCase1", args: args{commits: found}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BulkSaveCommits(tt.args.commits); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveCommits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
