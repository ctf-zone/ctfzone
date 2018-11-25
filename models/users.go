package models

import (
	"database/sql"
	"fmt"
	"time"

	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"

	"github.com/ctf-zone/ctfzone/internal/crypto"
)

// User represents user profile.
type User struct {
	ID           int64                  `db:"id,omitempty"       json:"id"`
	Name         string                 `db:"name"               json:"name"`
	Email        string                 `db:"email"              json:"email"`
	Password     string                 `db:"-"                  json:"password,omitempty"`
	PasswordHash string                 `db:"password_hash"      json:"-"`
	Extra        map[string]interface{} `db:"extra"              json:"extra"`
	IsActivated  bool                   `db:"is_activated"       json:"isActivated"`
	CreatedAt    time.Time              `db:"created_at"         json:"createdAt"`
	UpdatedAt    time.Time              `db:"updated_at"         json:"updatedAt"`
}

// UserP represents public user profile.
type UserP struct {
	ID    int64                  `db:"id"    json:"id"`
	Name  string                 `db:"name"  json:"name"`
	Extra map[string]interface{} `db:"extra" json:"extra"`
}

type Users []User

func (l Users) First() int64 {
	if len(l) == 0 {
		return 0
	}
	return l[0].ID
}

func (l Users) Last() int64 {
	if len(l) == 0 {
		return 0
	}
	return l[len(l)-1].ID
}

func (l Users) Len() int {
	return len(l)
}

// usersListOptions contains all avaliable options of UsersList method.
type usersListOptions struct {
	Pagination
	UsersFilters
}

func newUsersListOptions(options []usersListOption) *usersListOptions {
	o := &usersListOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

// usersListOption is a functional option type of UsersList method.
type usersListOption func(*usersListOptions)

// UsersFilters contains all avaliable filters of UsersList method.
type UsersFilters struct {
	ID        int64  `schema:"id"`
	Name      string `schema:"name"`
	Email     string `schema:"email"`
	CreatedAt struct {
		From time.Time `schema:"from"`
		To   time.Time `schema:"to"`
	} `schema:"createdAt"`
	IsActivated *bool                  `schema:"isActivated"`
	Extra       map[string]interface{} `schema:"extra"`
}

// UsersPagination option sets pagination params of UsersList method.
func UsersPagination(p Pagination) usersListOption {
	return func(params *usersListOptions) {
		if !p.IsZero() {
			params.Pagination = p
		} else {
			params.Pagination = defaultPagination
		}
	}
}

// UsersFilter option sets filters parameters of UsersList method.
func UsersFilter(f UsersFilters) usersListOption {
	return func(params *usersListOptions) {
		params.UsersFilters = f
	}
}

// UsersInsert saves the user into a database.
func (r *Repository) UsersInsert(o *User) error {

	if hash, err := crypto.HashPassword(o.Password); err != nil {
		return err
	} else {
		o.PasswordHash = hash
		o.Password = ""
	}

	o.CreatedAt = now()
	o.UpdatedAt = o.CreatedAt

	row, err := r.db.
		InsertInto("users").
		Values(o).
		Returning("id").
		QueryRow()

	if err != nil {
		return err
	}

	return row.Scan(&o.ID)
}

// UsersUpdate updates the user in a database.
func (r *Repository) UsersUpdate(o *User) error {

	o.UpdatedAt = now()

	query := r.db.
		Update("users").
		Set(map[string]interface{}{
			"name":         o.Name,
			"email":        o.Email,
			"extra":        o.Extra,
			"is_activated": o.IsActivated,
			"updated_at":   o.UpdatedAt,
		}).
		Where("id", o.ID)

	// Update password if nesessary.
	if o.Password != "" {
		if hash, err := crypto.HashPassword(o.Password); err != nil {
			return err
		} else {
			o.PasswordHash = hash
			query = query.Set("password_hash", hash)
		}
	}

	res, err := query.Exec()

	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n != 1 {
		return sql.ErrNoRows
	}

	return nil
}

// UsersOneByID retrives the user from a database by id.
func (r *Repository) UsersOneByID(id int64) (*User, error) {
	var o User

	err := r.db.
		SelectFrom("users").
		Where("id", id).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}

// UsersOneByEmail retrives the user from a database by email.
func (r *Repository) UsersOneByEmail(email string) (*User, error) {
	var o User

	err := r.db.
		SelectFrom("users").
		Where("email", email).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}

// UsersOneByName retrives the user from a database by id.
func (r *Repository) UsersOneByName(name string) (*User, error) {
	var o User

	err := r.db.
		SelectFrom("users").
		Where("name", name).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}

// UsersLogin returns user if email and login are correct and error otherwise.
func (r *Repository) UsersLogin(email, password string) (*User, error) {

	o, err := r.UsersOneByEmail(email)

	if err != nil {
		// We should pass valid bcrypt hash to bcrypt.CompareHashAndPassword
		// to prevent user enumeration by timing requests.
		_ = crypto.CheckPassword("$2a$10$8OcntURIAZI8nYqDEacSReBce1rqiFPPgEuTBAVu5YHCsd4pwv4E2", password)
		return nil, err
	}

	if err := crypto.CheckPassword(o.PasswordHash, password); err != nil {
		return nil, err
	}

	return o, nil
}

// UsersDelete removes the user from a database by id.
func (r *Repository) UsersDelete(id int64) error {
	res, err := r.db.
		DeleteFrom("users").
		Where("id", id).
		Exec()

	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n != 1 {
		return sql.ErrNoRows
	}

	return nil
}

// usersApplyFilters applies UsersFilters to UsersList query.
func usersApplyFilters(query sqlbuilder.Selector, f UsersFilters) sqlbuilder.Selector {

	// Name.
	if f.Name != "" {
		query = query.Where(
			udb.Cond{
				"name ILIKE": fmt.Sprintf("%%%s%%", f.Name),
			},
		)
	}

	// Email.
	if f.Email != "" {
		query = query.Where(
			udb.Cond{
				"email ILIKE": fmt.Sprintf("%%%s%%", f.Email),
			},
		)
	}

	// Date from and date to.
	if from, to := f.CreatedAt.From, f.CreatedAt.To; !from.IsZero() || !to.IsZero() {
		fromCond := udb.Cond{"created_at >": from}
		toCond := udb.Cond{"created_at <=": to}

		switch {
		case !from.IsZero() && !to.IsZero():
			query = query.Where(udb.And(fromCond, toCond))
		case !from.IsZero():
			query = query.Where(fromCond)
		case !to.IsZero():
			query = query.Where(toCond)
		}
	}

	// IsActivated
	if f.IsActivated != nil {
		query = query.Where("is_activated = ?", *f.IsActivated)
	}

	// Extra fields.
	if f.Extra != nil {
		for k, v := range f.Extra {
			query = query.Where("extra->>? IN ?", k, v)
		}
	}

	return query
}

// UsersList returns list of users.
func (r *Repository) UsersList(opts ...usersListOption) ([]User, *PagesInfo, error) {
	list := make(Users, 0)

	options := newUsersListOptions(opts)

	query := usersApplyFilters(
		r.db.SelectFrom("users"),
		options.UsersFilters,
	)

	if options.Pagination.IsZero() {
		return list, nil, handleErr(query.All(&list))
	}

	pageQuery := query.
		Paginate(options.Pagination.Count).
		Cursor("-id")

	if err := paginate(pageQuery, options.Pagination).All(&list); err != nil {
		return nil, nil, err
	}

	pi, err := r.pagesInfo(pageQuery, options.Pagination, list)

	if err != nil {
		return nil, nil, err
	}

	return list, pi, nil
}
