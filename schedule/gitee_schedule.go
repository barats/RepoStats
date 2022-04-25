package schedule

import (
	"fmt"
	"log"
	gitee_model "repostats/model/gitee"
	"repostats/network"
	gitee_storage "repostats/storage/gitee"
	"repostats/utils"
	"strings"
	"time"

	"github.com/remeh/sizedwaitgroup"
)

// 启动 Gitee 定时器
//
func StartGiteeSchedule() error {
	ticker := time.NewTicker(utils.GITEE_SCHEDULER_INTERVAL)
	for range ticker.C {
		log.Println("[RepoStats] Gitee Schedule Start.")
		if err := StarGiteeJobs(true); err != nil {
			log.Printf("[RepoStats] error while doing schedule jobs. %s", err)
		}
		log.Println("[RepoStats] Gitee Schedule Finish.")
	}
	return nil
}

// 启动 Gitee 任务
//
func StarGiteeJobs(wait bool) error {

	//检查 Grafana Token, Datasource
	grafanaToken, err := network.RetrieveGrafanaToken()
	if err != nil {
		return fmt.Errorf("无法获取 Grafana Token. \n %s", err)
	}

	datasource, err := network.RetrieveGrafanaDatasource()
	if err != nil {
		return fmt.Errorf("无法获取 Grafana 必须的数据源配置. \n %s", err)
	}

	folder, err := network.RetrieveGiteeRepostatsFolder()
	if err != nil {
		return fmt.Errorf("无法获取 Grafana 必须的 Folder . \n %s", err)
	}

	//检查 Gitee Token
	giteeToken, err := network.ValidGiteeToken()
	if err != nil {
		return fmt.Errorf("无法获取 Gitee Token. \n %s", err)
	}

	//获取所有需要爬取的仓库信息
	repos, err := gitee_storage.FindRepos()
	if err != nil {
		return fmt.Errorf("无法获取需要爬取的仓库. \n %s", err)
	}

	//抓取 Gitee 信息并存储到数据库
	wg := utils.WaitingGroup
	for _, repo := range repos {
		if !repo.EnableCrawl {
			continue
		}
		wg.Add()
		go GrabRepo(&wg, repo, giteeToken, grafanaToken, datasource, folder)
	}

	if wait {
		wg.Wait()
	}
	return nil
}

//抓取 Gitee 仓库信息并更新 Grafana 面板
//
func GrabRepo(wg *sizedwaitgroup.SizedWaitGroup, repo gitee_model.Repository,
	giteeToken network.OauthToken, grafanaToken network.GrafanaToken, grafanaDatasource network.GrafanaDatasource, grafanaFolder network.GrafanaFolder) error {

	if wg != nil {
		defer wg.Done()
	}

	log.Printf("[RepoStats] start to grab [%s]", repo.HTMLURL)
	str := strings.Split(repo.FullName, "/")
	repoInfo, err := network.GetGiteeRepo(str[0], str[1])
	if err != nil {
		log.Printf("[RepoStats] failed during GetGiteeRepo %s, %s", repo.HTMLURL, err)
		return err
	}

	repoInfo.EnableCrawl = true
	err = gitee_storage.BulkSaveRepos([]gitee_model.Repository{repoInfo}) //update newest repo info
	if err != nil {
		log.Printf("[RepoStats] failed during BulkSaveRepos %s,%s", repo.HTMLURL, err)
		return err
	}

	var users []gitee_model.User
	commits, err := network.GetGiteeCommits(str[0], str[1])
	if err != nil {
		log.Printf("[RepoStats] failed during GetGiteeCommits %s,%s", repo.HTMLURL, err)
	}
	for i := 0; i < len(commits); i++ {
		commits[i].RepoID = repo.ID
		users = append(users, commits[i].Author)
		users = append(users, commits[i].Committer)
	}
	if len(commits) > 0 {
		err = gitee_storage.BulkSaveCommits(commits)
		if err != nil {
			log.Printf("[RepoStats] failed during BulkSaveCommits %s, %s", repo.HTMLURL, err)
		}
	}

	issues, err := network.GetGiteeIssues(str[0], str[1])
	if err != nil {
		log.Printf("[RepoStats] failed during GetGiteeIssues %s, %s", repo.HTMLURL, err)
	}
	for i := 0; i < len(issues); i++ {
		issues[i].RepoID = int64(repo.ID)
		users = append(users, issues[i].User)
	}
	if len(issues) > 0 {
		err = gitee_storage.BulkSaveIssues(issues)
		if err != nil {
			log.Printf("[RepoStats] failed during BulkSaveIssues %s, %s", repo.HTMLURL, err)
		}
	}

	prs, err := network.GetGiteePullRequests(str[0], str[1])
	if err != nil {
		log.Printf("[RepoStats] failed during GetGiteePullRequests %s, %s", repo.HTMLURL, err)
	}
	for i := 0; i < len(prs); i++ {
		prs[i].RepoID = int64(repo.ID)
		users = append(users, prs[i].User)
	}
	usersNeededToSave := gitee_model.RemoveDuplicateUsers(users)
	if len(usersNeededToSave) > 0 {
		err = gitee_storage.BulkSaveUsers(usersNeededToSave)
		if err != nil {
			log.Printf("[RepoStats] failed during BulkSaveUsers %s, %s", repo.HTMLURL, err)
		}
	}
	if len(prs) > 0 {
		err = gitee_storage.BulkSavePullRequests(prs)
		if err != nil {
			log.Printf("[RepoStats] failed during BulkSavePullRequests %s,%s", repo.HTMLURL, err)
		}
	}

	stargazers, err := network.GetGiteeStargazers(str[0], str[1])
	if err != nil {
		log.Printf("[RepoStats] failed during GetGiteeStargazers %s, %s", repo.HTMLURL, err)
	}
	for i := 0; i < len(stargazers); i++ {
		stargazers[i].RepoID = int64(repo.ID)
	}
	if len(stargazers) > 0 {
		err = gitee_storage.BulkSaveStargazers(stargazers)
		if err != nil {
			log.Printf("[RepoStats] failed during BulkSaveStargazers %s, %s", repo.HTMLURL, err)
		}
	}

	err = CreateOrUpdateGrafanaRepo(repo, grafanaToken, grafanaFolder, grafanaDatasource)
	if err != nil {
		log.Printf("[RepoStats] failed during CreateGiteeRepoDashboard %s, %s", repo.HTMLURL, err)
	}

	log.Printf("[RepoStats] finish to grab [%s]", repo.HTMLURL)
	return nil
}

//创建或更新 Grafana 的项目视图面板
//
func CreateOrUpdateGrafanaRepo(repo gitee_model.Repository, token network.GrafanaToken,
	folder network.GrafanaFolder, datasource network.GrafanaDatasource) error {
	return network.CreateGiteeRepoDashboard(token, folder, datasource, repo)
}
