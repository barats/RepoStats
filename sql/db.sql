-- Database Structure For RepoStats
CREATE DATABASE repostats ENCODING 'UTF8';

\c repostats

CREATE TABLE public.users (
  id serial4 NOT NULL,
	account varchar(200) NOT NULL,
	password text NOT NULL,			
	CONSTRAINT users_pk PRIMARY KEY (id),
	CONSTRAINT users_account_un UNIQUE (account)
);

INSERT INTO public.users(account, "password") VALUES('repostats', 'EZ2zQjC3fqbkvtggy9p2YaJiLwx1kKPTJxvqVzowtx6t');
