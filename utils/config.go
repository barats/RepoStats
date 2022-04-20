// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package utils

import (
	"gopkg.in/ini.v1"
)

var (
	DatabaseConifg  DatabaseConfigInfo
	RepoStatsConfig RepoStatsConfigInfo
	GrafanaConfig   GrafanaConfigInfo
)

type GrafanaConfigInfo struct {
	Host     string
	Port     string
	User     string
	Password string
}

type RepoStatsConfigInfo struct {
	AdminPort int
	Debug     bool
	Version   string
	Build     string
}

type DatabaseConfigInfo struct {
	Host         string
	Port         int
	User         string
	Password     string
	DbName       string
	MaxOpenConns int
	MaxIdleConn  int
}

func InitConfig(file string) (*ini.File, error) {

	cfg, err := ini.Load(file)
	if err != nil {
		return nil, nil
	}

	section := cfg.Section("postgres")
	DatabaseConifg.Host = section.Key("host").String()
	DatabaseConifg.Port = section.Key("port").MustInt()
	DatabaseConifg.MaxOpenConns = section.Key("max_open_conn").MustInt()
	DatabaseConifg.MaxIdleConn = section.Key("max_idle_conn").MustInt()
	DatabaseConifg.User = section.Key("user").String()
	DatabaseConifg.Password = section.Key("password").String()
	DatabaseConifg.DbName = section.Key("database").String()

	repostatsSection := cfg.Section("repostats")
	RepoStatsConfig.Debug = repostatsSection.Key("debug").MustBool()
	RepoStatsConfig.AdminPort = repostatsSection.Key("admin_port").MustInt()
	RepoStatsConfig.Version = repostatsSection.Key("version").String()
	RepoStatsConfig.Build = repostatsSection.Key("build").String()

	grafanaSection := cfg.Section("grafana")
	GrafanaConfig.Host = grafanaSection.Key("host").String()
	GrafanaConfig.Port = grafanaSection.Key("port").String()
	GrafanaConfig.User = grafanaSection.Key("user").String()
	GrafanaConfig.Password = grafanaSection.Key("password").String()

	return cfg, err
}
