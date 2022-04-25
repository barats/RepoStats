package storage

import (
	"repostats/model"
	"repostats/utils"
	"strings"
)

func NewAdmin(account string, password string) error {
	query := `INSERT INTO public.users (account, "password") VALUES(:account,:password)`
	data, err := PasswordBase58Hash(password)
	if err != nil {
		return err
	}
	return DbNamedExec(query, model.Admin{Account: account, Password: data})
}

func UpdateAdmin(user model.Admin) error {
	query := `UPDATE public.users SET account = :account , "password" = :password WHERE id = :id`
	return DbNamedExec(query, user)
}

func FindAdminByAccount(account string) (model.Admin, error) {
	var user model.Admin
	query := `SELECT * FROM public.users u WHERE lower(u.account) = $1`
	return user, DbGet(query, &user, strings.ToLower(account))
}

func PasswordBase58Hash(password string) (string, error) {
	data, err := utils.Sha256Of(password)
	if err != nil {
		return "", err
	}
	return utils.Base58Encode(data), nil
}
