package services

import (
	"deni-be-crm/config"
	"deni-be-crm/internal/common"
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/repositories"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(email, password string) (string, error)
	Register(user *models.RegisterRequest) error
}

type AuthService struct {
	userRepo repositories.IUserRepository
}

func NewAuthService(repo repositories.IUserRepository) IAuthService {
	return &AuthService{userRepo: repo}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", common.NewHTTPError(http.StatusBadRequest, "invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", common.NewHTTPError(http.StatusBadRequest, "invalid credentials")
	}

	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(req *models.RegisterRequest) error {

	// check email exist
	emailExist, _ := s.userRepo.IsEmailExist(req.Email)
	fmt.Println(emailExist)
	if emailExist {
		return common.NewHTTPError(http.StatusBadRequest, "Email already registed")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if !models.IsValidRole(req.Role) {
		return common.NewHTTPError(http.StatusBadRequest, "Invalid role")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	return s.userRepo.Create(user)
}

func generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"role":     user.Role,
		"isLeader": user.IsLeader,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.GetEnv("JWT_SECRET", "default_secret")
	return token.SignedString([]byte(secret))
}
