package enrollment_test

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	courseSdk "github.com/iGuessImaDev/go_course_sdk/course/mock"
	userSdk "github.com/iGuessImaDev/go_course_sdk/user/mock"

	"github.com/iGuessImaDev/gocourse_domain/domain"
	"github.com/iGuessImaDev/gocourse_enrollment/internal/enrollment"
	"github.com/stretchr/testify/assert"
)

func TestService_GetAll(t *testing.T) {
	l := log.New(io.Discard, "", 0)
	t.Run("should return an error", func(t *testing.T) {
		var want error = errors.New("my error")
		var wantCounter int = 1
		var counter int = 0
		repo := &mockRepository{
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				counter++
				return nil, want
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)
		enrrollments, err := service.GetAll(context.Background(), enrollment.Filters{}, 0, 10)
		assert.Error(t, err)
		assert.Nil(t, enrrollments)
		assert.Equal(t, wantCounter, counter)
		assert.EqualError(t, want, err.Error())
	})

	t.Run("should return all enrollments", func(t *testing.T) {
		want := []domain.Enrollment{
			{
				ID:       "1",
				UserID:   "11",
				CourseID: "22",
				Status:   "P",
			},
		}
		var wantCounter int = 1
		var counter int = 0
		repo := &mockRepository{
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				counter++
				return []domain.Enrollment{
					{
						ID:       "1",
						UserID:   "11",
						CourseID: "22",
						Status:   "P",
					},
				}, nil
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)
		enrrollments, err := service.GetAll(context.Background(), enrollment.Filters{}, 0, 10)
		assert.Nil(t, err)
		assert.NotNil(t, enrrollments)
		assert.Equal(t, wantCounter, counter)
		assert.Equal(t, want, enrrollments)
	})
}

func TestService_Update(t *testing.T) {
	l := log.New(io.Discard, "", 0)
	t.Run("should return an error", func(t *testing.T) {
		var want error = errors.New("my error")
		var wantCounter int = 1
		var counter int = 0
		repo := &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				counter++
				return errors.New("my error")
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)
		status := "A"
		err := service.Update(context.Background(), "11", &status)
		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, counter)
		assert.EqualError(t, want, err.Error())
	})

	t.Run("should update an enrollment", func(t *testing.T) {
		var wantCounter int = 1
		var counter int = 0
		var wantStatus string = "A"
		var wantID string = "11"
		repo := &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				counter++
				assert.Equal(t, wantID, id)
				assert.NotNil(t, status)
				assert.Equal(t, wantStatus, *status)
				return nil
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)
		status := "A"
		err := service.Update(context.Background(), "11", &status)
		assert.Nil(t, err)
		assert.Equal(t, wantCounter, counter)
	})
}

func TestService_Count(t *testing.T) {
	l := log.New(io.Discard, "", 0)
	t.Run("should return an error", func(t *testing.T) {
		var want error = errors.New("my error")
		var wantCounter int = 1
		var counter int = 0
		repo := &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				counter++
				return 0, errors.New("my error")
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)
		count, err := service.Count(context.Background(), enrollment.Filters{})
		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, counter)
		assert.EqualError(t, want, err.Error())
		assert.Zero(t, count)
	})

	t.Run("should return the count of enrollments", func(t *testing.T) {
		var wantCounter int = 1
		var counter int = 0
		var wantCount int = 10
		repo := &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				counter++
				return 10, nil
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)
		count, err := service.Count(context.Background(), enrollment.Filters{})
		assert.Nil(t, err)
		assert.Equal(t, wantCounter, counter)
		assert.Equal(t, wantCount, count)
	})
}

func TestService_Create(t *testing.T) {
	l := log.New(io.Discard, "", 0)

	t.Run("should return an error in user sdk", func(t *testing.T) {
		var want error = errors.New("my error")
		var wantCounter int = 1
		var counter int = 0
		userSdk := &userSdk.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				counter++
				return nil, errors.New("my error")
			},
		}

		service := enrollment.NewService(l, userSdk, nil, nil)
		enrollment, err := service.Create(context.Background(), "11", "22")

		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, counter)
		assert.EqualError(t, want, err.Error())
		assert.Nil(t, enrollment)
	})

	t.Run("should return an error in course sdk", func(t *testing.T) {
		var want error = errors.New("my error")
		var wantCounter int = 2
		var counter int = 0
		userSdk := &userSdk.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				counter++
				return nil, nil
			},
		}

		courseSdk := &courseSdk.CourseSdkMock{
			GetMock: func(id string) (*domain.Course, error) {
				counter++
				return nil, errors.New("my error")
			},
		}

		service := enrollment.NewService(l, userSdk, courseSdk, nil)
		enrollment, err := service.Create(context.Background(), "11", "22")

		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, counter)
		assert.EqualError(t, want, err.Error())
		assert.Nil(t, enrollment)
	})

	t.Run("should return an error in repository", func(t *testing.T) {
		var want error = errors.New("my error")
		var wantCounter int = 3
		var counter int = 0
		userSdk := &userSdk.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				counter++
				return nil, nil
			},
		}

		courseSdk := &courseSdk.CourseSdkMock{
			GetMock: func(id string) (*domain.Course, error) {
				counter++
				return nil, nil
			},
		}

		repo := &mockRepository{
			CreateMock: func(ctx context.Context, enroll *domain.Enrollment) error {
				counter++
				return errors.New("my error")
			},
		}

		service := enrollment.NewService(l, userSdk, courseSdk, repo)
		enrollment, err := service.Create(context.Background(), "11", "22")

		assert.NotNil(t, err)
		assert.Equal(t, wantCounter, counter)
		assert.EqualError(t, want, err.Error())
		assert.Nil(t, enrollment)
	})

	t.Run("should create an enrollment", func(t *testing.T) {
		var wantCounter int = 3
		var counter int = 0
		var wantUserID string = "11"
		var wantCourseID string = "22"
		wantStatus := domain.Pending
		var wantID string = "123"
		userSdk := &userSdk.UserSdkMock{
			GetMock: func(id string) (*domain.User, error) {
				counter++
				assert.Equal(t, wantUserID, id)
				return nil, nil
			},
		}
		courseSdk := &courseSdk.CourseSdkMock{
			GetMock: func(id string) (*domain.Course, error) {
				counter++
				assert.Equal(t, wantCourseID, id)
				return nil, nil
			},
		}
		repo := &mockRepository{
			CreateMock: func(ctx context.Context, enroll *domain.Enrollment) error {
				counter++
				enroll.ID = "123"
				return nil
			},
		}

		service := enrollment.NewService(l, userSdk, courseSdk, repo)
		enrollment, err := service.Create(context.Background(), "11", "22")

		assert.Nil(t, err)
		assert.Equal(t, wantCounter, counter)
		assert.NotNil(t, enrollment)
		assert.Equal(t, wantID, enrollment.ID)
		assert.Equal(t, wantUserID, enrollment.UserID)
		assert.Equal(t, wantCourseID, enrollment.CourseID)
		assert.Equal(t, wantStatus, enrollment.Status)
	})
}
