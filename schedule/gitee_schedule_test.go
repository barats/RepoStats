package schedule

import (
	"log"
	"repostats/storage"
	"repostats/utils"
	"testing"
)

func TestStarGiteeJobs(t *testing.T) {

	testSetup(t)
	defer testTeardown(t)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "TestCase after OH", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StarGiteeJobs(true); (err != nil) != tt.wantErr {
				t.Errorf("StarGiteeJobs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func testSetup(t *testing.T) {
	utils.InitConfig("../repostats.ini")
	storage.InitDatabaseService()
	log.Println("test start -->")
}

func testTeardown(t *testing.T) {
	log.Println("test done <--")
}
