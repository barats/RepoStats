// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package utils

import (
	"crypto/sha256"
	"errors"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/remeh/sizedwaitgroup"
)

const (
	MAX_ROUTINE_NUMBER       = 20
	GITEE_SCHEDULER_INTERVAL = 6 * time.Hour
)

var (
	Version      = "1.0"
	Build        = "2204252028"
	WaitingGroup = sizedwaitgroup.SizedWaitGroup{}
)

func InitWaitingGruop() {
	WaitingGroup = sizedwaitgroup.New(MAX_ROUTINE_NUMBER)
}

func ExitOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func EmptyString(str string) bool {
	return reflect.DeepEqual("", strings.TrimSpace(str))
}

// Write file to RepoStats home path
//
func WriteRepoStatsFile(file string, data []byte) error {
	appHome, err := appHomeDir()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(appHome, "/", file), data, os.ModePerm)
}

// Read file content from RepoStats home path
//
func ReadRepoStatsFile(file string) ([]byte, error) {
	appHome, err := appHomeDir()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filepath.Join(appHome, "/", file))
}

// Get App Home Dir
//
func appHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	appPath := filepath.Join(home, ".repostats")
	return appPath, os.MkdirAll(appPath, os.ModePerm)
}

// 给定 url 解析 Gitee 的用户名和仓库名
//
func ParseGiteeRepoUrl(urlStr string) (string, string, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}

	if !strings.EqualFold(url.Host, "gitee.com") {
		return "", "", errors.New("不是 Gitee 链接")
	}

	arr := strings.Split(url.Path, "/")

	if len(arr) == 3 {
		return arr[1], arr[2], nil
	} else if len(arr) == 2 {
		return arr[1], "", nil
	} else {
		return "", "", errors.New("无法解析该链接")
	}
}

func Sha256Of(input string) ([]byte, error) {
	algorithm := sha256.New()
	_, err := algorithm.Write([]byte(strings.TrimSpace(input)))
	if err != nil {
		return nil, err
	}
	return algorithm.Sum(nil), nil
}

func Base58Encode(data []byte) string {
	return base58.Encode(data)
}
