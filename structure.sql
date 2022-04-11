-- Database Structure For RepoStats
CREATE DATABASE repostats ENCODING 'UTF8';

-- Connect to database repostats
\c repostats


CREATE TABLE public.gitee_repos ( id BIGINT NOT NULL,
	full_name VARCHAR(255) NOT NULL,
	human_name VARCHAR(255) NOT NULL,
	owner_id BIGINT NOT NULL,
	html_url VARCHAR(500) NOT NULL,
	ssh_url VARCHAR(500) NOT NULL,
	recommend BOOLEAN NOT NULL DEFAULT false,
	gvp BOOLEAN NOT NULL DEFAULT false,
	homepage VARCHAR(500),
	language VARCHAR(500),
	forks_count BIGINT NOT NULL DEFAULT 0,
	stargazers_count BIGINT NOT NULL DEFAULT 0,
	watchers_count BIGINT NOT NULL DEFAULT 0,
	open_issues_count BIGINT NOT NULL DEFAULT 0,
	has_issues BOOLEAN NOT NULL DEFAULT false,
	has_wiki BOOLEAN NOT NULL DEFAULT false,
	issue_comment BOOLEAN NOT NULL DEFAULT false,
	can_comment BOOLEAN NOT NULL DEFAULT false,
	pull_requests_enabled BOOLEAN NOT NULL DEFAULT false,
	license VARCHAR(500),
	project_creator VARCHAR(500),
	pushed_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE
);
ALTER TABLE public.gitee_repos ADD CONSTRAINT uni_gitee_repos_id UNIQUE (id);

-- Users
CREATE TABLE public.gitee_users ( id BIGINT NOT NULL,
	login VARCHAR(100) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	html_url VARCHAR(255) NOT NULL,
	"type" VARCHAR(50) NULL
);
ALTER TABLE public.gitee_users ADD CONSTRAINT uni_gitee_users_id UNIQUE (id);

-- Collaborators
CREATE TABLE public.gitee_collaborators ( "user_id" BIGINT NOT NULL, repo_id BIGINT NOT NULL);
ALTER TABLE public.gitee_collaborators ADD CONSTRAINT uni_gitee_rcs UNIQUE (user_id,repo_id);

-- Issues
CREATE TABLE public.gitee_issues ( id BIGINT NOT NULL,
	repo_id BIGINT NOT NULL,
	"user_id" BIGINT NOT NULL,
	html_url VARCHAR(500) NOT NULL,
	"number" VARCHAR(40) NOT NULL,
	"state" VARCHAR(40) NOT NULL,
	scheduled_time INT, comments INT, priority INT, issue_type VARCHAR(40),
	issue_state VARCHAR(40),
	finished_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE
);
ALTER TABLE public.gitee_issues ADD CONSTRAINT uni_gitee_issue UNIQUE (id);

-- Pull requests
CREATE TABLE public.gitee_pullrequests ( id BIGINT NOT NULL,
	repo_id BIGINT NOT NULL,
	"user_id" BIGINT NOT NULL,
	html_url VARCHAR(500) NOT NULL,
	"number" VARCHAR(40) NOT NULL,
	"state" VARCHAR(40) NOT NULL,
	finished_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE,
	closed_at TIMESTAMP WITH TIME ZONE,
	merged_at TIMESTAMP WITH TIME ZONE,
	mergeable BOOLEAN, can_merge_check BOOLEAN
);
ALTER TABLE public.gitee_pullrequests ADD CONSTRAINT uni_gitee_prs UNIQUE (id);

-- Commits
CREATE TABLE public.gitee_commits ( sha VARCHAR(100) NOT NULL,
	repo_id BIGINT NOT NULL,
	author BIGINT NOT NULL,
	html_url VARCHAR(500) NOT NULL,
	commiter BIGINT NOT NULL,
	commit_at TIMESTAMP WITH TIME ZONE
);
ALTER TABLE public.gitee_commits ADD CONSTRAINT uni_gitee_commits UNIQUE(sha);

-- Stargazers
CREATE TABLE public.gitee_stargazers ( user_id BIGINT NOT NULL,
	repo_id BIGINT NOT NULL,
	star_at TIMESTAMP WITH TIME ZONE NOT NULL
);
ALTER TABLE public.gitee_stargazers ADD CONSTRAINT uni_gitee_stargazers UNIQUE(user_id,repo_id);

-- Create schema and define views for Grafana
CREATE SCHEMA gitee_state;

-- Issue Type Detail
CREATE VIEW gitee_state.issue_type_detail AS
SELECT i.ISSUE_TYPE AS issue_type,
	count(i.ID) AS total_count,
	NOW() AS time
FROM public.GITEE_ISSUES i
GROUP BY i.ISSUE_TYPE;

-- Issue State Detail
CREATE VIEW gitee_state.issue_state_detail AS
SELECT i.STATE AS issue_state,
	count(i.ID) AS total_count,
	NOW() AS time
FROM public.GITEE_ISSUES i
GROUP BY i.STATE;

-- Repo Detail
CREATE VIEW gitee_state.repo_detail AS
SELECT 
	r.id AS repo_id,
	r.human_name AS repo_name,
	r.stargazers_count AS star_count,
	r.forks_count AS fork_count,
	r.watchers_count AS watch_count,
	NOW() AS time    
FROM gitee_repos r;

-- Issue by date
CREATE VIEW gitee_state.issue_by_date AS
SELECT date(i.CREATED_AT) AS time,
	count(i.ID) AS issue_count
FROM public.GITEE_ISSUES i
GROUP BY date(i.CREATED_AT)
ORDER BY time DESC;

-- Commit by date
CREATE VIEW gitee_state.commit_by_date AS
SELECT date(c.COMMIT_AT) AS time,
	count(c.SHA) AS commit_count
FROM public.GITEE_COMMITS c
GROUP BY date(c.COMMIT_AT)
ORDER BY time DESC;

-- PR by date
CREATE VIEW gitee_state.pr_by_date AS
SELECT date(pr.CREATED_AT) AS time,
	count(pr.ID) AS pr_count
FROM public.GITEE_PULLREQUESTS pr
GROUP BY date(pr.CREATED_AT)
ORDER BY time DESC;

-- Star by date
CREATE VIEW gitee_state.start_by_date AS
SELECT date(star.STAR_AT) AS time,
	count(star.USER_ID) AS star_count
FROM public.GITEE_STARGAZERS star
GROUP BY date(star.STAR_AT)
ORDER BY time DESC;

-- Star DEtail
CREATE OR REPLACE VIEW gitee_state.star_detail
AS SELECT u.id AS user_id,
    u.name AS user_name,
    r.id AS repo_id,
    r.human_name,
    s.star_at AS "time"
   FROM gitee_repos r,
    gitee_users u,
    gitee_stargazers s
  WHERE s.user_id = u.id AND s.repo_id = r.id;

-- Create Roles & Grant Privileges
CREATE USER repostats_admin WITH PASSWORD '_Rep0stats^123';
CREATE USER repostats_readonly WITH PASSWORD '1R_repoSt(ats987';
ALTER DATABASE repostats OWNER TO repostats_admin;
ALTER SCHEMA public OWNER TO repostats_admin;
ALTER SCHEMA gitee_state OWNER TO repostats_admin;

GRANT ALL PRIVILEGES ON DATABASE repostats TO repostats_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO repostats_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA gitee_state TO repostats_admin;

GRANT CONNECT ON DATABASE repostats TO repostats_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO repostats_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA gitee_state TO repostats_readonly;

