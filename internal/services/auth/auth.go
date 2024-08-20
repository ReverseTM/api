package auth

import (
	"api/internal/domain/models"
	myerrors "api/internal/errors"
	"api/internal/lib/jwt"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type AuthService struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	secretKey    string
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(username string, password []byte) (int64, error)
}

type UserProvider interface {
	User(username string) (models.User, error)
}

func NewAuthService(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	secretKey string,
	tokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		secretKey:    secretKey,
		tokenTTL:     tokenTTL,
	}
}

func (a *AuthService) Login(username, password string) (string, error) {
	const op = "services.auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(username)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotFound) {
			log.Error("user not found")
			return "", fmt.Errorf("%s: %w", op, myerrors.ErrInvalidCredentials)
		}

		log.Error("failed to get user")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Error("invalid password")
		return "", fmt.Errorf("%s: %w", op, myerrors.ErrInvalidCredentials)
	}

	token, err := jwt.GenerateToken(user, a.secretKey, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AuthService) Register(username, password string) (int64, error) {
	const op = "services.auth.Register"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to register user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(username, passwordHash)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserExists) {
			log.Error("user already exists")
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save user")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
