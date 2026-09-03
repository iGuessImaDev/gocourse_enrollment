package enrollment

import (
	"context"
	"log"

	"github.com/iGuessImaDev/gocourse_domain/domain"
)

type (
	Filters struct {
		UserID   string
		CourseID string
	}

	Service interface {
		Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error) {
	enroll := domain.Enrollment{
		UserId:   userID,
		CourseId: courseID,
		Status:   "P",
	}

	if err := s.repo.Create(ctx, &enroll); err != nil {
		s.log.Println(err)
		return nil, err
	}

	return &enroll, nil
}
