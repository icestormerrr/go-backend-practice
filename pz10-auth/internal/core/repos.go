package core

type UserRepo interface {
	GetById(id int64) (User, error)
	CheckPassword(email, pass string) (User, error)
}

type SessionRepo interface {
	SetRefreshToken(userID int64, refreshToken string) error
	GetRefreshToken(userID int64) (string, error)
	IncLoginAttempts(email string) (int64, error)
	ResetLoginAttempts(email string) error
}
