-- Database Structure For RepoStats
CREATE DATABASE repostats ENCODING 'UTF8';

CREATE TABLE public.users (
  id serial4 NOT NULL,
	account varchar(200) NOT NULL,
	password text NOT NULL,			
	CONSTRAINT users_pk PRIMARY KEY (id),
	CONSTRAINT users_account_un UNIQUE (account)
);