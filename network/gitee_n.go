// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	gitee_model "repostats/model/gitee"
	"repostats/utils"
	"strconv"
	"time"
)

// 获取组织下的所有公开仓库
//
// 调用此方法之前，务必确保是组织帐号
func GetOrgRepos(org string) ([]gitee_model.Repository, error) {
	token, err := validGiteeToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/orgs/%s/repos", gitee_model.GITEE_OAUTH_V5PREFIX, org)
	page := gitee_model.GITEE_API_START_PAGE
	allRepos := []gitee_model.Repository{}
	for {
		page += 1
		code, rs, err := HttpGet(token.AccessToken, url, nil, map[string]string{
			"type":     "public",
			"page":     strconv.Itoa(page),
			"per_page": strconv.Itoa(gitee_model.GITEE_API_PAGE_SIZE),
		})

		if err != nil {
			return allRepos, err
		}

		if code != http.StatusOK {
			return allRepos, errors.New("unexpected StatusCode")
		}

		var foundRepos = []gitee_model.Repository{}
		err = json.Unmarshal([]byte(rs), &foundRepos)
		if err != nil {
			return allRepos, err
		}

		if len(foundRepos) > 0 {
			allRepos = append(allRepos, foundRepos...)
			continue
		}

		break
	} //end of for
	return allRepos, nil
}

// 获取个人用户名下的所有公开仓库
//
// 调用此方法之前，务必确保是个人帐号
func GetUserRepos(name string) ([]gitee_model.Repository, error) {
	token, err := validGiteeToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/users/%s/repos", gitee_model.GITEE_OAUTH_V5PREFIX, name)
	page := gitee_model.GITEE_API_START_PAGE
	allRepos := []gitee_model.Repository{}
	for {
		page += 1
		code, res, err := HttpGet(token.AccessToken, url, nil, map[string]string{
			"page":     strconv.Itoa(page),
			"per_page": strconv.Itoa(gitee_model.GITEE_API_PAGE_SIZE),
			"type":     "all",
		})

		if err != nil {
			return allRepos, err
		}

		if code != http.StatusOK {
			return allRepos, errors.New("unexpected StatusCode")
		}

		var foundRepos = []gitee_model.Repository{}
		err = json.Unmarshal([]byte(res), &foundRepos)
		if err != nil {
			return allRepos, err
		}

		if len(foundRepos) > 0 {
			allRepos = append(allRepos, foundRepos...)
			continue
		}

		break
	} //end of for

	return allRepos, nil
}

//获取指定仓库的 issue
//
//
func GetIssues(owner string, repo string) ([]gitee_model.Issue, error) {

	token, err := validGiteeToken()
	if err != nil {
		return nil, err
	}

	var foundIssues = []gitee_model.Issue{}
	page := gitee_model.GITEE_API_START_PAGE
	for {
		page += 1
		url := fmt.Sprintf("%s/repos/%s/%s/issues", gitee_model.GITEE_OAUTH_V5PREFIX, owner, repo)
		code, rs, err := HttpGet(token.AccessToken, url, nil, map[string]string{
			"page":     strconv.Itoa(page),
			"per_page": strconv.Itoa(gitee_model.GITEE_API_PAGE_SIZE),
			"state":    "all",
		})

		if err != nil {
			return foundIssues, err
		}

		if code != http.StatusOK {
			return foundIssues, fmt.Errorf("GrabIssue failed during network. Status Code: %d", code)
		}

		var issues = []gitee_model.Issue{}
		e := json.Unmarshal([]byte(rs), &issues)
		if e != nil {
			return foundIssues, err
		}

		if len(issues) > 0 {
			foundIssues = append(foundIssues, issues...)
			continue
		}
		break
	} //end of for
	return foundIssues, nil
}

// 从仓库中获取提交记录
//
// 从制定的 owner 和 repo 中获取全部提交
func GetCommits(owner string, repo string) ([]gitee_model.Commit, error) {
	token, err := validGiteeToken()
	if err != nil {
		return nil, err
	}
	var allCommits = []gitee_model.Commit{}
	page := gitee_model.GITEE_API_START_PAGE
	for {
		page += 1
		url := fmt.Sprintf("%s/repos/%s/%s/commits", gitee_model.GITEE_OAUTH_V5PREFIX, owner, repo)
		code, rs, err := HttpGet(token.AccessToken, url, nil, map[string]string{
			"page":     strconv.Itoa(page),
			"per_page": strconv.Itoa(gitee_model.GITEE_API_PAGE_SIZE),
		})

		if err != nil {
			return allCommits, err
		}

		if code != http.StatusOK {
			return allCommits, fmt.Errorf("GrabCommit failed during network. Status Code: %d", code)
		}

		var commits = []gitee_model.Commit{}
		e := json.Unmarshal([]byte(rs), &commits)
		if e != nil {
			log.Printf("GrabCommit Failed during json parse. %s", e)
			return allCommits, e
		}

		if len(commits) > 0 {
			allCommits = append(allCommits, commits...)
			continue
		}
		break
	} //end of for
	return allCommits, nil
}

// 获取一个可用、有效的 token
//
// 先从本地配置文件中获取 access_token ，如果该 access_token 已失效，则调用 refreshGiteeToken() 更新
func validGiteeToken() (OauthToken, error) {
	var token OauthToken
	token, err := retrieveGiteeToken()
	if err != nil {
		return token, err
	}

	if time.Now().Unix() >= (token.CreatedAt + token.ExpiresIn) {
		err := refreshGiteeToken(&token)
		if err != nil {
			return token, err
		}
	}

	return token, nil
}

// 从本地配置文件中获取 access_token
//
// 从 ~/.repostats/{gitee_token_file}.json 中获取本地配置文件中的 access_token
func retrieveGiteeToken() (OauthToken, error) {
	var giteeOauth OauthToken
	if data, err := utils.ReadRepoStatsFile(gitee_model.GITEE_TOKEN_FILE); err != nil {
		return giteeOauth, err
	} else {
		return giteeOauth, json.Unmarshal(data, &giteeOauth)
	}
}

// 更新 access_token
//
// 使用已存在的 refresh_token 更新 access_token
func refreshGiteeToken(token *OauthToken) error {
	tokenUrl := fmt.Sprintf("%s?grant_type=refresh_token&refresh_token=%s", gitee_model.GITEE_OAUTH_TOKEN_URL, token.RefreshToken)
	rc, rs, err := HttpPost(token.AccessToken, tokenUrl, nil, nil)
	if err != nil {
		return err
	}

	if rc == http.StatusOK {
		return utils.WriteRepoStatsFile(gitee_model.GITEE_TOKEN_FILE, []byte(rs))
	}

	return json.Unmarshal([]byte(rs), &token)
}
