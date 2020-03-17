package repo

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pytorchtw/go-jwt-server/repo/models"
	"github.com/volatiletech/sqlboiler/boil"
)

type Repo interface {
	/*
		List    func(dest interface{}, query string) error
	*/
	GetUser(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) (int64, error)
	DeleteUser(email string) (int64, error)
	Close() error
}

type DBRepo struct {
	db  *sql.DB
	ctx context.Context
}

func NewDBRepo(db *sql.DB) *DBRepo {
	repo := DBRepo{}
	boil.SetDB(db)
	repo.db = db
	repo.ctx = context.Background()
	return &repo
}

func (repo *DBRepo) GetUser(email string) (*models.User, error) {
	user, err := models.Users(models.UserWhere.Email.EQ(email)).One(repo.ctx, repo.db)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *DBRepo) CreateUser(user *models.User) error {
	err := user.Insert(repo.ctx, repo.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBRepo) UpdateUser(user *models.User) (int64, error) {
	rowsAff, err := user.Update(repo.ctx, repo.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return rowsAff, nil
}

func (repo *DBRepo) DeleteUser(email string) (int64, error) {
	user, err := repo.GetUser(email)
	rowsAff, err := user.Delete(repo.ctx, repo.db)
	if err != nil {
		return 0, err
	}
	return rowsAff, nil
}

func (repo *DBRepo) Close() error {
	err := repo.db.Close()
	return err
}
