-- Connect to database repostats
\c repostats

-- Create schema and define views
CREATE SCHEMA gitee;

-- Repos
CREATE TABLE gitee.repos ( 
	id int8 NOT NULL,
	full_name VARCHAR(1000),
	human_name VARCHAR(1000),
	path VARCHAR(1000),
	name VARCHAR(1000),
	url VARCHAR(1000),
	owner_id int8,
	assigner_id int8,
	description VARCHAR(1000),
	html_url VARCHAR(2000),
	ssh_url VARCHAR(2000),
	forked_repo BOOLEAN,
	default_branch VARCHAR(1000),		
	forks_count INT ,
	stargazers_count INT,
	watchers_count INT,
	license VARCHAR(1000),
	pushed_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE,
	updated_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT uni_gitee_repo_id UNIQUE (id)
);

-- Commits
CREATE TABLE gitee.commits (
	sha VARCHAR(80) NOT NULL,
	repo_id int8 NOT NULL,	
	html_url VARCHAR(2000),	
	author_name VARCHAR(500),
	author_email VARCHAR(500),
	author_date TIMESTAMP WITH TIME ZONE,
	committer_name VARCHAR(200) ,
	committer_email VARCHAR(200) ,
	committer_date TIMESTAMP WITH TIME ZONE ,
	detail_message TEXT,
	tree VARCHAR(80),	
	CONSTRAINT uni_gitee_commits UNIQUE (sha,repo_id)
);

-- Issues
CREATE TABLE gitee.issues ( 
	id int8 NOT NULL,		
	"user_id" int8 NOT NULL,
	repo_id int8 NOT NULL,
	html_url VARCHAR(2000),
	"number" VARCHAR(100),
	"state" VARCHAR(100),
	title VARCHAR(1000),	
	finished_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE,
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
	id int8 NOT NULL,
	repo_id int8 NOT NULL,
	"user_id" int8 NOT NULL,		
	html_url VARCHAR(2000),
	diff_url VARCHAR(2000),
	patch_url VARCHAR(2000),
	"number" int8,
	"state" VARCHAR(100),
	created_at TIMESTAMP WITH TIME ZONE,
	updated_at TIMESTAMP WITH TIME ZONE,
	closed_at TIMESTAMP WITH TIME ZONE,
	merged_at TIMESTAMP WITH TIME ZONE,
	mergeable BOOLEAN, 
	can_merge_check BOOLEAN,
	title varchar(1000),
	head_label VARCHAR(100),
	head_ref VARCHAR(100),
	head_sha VARCHAR(100),
	head_user_id int8,	
	head_repo_id int8,		
	CONSTRAINT uni_gitee_prs UNIQUE (id,repo_id)
);

-- Users
CREATE TABLE gitee.users ( 
	id int8 NOT NULL,
	login VARCHAR(500),
	"name" VARCHAR(1000),
	avatar_url VARCHAR(1000),
	html_url VARCHAR(1000),
	remark VARCHAR(1000), 
	"type" VARCHAR(500),
	email VARCHAR(1000),
	created_at TIMESTAMP WITH TIME ZONE,
	CONSTRAINT uni_gitee_users_id UNIQUE (id)
);

-- Stargazers
CREATE TABLE gitee.stargazers ( 
	user_id int8 NOT NULL,
	repo_id int8 NOT NULL,
	star_at TIMESTAMP WITH TIME ZONE NOT NULL,
	CONSTRAINT uni_gitee_stargazers UNIQUE(user_id,repo_id)
);

-- Collaborators
CREATE TABLE gitee.collaborators (
	 "user_id" int8 NOT NULL, 
	 repo_id int8 NOT NULL,
	 CONSTRAINT uni_gitee_rcs UNIQUE (user_id,repo_id)
);


