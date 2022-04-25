package storage

import (
	"reflect"
	"repostats/utils"
	"testing"
)

func TestNewAdmin(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		account  string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCase1", args: args{account: "repostats", password: "-2aDzm=0(ln_9^1"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewAdmin(tt.args.account, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("NewAdmin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func testSetup(t *testing.T) {
	_, err := utils.InitConfig("../repostats.ini")
	if err != nil {
		t.Error(err)
		utils.ExitOnError(err)
	}

	_, err = InitDatabaseService()
	if err != nil {
		t.Error(err)
		utils.ExitOnError(err)
	}
}

func testTeardown(t *testing.T) {
	DbClose()
}

func TestFindAdminByAccount(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	type args struct {
		account string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "TestCase 2", args: args{account: "repostats"}, want: "EZ2zQjC3fqbkvtggy9p2YaJiLwx1kKPTJxvqVzowtx6t", wantErr: false},
		// {name: "TestCase 2", args: args{account: "repostats1"}, want: "ss", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindAdminByAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAdminByAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Password, tt.want) {
				t.Errorf("FindAdminByAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
