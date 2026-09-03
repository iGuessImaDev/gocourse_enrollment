package enrollment

import (
	"context"
	"log"

	"github.com/iGuessImaDev/gocourse_domain/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, enroll *domain.Enrollment) error
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(db *gorm.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (repo *repo) Create(ctx context.Context, enroll *domain.Enrollment) error {
	if err := repo.db.WithContext(ctx).Create(enroll).Error; err != nil {
		repo.log.Println(err)
		return err
	}
	repo.log.Println("enrollment created with id ", enroll.ID)
	return nil
}
