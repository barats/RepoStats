package model

import "reflect"

type Admin struct {
	ID       int    `db:"id"`
	Account  string `db:"account"`
	Password string `db:"password"`
}

func (user Admin) IsEmpty() bool {
	return reflect.DeepEqual(user, Admin{})
}
