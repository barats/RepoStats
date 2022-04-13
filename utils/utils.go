// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package utils

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

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
