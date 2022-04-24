-- Connect to database repostats
\c repostats

-- Create schema and define views
CREATE SCHEMA gitee;

-- Repos
CREATE TABLE gitee.repos ( 
	id INT8 NOT NULL,
	full_name VARCHAR(1000),
	human_name VARCHAR(1000),
	path VARCHAR(1000),
	name VARCHAR(1000),
	url VARCHAR(1000),
	owner_id INT8,
	assigner_id INT8,
	description VARCHAR(1000),
	html_url VARCHAR(2000),
	ssh_url VARCHAR(2000),
	forked_repo BOOLEAN,
	default_branch VARCHAR(1000),		
	forks_count INT,
	stargazers_count INT,
	watchers_count INT,
	license VARCHAR(1000),
	pushed_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE,
	enable_crawl BOOLEAN NOT NULL DEFAULT true,
  CONSTRAINT uni_gitee_repo_id UNIQUE (id)
);

-- Commits
CREATE TABLE gitee.commits (
	sha VARCHAR(80) NOT NULL,
	repo_id INT8 NOT NULL,	
	html_url VARCHAR(2000),	
	author_name VARCHAR(500),
	author_email VARCHAR(500),
	author_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	committer_name VARCHAR(200),
	committer_email VARCHAR(200),
	committer_date TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT NOW(),
	detail_message TEXT,
	tree VARCHAR(80),	
	CONSTRAINT uni_gitee_commits UNIQUE (sha,repo_id)
);

-- Issues
CREATE TABLE gitee.issues ( 
	id INT8 NOT NULL,		
	"user_id" INT8 NOT NULL,
	repo_id INT8 NOT NULL,
	html_url VARCHAR(2000),
	"number" VARCHAR(100),
	"state" VARCHAR(100),
	title VARCHAR(1000),	
	finished_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE,	
	plan_started_at TIMESTAMP WITH TIME ZONE,			
	comments INT, 
	priority INT, 
	issue_type VARCHAR(100),
	issue_state VARCHAR(100),
	security_hole BOOLEAN,
	CONSTRAINT uni_gitee_issue_id UNIQUE (id,repo_id)
);

-- Pull requests
CREATE TABLE gitee.pull_requests ( 
	id INT8 NOT NULL,
	repo_id INT8 NOT NULL,
	"user_id" INT8 NOT NULL,		
	html_url VARCHAR(2000),
	diff_url VARCHAR(2000),
	patch_url VARCHAR(2000),
	"number" INT8,
	"state" VARCHAR(100),
	created_at TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE,
	closed_at TIMESTAMP WITH TIME ZONE,
	merged_at TIMESTAMP WITH TIME ZONE,
	mergeable BOOLEAN, 
	can_merge_check BOOLEAN,
	title varchar(1000),
	head_label VARCHAR(100),
	head_ref VARCHAR(100),
	head_sha VARCHAR(100),
	head_user_id INT8,	
	head_repo_id INT8,		
	CONSTRAINT uni_gitee_prs UNIQUE (id,repo_id)
);

-- Users
CREATE TABLE gitee.users ( 
	id INT8 NOT NULL,
	login VARCHAR(500),
	"name" VARCHAR(1000),
	avatar_url VARCHAR(1000),
	html_url VARCHAR(1000),
	remark VARCHAR(1000), 
	"type" VARCHAR(500),
	email VARCHAR(1000),
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	CONSTRAINT uni_gitee_users_id UNIQUE (id)
);

-- Stargazers
CREATE TABLE gitee.stargazers ( 
	user_id INT8 NOT NULL,
	repo_id INT8 NOT NULL,
	star_at TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT NOW(),
	CONSTRAINT uni_gitee_stargazers UNIQUE(user_id,repo_id)
);

-- Collaborators
CREATE TABLE gitee.collaborators (
	 "user_id" INT8 NOT NULL, 
	 repo_id INT8 NOT NULL,
	 can_pull BOOLEAN,
	 can_push BOOLEAN,
	 can_admin BOOLEAN,
	 CONSTRAINT uni_gitee_rcs UNIQUE (user_id,repo_id)
);

-- Create schema and define views for Grafana
CREATE SCHEMA gitee_state;

