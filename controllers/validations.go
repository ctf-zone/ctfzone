package controllers

import "github.com/ctf-zone/ctfzone/models"

func userValidate(db *models.Repository, u *models.User) (bool, map[string][]string) {

	errs := make(map[string][]string)

	if _, err := db.UsersOneByName(u.Name); err == nil {
		errs["name"] = []string{"User with such name already exists"}
	}

	if _, err := db.UsersOneByEmail(u.Email); err == nil {
		errs["email"] = []string{"User with such email already exists"}
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, nil
}
