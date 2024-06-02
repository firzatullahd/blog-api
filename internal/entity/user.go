package entity

import "time"

type User struct {
	ID       uint64 `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Name     string `db:"name"`
	Role     string `db:"role"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type UserRole int

const (
	RoleUser UserRole = iota + 1
	RoleAdmin
)

func (e UserRole) String() string {
	switch e {
	case RoleUser:
		return "user"
	case RoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}
