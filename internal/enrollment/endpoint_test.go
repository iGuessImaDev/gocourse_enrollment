package enrollment_test

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"testing"

	courseSdk "github.com/iGuessImaDev/go_course_sdk/course"
	userSdk "github.com/iGuessImaDev/go_course_sdk/user"

	mockCourseSdk "github.com/iGuessImaDev/go_course_sdk/course/mock"
	mockUserSdk "github.com/iGuessImaDev/go_course_sdk/user/mock"

	"github.com/iGuessImaDev/go_lib_response/response"
	"github.com/iGuessImaDev/gocourse_domain/domain"
	"github.com/iGuessImaDev/gocourse_enrollment/internal/enrollment"
	"github.com/stretchr/testify/assert"
)

func TestCreateEndpoint(t *testing.T) {
	l := log.New(io.Discard, "", 0)
	t.Run("should return bad request when user id is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrUserIDRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})
	t.Run("should return bad request when course id is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{UserID: "123"})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrCourseIDRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	obj := []struct {
		tag            string
		repositoryMock enrollment.Repository
		userSdkMock    userSdk.Transport
		courseSdkMock  courseSdk.Transport
		wantErr        error
		wantCode       int
		wantResponse   *domain.Enrollment
	}{
		{
			tag: "should return an error if user sdk returns an unexpected error",
			userSdkMock: &mockUserSdk.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if user does not exist",
			userSdkMock: &mockUserSdk.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, userSdk.ErrNotFound{Message: "user not found"}
				},
			},
			wantErr:  userSdk.ErrNotFound{Message: "user not found"},
			wantCode: http.StatusNotFound,
		},
		{
			tag: "should return an error if course sdk returns an unexpected error",
			userSdkMock: &mockUserSdk.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &mockCourseSdk.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if course does not exist",
			userSdkMock: &mockUserSdk.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &mockCourseSdk.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, courseSdk.ErrNotFound{Message: "course not found"}
				},
			},
			wantErr:  courseSdk.ErrNotFound{Message: "course not found"},
			wantCode: http.StatusNotFound,
		},
		{
			tag: "should return an error if repository returns an unexpected error",
			userSdkMock: &mockUserSdk.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &mockCourseSdk.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, nil
				},
			},
			repositoryMock: &mockRepository{
				CreateMock: func(ctx context.Context, enroll *domain.Enrollment) error {
					return errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return an enrollment",
			userSdkMock: &mockUserSdk.UserSdkMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &mockCourseSdk.CourseSdkMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, nil
				},
			},
			repositoryMock: &mockRepository{
				CreateMock: func(ctx context.Context, enroll *domain.Enrollment) error {
					enroll.ID = "10010"
					return nil
				},
			},
			wantCode: http.StatusCreated,
			wantResponse: &domain.Enrollment{
				ID:       "10010",
				UserID:   "1",
				CourseID: "4",
				Status:   "P",
			},
		},
	}

	for _, obj := range obj {
		t.Run(obj.tag, func(t *testing.T) {
			service := enrollment.NewService(l, obj.userSdkMock, obj.courseSdkMock, obj.repositoryMock)
			endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
			resp, err := endpoint.Create(context.Background(), enrollment.CreateReq{UserID: "1", CourseID: "4"})

			if obj.wantErr != nil {
				assert.NotNil(t, err)
				assert.Nil(t, resp)

				r := err.(response.Response)
				assert.EqualError(t, obj.wantErr, r.Error())
				assert.Equal(t, obj.wantCode, r.StatusCode())
			} else {
				assert.NotNil(t, resp)
				assert.Nil(t, err)

				r := resp.(response.Response)
				assert.Equal(t, obj.wantCode, r.StatusCode())
				assert.Empty(t, r.Error())

				enrollment := r.GetData().(*domain.Enrollment)
				assert.Equal(t, obj.wantResponse.ID, enrollment.ID)
				assert.Equal(t, obj.wantResponse.UserID, enrollment.UserID)
				assert.Equal(t, obj.wantResponse.CourseID, enrollment.CourseID)
				assert.Equal(t, obj.wantResponse.Status, enrollment.Status)
			}
		})
	}
}

func TestGetAllEndpoint(t *testing.T) {
	l := log.New(io.Discard, "", 0)
	t.Run("should return an error if Count returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 0, errors.New("unexpected error")
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{Page: 0})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return an error if meta returns a parsing error", func(t *testing.T) {
		wantErr := errors.New("strconv.Atoi: parsing \"invalid number\": invalid syntax")
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 3, nil
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimPageDef: "invalid number"})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{Page: 0})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return an error if GetAll repository returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 3, nil
			},
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				return nil, errors.New("unexpected error")
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimPageDef: "10"})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{Page: 0})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return all enrollments", func(t *testing.T) {
		wantEnrollments := []domain.Enrollment{
			{ID: "1", UserID: "11", CourseID: "111", Status: "P"},
			{ID: "2", UserID: "22", CourseID: "222", Status: "P"},
			{ID: "3", UserID: "33", CourseID: "333", Status: "P"},
		}
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			CountMock: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 3, nil
			},
			GetAllMock: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				return []domain.Enrollment{
					{ID: "1", UserID: "11", CourseID: "111", Status: "P"},
					{ID: "2", UserID: "22", CourseID: "222", Status: "P"},
					{ID: "3", UserID: "33", CourseID: "333", Status: "P"},
				}, nil
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimPageDef: "10"})
		resp, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{Page: 0})
		assert.Nil(t, err)

		r := resp.(response.Response)
		assert.Equal(t, http.StatusOK, r.StatusCode())
		assert.Empty(t, r.Error())

		enrollments := r.GetData().([]domain.Enrollment)
		assert.Equal(t, wantEnrollments, enrollments)
	})
}

func TestUpdateEndpoint(t *testing.T) {
	l := log.New(io.Discard, "", 0)
	t.Run("should return bad request when status is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		status := ""
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrStatusRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("should return an error if repository returns a not found error", func(t *testing.T) {
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				return enrollment.ErrNotFound{EnrollmentID: id}
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "A"
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrNotFound{EnrollmentID: "20"}, resp.Error())
		assert.Equal(t, http.StatusNotFound, resp.StatusCode())
	})

	t.Run("should return an error if repository returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				return errors.New("unexpected error")
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "A"
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, wantErr, resp.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("should return success", func(t *testing.T) {
		service := enrollment.NewService(l, nil, nil, &mockRepository{
			UpdateMock: func(ctx context.Context, id string, status *string) error {
				assert.Equal(t, "20", id)
				assert.NotNil(t, status)
				assert.Equal(t, "A", *status)
				return nil
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "A"
		resp, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Nil(t, err)

		r := resp.(response.Response)
		assert.Equal(t, http.StatusOK, r.StatusCode())
		assert.Empty(t, r.Error())
		assert.Nil(t, r.GetData())
	})
}