-- repos list
CREATE OR REPLACE  VIEW gitee_state.repos_list AS
SELECT 
	'Gitee' AS platform,
	r.id AS repo_id,
	r.human_name AS repo_name,
	r.stargazers_count AS star_count,
	r.forks_count AS fork_count,
	r.watchers_count AS watch_count,
	r.created_at AS created_at,
	r.license AS license,
	NOW() AS time    
FROM gitee.repos r;

-- repos stat
CREATE OR REPLACE VIEW gitee_state.repos_stat AS 
SELECT 
	COUNT(r.id) AS total_repos_count,
	SUM(r.stargazers_count) AS total_star_count,
	SUM(r.forks_count) AS total_forks_count,
	SUM(r.watchers_count) AS total_watchers_count,
	NOW() AS time
FROM  gitee.repos r;

-- commits stat
CREATE OR REPLACE VIEW gitee_state.commits_stat AS 
SELECT 
	COUNT(DISTINCT(c.sha)) AS total_commits_count,
	COUNT(DISTINCT(c.author_email)) AS total_author_email,
	COUNT(DISTINCT(c.committer_email)) AS total_committer_email,
	NOW() AS time
FROM  gitee.commits c;

-- issues stat 
CREATE OR REPLACE VIEW gitee_state.issues_stat AS 
SELECT 
	COUNT(i.id) AS total_issue_count,
	COUNT(DISTINCT(i.user_id)) AS total_issue_user,
	(SELECT COUNT(id) FROM gitee.issues WHERE state = 'open') AS open_count,
	(SELECT COUNT(id) FROM gitee.issues WHERE state = 'rejected') AS rejected_count,
	(SELECT COUNT(id) FROM gitee.issues WHERE state = 'closed') AS closed_count,	
	(SELECT COUNT(id) FROM gitee.issues WHERE state = 'progressing') AS progressing_count,
	NOW() AS time
FROM gitee.issues i;

-- pull requests stats
CREATE OR REPLACE VIEW gitee_state.pr_stat AS 
SELECT 
	COUNT(pr.id) AS total_pr_count,
	COUNT(DISTINCT(pr.user_id)) AS total_user_count,
	(SELECT COUNT(id) FROM gitee.pull_requests WHERE state = 'open' ) AS open_count,
	(SELECT COUNT(id) FROM gitee.pull_requests WHERE state = 'merged' ) AS merged_count,
	(SELECT COUNT(id) FROM gitee.pull_requests WHERE state = 'closed' ) AS closed_count,
	NOW() AS time
FROM gitee.pull_requests pr;


-- commits stat by repo
CREATE OR REPLACE VIEW gitee_state.commits_stats_by_repo AS 
SELECT 
	r.id AS repo_id,
	COUNT(c.sha) AS commits_count,
	COUNT(DISTINCT(c.author_email)) AS author_count,
	COUNT(DISTINCT(c.committer_email)) AS committer_count,
	NOW() AS time
FROM  gitee.repos r LEFT JOIN gitee.commits c ON r.id = c.repo_id 
GROUP BY r.id;

-- issues stat by repo 
CREATE OR REPLACE VIEW gitee_state.issues_stats_by_repo AS 
SELECT 
	r.id AS repo_id,
	COUNT(i.id) AS issue_count,
	(SELECT COUNT(DISTINCT(user_id)) FROM gitee.issues WHERE repo_id = r.id) AS issue_user_count,
	(SELECT COUNT(id) FROM gitee.issues WHERE repo_id = r.id AND state = 'open') AS open_cuont,
	(SELECT COUNT(id) FROM gitee.issues WHERE repo_id = r.id AND state = 'rejected') AS rejected_cuont,
	(SELECT COUNT(id) FROM gitee.issues WHERE repo_id = r.id AND state = 'closed') AS closed_cuont,
	(SELECT COUNT(id) FROM gitee.issues WHERE repo_id = r.id AND state = 'progressing') AS progressing_cuont,
	NOW() AS time
FROM gitee.repos r LEFT JOIN gitee.issues i ON r.id = i.repo_id
GROUP BY r.id;

