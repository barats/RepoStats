-- Connect to database repostats
\c repostats

-- Create schema and define views
CREATE SCHEMA gitee;

-- Repos
CREATE TABLE gitee.repos ( 
	id BIGINT NOT NULL,
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
	license VARCHAR(500),
	project_creator VARCHAR(500),
	pushed_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE
);
ALTER TABLE gitee.repos ADD CONSTRAINT uni_gitee_repos_id UNIQUE (id);

-- Users
CREATE TABLE gitee.users ( 
	id BIGINT NOT NULL,
	login VARCHAR(100) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	html_url VARCHAR(255) NOT NULL,
	"type" VARCHAR(50) NULL
);
ALTER TABLE gitee.users ADD CONSTRAINT uni_gitee_users_id UNIQUE (id);

-- Collaborators
CREATE TABLE gitee.collaborators (
	 "user_id" BIGINT NOT NULL, 
	 repo_id BIGINT NOT NULL
);
ALTER TABLE gitee.collaborators ADD CONSTRAINT uni_gitee_rcs UNIQUE (user_id,repo_id);

-- Issues
CREATE TABLE gitee.issues ( 
	id BIGINT NOT NULL,
	repo_id BIGINT NOT NULL,
	"user_id" BIGINT NOT NULL,
	html_url VARCHAR(500) NOT NULL,
	"number" VARCHAR(40) NOT NULL,
	"state" VARCHAR(40) NOT NULL,
	scheduled_time INT, 
	comments INT, 
	priority INT, 
	issue_type VARCHAR(40),
	issue_state VARCHAR(40),
	finished_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE
);
ALTER TABLE gitee.issues ADD CONSTRAINT uni_gitee_issue UNIQUE (id);

-- Pull requests
CREATE TABLE gitee.pullrequests ( 
	id BIGINT NOT NULL,
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
	mergeable BOOLEAN, 
	can_merge_check BOOLEAN
);
ALTER TABLE gitee.pullrequests ADD CONSTRAINT uni_gitee_prs UNIQUE (id);

-- Commits
CREATE TABLE gitee.commits ( 
	sha VARCHAR(100) NOT NULL,
	repo_id BIGINT NOT NULL,
	author BIGINT NOT NULL,
	html_url VARCHAR(500) NOT NULL,
	commiter BIGINT NOT NULL,
	commit_at TIMESTAMP WITH TIME ZONE
);
ALTER TABLE gitee.commits ADD CONSTRAINT uni_gitee_commits UNIQUE(sha,repo_id);

-- Stargazers
CREATE TABLE gitee.stargazers ( 
	user_id BIGINT NOT NULL,
	repo_id BIGINT NOT NULL,
	star_at TIMESTAMP WITH TIME ZONE NOT NULL
);
ALTER TABLE gitee.stargazers ADD CONSTRAINT uni_gitee_stargazers UNIQUE(user_id,repo_id);

