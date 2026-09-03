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
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
		Update(ctx context.Context, id string, status *string) error
		Count(ctx context.Context, filters Filters) (int, error)
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

func (repo *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	var e []domain.Enrollment

	tx := repo.db.WithContext(ctx).Model(&e)
	tx = applyFilters(tx, filters)
	tx = tx.Offset(offset).Limit(limit)
	result := tx.Order("created_at desc").Find(&e)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}

	return e, nil
}

func (repo *repo) Update(ctx context.Context, id string, status *string) error {
	values := make(map[string]interface{})

	if status != nil {
		values["status"] = *status
	}
	result := repo.db.WithContext(ctx).Model(&domain.Enrollment{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		repo.log.Printf("enrollment %s doesn't exist", id)
		return ErrNotFound{id}
	}
	return nil
}

func (repo *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.Enrollment{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, err
	}
	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.UserID != "" {
		tx = tx.Where("user_id = ?", filters.UserID)
	}

	if filters.CourseID != "" {
		tx = tx.Where("course_id = ?", filters.CourseID)
	}

	return tx
}
