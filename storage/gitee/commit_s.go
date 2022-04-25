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
)

const max_array_size = 60000

var commitQueryPrefix = `SELECT c.author_name AS "author.name",c.author_email AS "author.email", 
c.author_name AS "commit.author.name", c.author_email AS "commit.author.email", c.author_date AS "commit.author.date",
c.committer_name AS "committer.name",c.committer_email AS "committer.email", 
c.committer_name AS "commit.committer.name", c.committer_email AS "commit.committer.email", c.committer_date AS "commit.committer.date",
c.detail_message AS "commit.message", c.tree AS "commit.tree.sha",c.* FROM gitee.commits c `

func BulkSaveCommits(commits []gitee_model.Commit) error {
	query := `INSERT INTO gitee.commits (repo_id, sha, html_url, author_name, author_email, author_date, committer_name, 
		committer_email, committer_date, detail_message,tree)
	VALUES(:repo_id,:sha,:html_url,:commit.author.name, :commit.author.email,:commit.author.date,
		 :commit.committer.name,:commit.committer.email,:commit.committer.date,:commit.message,:commit.tree.sha)
	ON CONFLICT (repo_id,sha) DO UPDATE SET repo_id=EXCLUDED.repo_id,html_url=EXCLUDED.html_url,
	author_name=EXCLUDED.author_name,author_email=EXCLUDED.author_email,author_date=EXCLUDED.author_date,committer_name=EXCLUDED.committer_name,
	committer_email=EXCLUDED.committer_email,committer_date=EXCLUDED.committer_date,detail_message=EXCLUDED.detail_message,tree=EXCLUDED.tree`
	if len(commits) > max_array_size {
		nc := splitCommitsArray(commits)
		for i := 0; i < len(nc); i++ {
			storage.DbNamedExec(query, nc[i])
		}
		return nil
	}
	return storage.DbNamedExec(query, commits)
}

func FindTotalCommitsCount() (int, error) {
	var count int
	query := `SELECT count(c.sha) FROM gitee.commits c`
	return count, storage.DbGet(query, &count)
}

func FindPagedCommits(page, size int) ([]gitee_model.Commit, error) {
	if page < 1 {
		page = 1
	}
	commits := []gitee_model.Commit{}
	query := commitQueryPrefix + ` ORDER BY c.author_date DESC LIMIT $1 OFFSET $2`
	offset := (page - 1) * size
	return commits, storage.DbSelect(query, &commits, size, offset)
}

func FindCommits() ([]gitee_model.Commit, error) {
	found := []gitee_model.Commit{}
	query := commitQueryPrefix + ` ORDER BY c.author_date DESC`
	err := storage.DbSelect(query, &found)
	return found, err
}

func FindCommitsCountByAuthorEmail(email string) (int, error) {
	var count int
	query := `SELECT count(c.sha) FROM gitee.commits c WHERE c.author_email = $1`
	return count, storage.DbGet(query, &count, email)
}

func FindCommitsCountByCommitterEmail(email string) (int, error) {
	var count int
	query := `SELECT count(c.sha) FROM gitee.commits c WHERE c.committer_email = $1`
	return count, storage.DbGet(query, &count, email)
}

func FindPagedCommitsByAuthorEmail(email string, page, size int) ([]gitee_model.Commit, error) {
	if page < 1 {
		page = 1
	}
	commits := []gitee_model.Commit{}
	query := commitQueryPrefix + ` WHERE c.author_email = $1 ORDER BY c.author_date DESC LIMIT $2 OFFSET $3`
	offset := (page - 1) * size
	return commits, storage.DbSelect(query, &commits, email, size, offset)
}

func FindPagedCommitsByCommitterEmail(email string, page, size int) ([]gitee_model.Commit, error) {
	if page < 1 {
		page = 1
	}
	commits := []gitee_model.Commit{}
	query := commitQueryPrefix + ` WHERE c.committer_email = $1 ORDER BY c.author_date DESC LIMIT $2 OFFSET $3`
	offset := (page - 1) * size
	return commits, storage.DbSelect(query, &commits, email, size, offset)
}

func FindCommitBySha(sha string) (gitee_model.Commit, error) {
	found := gitee_model.Commit{}
	query := commitQueryPrefix + ` WHERE c.sha = $1 ORDER BY c.author_date DESC`
	err := storage.DbGet(query, &found, sha)
	return found, err
}

func FindCommitsByRepoID(repoID int) ([]gitee_model.Commit, error) {
	found := []gitee_model.Commit{}
	query := commitQueryPrefix + ` WHERE c.repo_id = $1 ORDER BY c.author_date DESC`
	err := storage.DbSelect(query, &found, repoID)
	return found, err
}

func DeleteCommitBySha(sha string) error {
	query := `DELETE FROM gitee.commits WHERE sha = $1`
	return storage.DbExec(query, sha)
}

func splitCommitsArray(slice []gitee_model.Commit) [][]gitee_model.Commit {
	var chunks [][]gitee_model.Commit
	for i := 0; i < len(slice); i += max_array_size {
		end := i + max_array_size

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
