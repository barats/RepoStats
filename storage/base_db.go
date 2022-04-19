// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package storage

import (
	"database/sql"
	"fmt"
	"repostats/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	DEFAULT_PAGE_NUMBER   = 1
	DEFAULT_PAGE_SIZE     = 15
	DEFAULT_MAX_PAGE_SIZE = 50
)

var dbService = &DatabaseService{}

type DatabaseService struct {
	Connection *sqlx.DB
}

func InitDatabaseService() (*DatabaseService, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		utils.DatabaseConifg.Host, utils.DatabaseConifg.Port, utils.DatabaseConifg.User,
		utils.DatabaseConifg.Password, utils.DatabaseConifg.DbName)
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return dbService, err
	}
	conn.SetMaxOpenConns(utils.DatabaseConifg.MaxOpenConns)
	conn.SetMaxIdleConns(utils.DatabaseConifg.MaxIdleConn)
	conn.SetConnMaxLifetime(0) //always REUSE

	// Unsafe returns a version of DB which will silently succeed to scan when
	// columns in the SQL result have no fields in the destination struct.
	dbService.Connection = conn.Unsafe()
	return dbService, nil
}

func DbNamedExec(query string, args interface{}) error {
	_, err := dbService.Connection.NamedExec(query, args)
	return err
}

func DbExec(query string, args ...interface{}) error {
	_, err := dbService.Connection.Exec(query, args...)
	return err
}

func DbExecTx(query string, args ...interface{}) error {
	tx, err := dbService.Connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	stmt, err := tx.Prepare(dbService.Connection.Rebind(query))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

func DbGet(query string, dest interface{}, args ...interface{}) error {
	err := dbService.Connection.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func DbSelect(query string, dest interface{}, args ...interface{}) error {
	return dbService.Connection.Select(dest, query, args...)
}

func DbClose() {
	dbService.Connection.Close()
}
