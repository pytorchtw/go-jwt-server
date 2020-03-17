package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pytorchtw/go-jwt-server/repo"
	"github.com/pytorchtw/go-jwt-server/repo/models"
	"log"
	"time"
)

type Token struct {
	Token string `json:"token"`
}

const (
	SecretKey = "this is a secretkey"
)

type Service interface {
	CreateUser(user *models.User) (int, error)
	VerifyUser(email string, pass string) (bool, error)
	/*
		Register(user *User) (entity.ID, error)
		ForgotPassword(user *User) error
		ChangePassword(user *User, password string) error
		Validate(user *User) error
		Auth(user *User, password string) error
		IsValid(user *User) bool
		GetRepo() Repository
	*/
	CreateToken() (*Token, error)
}

type WebService struct {
	Repo repo.Repo
}

func NewWebService(db *sql.DB) *WebService {
	ws := WebService{}
	ws.Repo = repo.NewDBRepo(db)
	return &ws
}

func (ws *WebService) CreateUser(user *models.User) error {
	tmp, err := ws.Repo.GetUser(user.Email)
	if err != nil {
		log.Println(err.Error())
	}
	if tmp != nil {
		// if user is non-existing then should have err
		return errors.New("user already exists")
	}

	err = ws.Repo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (ws *WebService) CreateToken() (*Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return nil, err
	}
	val := &Token{tokenString}
	return val, nil
}

func (ws *WebService) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("error parsing claim")
	}

	fmt.Println(claims["foo"], claims["nbf"])
	return nil
}

func (ws *WebService) VerifyUser(email string, pass string) (bool, error) {
	user, err := ws.Repo.GetUser(email)
	if err != nil {
		return false, err
	}
	if user.Password != pass {
		return false, errors.New("wrong password")
	}
	return true, nil
}
