package auth

import (
	"context"
	"log/slog"
	"sso/internal/domain/models"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system and returns
//
//	a JWT token if successful
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appId int,
) (string, error) {
	// Check if user exists
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		return "", err
	}
	panic("not implemented")
	// Check if user has access to the app
}

// RegisterNewUser registers a new user in the system and returns the user's ID
// If the user is already registered, it returns an error
//
//	userId int64, err error
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	panic("not implemented")
}

// IsAdmin checks if the user is an admin
//
//	isAdmin bool, err error
func (a *Auth) IsAdmin(
	ctx context.Context,
	userId int64,
) (bool, error) {
	panic("not implemented")
}