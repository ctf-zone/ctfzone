package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ctf-zone/ctfzone/internal/crypto"
)

// User represents user profile.
type User struct {
	ID           int64     `db:"id,omitempty"       json:"id"`
	Name         string    `db:"name"               json:"name"`
	Email        string    `db:"email"              json:"email"`
	Password     string    `db:"-"                  json:"password,omitempty"`
	PasswordHash string    `db:"password_hash"      json:"-"`
	Extra        Extra     `db:"extra"              json:"extra"`
	IsActivated  bool      `db:"is_activated"       json:"isActivated"`
	CreatedAt    time.Time `db:"created_at"         json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at"         json:"updatedAt"`
}

// TODO: change to struct when done
type Extra map[string]interface{}

func (a Extra) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Extra) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

// UserP represents public user profile.
type UserP struct {
	ID    int64  `db:"id"    json:"id"`
	Name  string `db:"name"  json:"name"`
	Extra Extra  `db:"extra" json:"extra"`
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

	stmt, err := r.db.PrepareNamed(
		"INSERT INTO users (name, email, password_hash, extra, is_activated, created_at, updated_at) " +
			"VALUES(:name, :email, :password_hash, :extra, :is_activated, :created_at, :updated_at) " +
			"RETURNING id")

	if err != nil {
		return err
	}

	return stmt.QueryRowx(o).Scan(&o.ID)
}

// UsersUpdate updates the user in a database.
func (r *Repository) UsersUpdate(o *User) error {

	o.UpdatedAt = now()

	query := "UPDATE users SET " +
		"name = :name, " +
		"email = :email, " +
		"extra = :extra, " +
		"is_activated = :is_activated, " +
		"updated_at = :updated_at"

	// Update password if nesessary.
	if o.Password != "" {
		if hash, err := crypto.HashPassword(o.Password); err != nil {
			return err
		} else {
			o.PasswordHash = hash
			query += ", password_hash = :password_hash"
		}
	}

	query += " WHERE id = :id RETURNING updated_at"

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}

	return stmt.QueryRowx(o).Scan(&o.UpdatedAt)
}

// UsersOneByID retrives the user from a database by id.
func (r *Repository) UsersOneByID(id int64) (*User, error) {
	var o User

	err := r.db.Get(&o, "SELECT * FROM users WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

// UsersOneByEmail retrives the user from a database by email.
func (r *Repository) UsersOneByEmail(email string) (*User, error) {
	var o User

	err := r.db.Get(&o, "SELECT * FROM users WHERE email = $1", email)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

// UsersOneByName retrives the user from a database by id.
func (r *Repository) UsersOneByName(name string) (*User, error) {
	var o User

	err := r.db.Get(&o, "SELECT * FROM users WHERE name = $1", name)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

// UsersLogin returns user if email and login are correct and error otherwise.
func (r *Repository) UsersLogin(email, password string) (*User, error) {

	o, err := r.UsersOneByName(email)

	if err != nil {
		// TODO: check
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
	return r.db.QueryRow("DELETE FROM users WHERE id = $1 RETURNING id", id).Scan(&id)
}

// TODO: rename
// usersApplyFilters applies UsersFilters to UsersList query.
func usersApplyFilters(query string, f UsersFilters) (string, map[string]interface{}) {
	cond := make([]string, 0)
	params := make(map[string]interface{})

	// Name.
	if f.Name != "" {
		cond = append(cond, "name ILIKE :name")
		params["name"] = fmt.Sprintf("%%%s%%", f.Name)
	}

	// Email.
	if f.Email != "" {
		cond = append(cond, "name ILIKE :email")
		params["email"] = fmt.Sprintf("%%%s%%", f.Email)
	}

	// Date from.
	if !f.CreatedAt.From.IsZero() {
		cond = append(cond, ":created_at >= :from")
		params["from"] = f.CreatedAt.From
	}

	// Date to.
	if !f.CreatedAt.To.IsZero() {
		cond = append(cond, ":created_at <= :to")
		params["from"] = f.CreatedAt.To
	}

	// IsActivated
	if f.IsActivated != nil {
		cond = append(cond, "name ILIKE :is_activated")
		params["is_activated"] = *f.IsActivated
	}

	// Extra fields.
	// if f.Extra != nil {
	// TODO: whitelist
	// for k, v := range f.Extra {
	// 	key := fmt.Sprintf("extra_%s_key", k)
	// 	value := fmt.Sprintf("extra_%s_value", k)
	// 	cond = append(cond, fmt.Sprintf("extra ->> :%s = :%s", key, value))
	// 	params[key] = k
	// 	params[value] = v
	// }
	// extra, _ := json.Marshal(f.Extra)
	// params["extra"] = extra
	// }

	if len(cond) > 0 {
		query += " WHERE " + strings.Join(cond, " AND ")
	}

	return query, params
}

// UsersList returns list of users.
func (r *Repository) UsersList(opts ...usersListOption) ([]User, *PagesInfo, error) {
	list := make(Users, 0)

	options := newUsersListOptions(opts)

	query, params := usersApplyFilters("SELECT * FROM users", options.UsersFilters)
	pageQuery, params := paginate(query, params, options.Pagination)

	stmt, err := r.db.PrepareNamed(pageQuery)
	if err != nil {
		return nil, nil, err
	}

	if err := stmt.Select(&list, params); err != nil {
		return nil, nil, err
	}

	pi, err := r.pagesInfo(query, params, options.Pagination, list)

	if err != nil {
		return nil, nil, err
	}

	return list, pi, nil
}
