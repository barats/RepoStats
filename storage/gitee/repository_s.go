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
	"repostats/storage"
	"strings"
)

var repoQueryPrefix = `SELECT r.owner_id AS "owner.id", r.assigner_id AS "assigner.id", r.*  FROM gitee.repos r `

func BulkSaveRepos(repos []gitee_model.Repository) error {
	query := `INSERT INTO gitee.repos (id, full_name, human_name, path,name, url, owner_id,assigner_id, description, 
		html_url, ssh_url,forked_repo,default_branch, forks_count, stargazers_count, watchers_count,license, pushed_at, created_at, updated_at,enable_crawl)
	VALUES(:id, :full_name, :human_name, :path, :name, :url, :owner.id, :assigner.id, :description, :html_url, :ssh_url, :forked_repo,
		:default_branch, :forks_count, :stargazers_count, :watchers_count, :license, :pushed_at, :created_at, :updated_at,:enable_crawl)
	ON CONFLICT (id) DO UPDATE SET id=EXCLUDED.id,full_name=EXCLUDED.full_name,human_name=EXCLUDED.human_name,path=EXCLUDED.path,
		url=EXCLUDED.url,owner_id=EXCLUDED.owner_id,assigner_id=EXCLUDED.assigner_id, description=EXCLUDED.description,
		html_url=EXCLUDED.html_url,ssh_url=EXCLUDED.ssh_url,forked_repo=EXCLUDED.forked_repo, forks_count=EXCLUDED.forks_count,
		stargazers_count=EXCLUDED.stargazers_count, watchers_count=EXCLUDED.watchers_count,
		license=EXCLUDED.license,pushed_at=EXCLUDED.pushed_at,created_at=EXCLUDED.created_at,updated_at=EXCLUDED.updated_at,enable_crawl=EXCLUDED.enable_crawl`
	return storage.DbNamedExec(query, repos)
}

func FindRepos() ([]gitee_model.Repository, error) {
	repos := []gitee_model.Repository{}
	query := repoQueryPrefix + ` ORDER BY r.id DESC`
	err := storage.DbSelect(query, &repos)
	return repos, err
}

func FindPagedRepos(page, size int) ([]gitee_model.Repository, error) {
	if page < 1 {
		page = 1
	}
	repos := []gitee_model.Repository{}
	query := repoQueryPrefix + ` ORDER BY r.id DESC LIMIT $1 OFFSET $2`
	offset := (page - 1) * size
	return repos, storage.DbSelect(query, &repos, size, offset)
}

func FindRepoByID(repoID int) (gitee_model.Repository, error) {
	found := gitee_model.Repository{}
	query := repoQueryPrefix + ` WHERE r.id = $1`
	err := storage.DbGet(query, &found, repoID)
	return found, err
}

func DeleteRepo(repoID int) error {
	query := `DELETE FROM gitee.repos WHERE id = $1`
	return storage.DbExec(query, repoID)
}

func FindTotalReposCount() (int, error) {
	var count int
	query := `SELECT count(r.id) FROM gitee.repos r`
	return count, storage.DbGet(query, &count)
}

func FindReposCountByName(name string) (int, error) {
	var count int
	query := `SELECT count(r.id) FROM gitee.repos r WHERE lower(r."path") LIKE $1 OR lower(r."name") LIKE $2`
	return count, storage.DbGet(query, &count, "%"+strings.ToLower(name)+"%", "%"+strings.ToLower(name)+"%")
}

func FindPagedReposByName(name string, page, size int) ([]gitee_model.Repository, error) {
	if page < 1 {
		page = 1
	}
	repos := []gitee_model.Repository{}
	query := repoQueryPrefix + `  WHERE lower(r."path") LIKE $1 OR lower(r."name") LIKE $2 ORDER BY r.id DESC LIMIT $3 OFFSET $4`
	offset := (page - 1) * size
	return repos, storage.DbSelect(query, &repos, "%"+strings.ToLower(name)+"%", "%"+strings.ToLower(name)+"%", size, offset)
}