-- pull requests stat by repo 
CREATE OR REPLACE VIEW gitee_state.pr_stats_by_repo AS 
SELECT 
	r.id AS repo_id,
	COUNT(pr.id) AS pr_count,
	(SELECT COUNT(DISTINCT(user_id)) FROM gitee.pull_requests WHERE repo_id = r.id ) AS total_user_count,
	(SELECT COUNT(id) FROM gitee.pull_requests WHERE state = 'open' AND repo_id = r.id ) AS open_count,
	(SELECT COUNT(id) FROM gitee.pull_requests WHERE state = 'merged' AND repo_id = r.id ) AS merged_count,
	(SELECT COUNT(id) FROM gitee.pull_requests WHERE state = 'closed' AND repo_id = r.id ) AS closed_count,
	NOW() AS time
FROM gitee.repos r LEFT JOIN gitee.pull_requests pr ON r.id = pr.repo_id
GROUP BY r.id ;


-- commit author ranking by repo 
CREATE OR REPLACE VIEW gitee_state.commit_author_rank_by_repo AS 
SELECT	
	c.repo_id AS repo_id,
	c.author_email AS author_email,
	COUNT(c.sha) AS commit_count,
	NOW() AS time 
FROM gitee.commits c 
GROUP BY c.repo_id, c.author_email
ORDER BY COUNT(c.sha) DESC;

-- commit committer ranking by repo 
CREATE OR REPLACE VIEW gitee_state.commit_committer_rank_by_repo AS 
SELECT	
	c.repo_id AS repo_id,
	c.committer_email AS committer_email,
	COUNT(c.sha) AS commit_count,
	NOW() AS time 
FROM gitee.commits c 
GROUP BY c.repo_id, c.committer_email
ORDER BY COUNT(c.sha) DESC;


-- pr merge time hours diff
CREATE OR REPLACE VIEW gitee_state.pr_merge_hours_diff AS 
SELECT
	pr.id AS pr_id,
	pr.repo_id AS repo_id,
	pr.created_at AS created_at,
	pr.merged_at AS merged_at,
	EXTRACT(EPOCH FROM (pr.merged_at - pr.created_at))/3600  AS hours_diff
FROM gitee.pull_requests pr WHERE pr.mergeable = TRUE AND pr.state = 'merged';


-- issue close time hours diff
CREATE OR REPLACE VIEW gitee_state.issue_close_hours_diff AS 
SELECT
	iss.repo_id AS repo_id,
	iss.id AS issue_id,
	iss.created_at AS created_at,
	iss.finished_at AS finished_at,
	EXTRACT(EPOCH FROM (iss.finished_at - iss.created_at))/3600  AS hours_diff
FROM gitee.issues iss WHERE iss.state = 'closed' OR iss.state = 'rejected';

-- commit list 
CREATE OR REPLACE VIEW gitee_state.commits_list AS 
SELECT
	'Gitee' AS platform,
	r.id AS repo_id,
	r.full_name AS repo_name,
	c.sha AS sha,	
	c.detail_message AS message, 
	c.author_name AS author_name,
	c.author_email AS author_email,
	c.author_date AS author_date,
	c.committer_name AS committer_name,
	c.committer_email AS committer_email,
	c.committer_date AS committer_date
FROM gitee.commits c , gitee.repos r 
WHERE c.repo_id = r.id;

-- issue list 
CREATE OR REPLACE VIEW gitee_state.issues_list AS 
SELECT 
	'Gitee' AS platform,
	r.id AS repo_id,
	r.full_name AS repo_name,
	iss.state AS issue_state,
	u."name" AS user_name,
	iss.title AS title,
	iss.created_at AS created_at		
FROM gitee.issues iss , gitee.repos r , gitee.users u 
WHERE iss.repo_id = r.id AND iss.user_id = u.id;

-- pull request list 
CREATE OR REPLACE VIEW gitee_state.prs_list AS 
SELECT
	'Gitee' AS platform,
	pr.id AS pr_id,
	pr.repo_id AS repo_id,
	pr."number" AS pr_number,
	pr.created_at AS created_at,
	pr.title AS title,
	pr.mergeable AS mergeable,
	r.full_name AS repo_name,
	u."name" AS user_name
FROM gitee.pull_requests pr , gitee.repos r, gitee.users u 
WHERE pr.repo_id = r.id AND pr.user_id = u.id;