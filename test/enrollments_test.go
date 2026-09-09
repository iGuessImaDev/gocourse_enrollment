package test

import (
	"net/http"
	"testing"

	"github.com/iGuessImaDev/gocourse_domain/domain"
	"github.com/iGuessImaDev/gocourse_enrollment/internal/enrollment"
	"github.com/stretchr/testify/assert"
)

type dataResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}

func TestEnrol(t *testing.T) {
	t.Run("should create an enrollment and get it", func(t *testing.T) {
		bodyRequest := enrollment.CreateReq{
			UserID:   "11-test",
			CourseID: "12-test",
		}

		resp := cli.Post("/enrollments", bodyRequest)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		dataCreated := domain.Enrollment{}
		dRespCreated := dataResponse{Data: &dataCreated}
		err := resp.FillUp(&dRespCreated)
		assert.Nil(t, err)

		assert.Equal(t, "success", dRespCreated.Message)
		assert.Equal(t, http.StatusCreated, dRespCreated.Status)

		assert.Equal(t, "11-test", dataCreated.UserID)
		assert.Equal(t, "12-test", dataCreated.CourseID)
		assert.NotEmpty(t, dataCreated.ID)

		resp = cli.Get("/enrollments?user_id=" + dataCreated.UserID + "&course_id=" + dataCreated.CourseID)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var dataGetAll []domain.Enrollment
		dRespGetAll := dataResponse{Data: &dataGetAll}
		err = resp.FillUp(&dRespGetAll)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, dRespGetAll.Status)
		assert.Equal(t, "success", dRespGetAll.Message)
		assert.Equal(t, 1, len(dataGetAll))
		assert.Equal(t, dataCreated.ID, dataGetAll[0].ID)
		assert.Equal(t, dataCreated.UserID, dataGetAll[0].UserID)
		assert.Equal(t, dataCreated.CourseID, dataGetAll[0].CourseID)
		assert.Equal(t, domain.Pending, dataGetAll[0].Status)
	})

	t.Run("update an enrollment", func(t *testing.T) {
		bodyRequest := enrollment.CreateReq{
			UserID:   "22-test",
			CourseID: "23-test",
		}

		resp := cli.Post("/enrollments", bodyRequest)
		dataCreated := domain.Enrollment{}
		dRespCreated := dataResponse{Data: &dataCreated}
		err := resp.FillUp(&dRespCreated)
		assert.Nil(t, err)

		status := "A"
		resp = cli.Patch("/enrollments/"+dataCreated.ID, enrollment.UpdateReq{Status: &status})
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp = cli.Get("/enrollments?user_id=" + dataCreated.UserID + "&course_id=" + dataCreated.CourseID)
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var dataGetAll []domain.Enrollment
		dRespGetAll := dataResponse{Data: &dataGetAll}
		err = resp.FillUp(&dRespGetAll)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, dRespGetAll.Status)
		assert.Equal(t, "success", dRespGetAll.Message)
		assert.Equal(t, 1, len(dataGetAll))
		assert.Equal(t, dataCreated.ID, dataGetAll[0].ID)
		assert.Equal(t, dataCreated.UserID, dataGetAll[0].UserID)
		assert.Equal(t, dataCreated.CourseID, dataGetAll[0].CourseID)
		assert.Equal(t, domain.Active, dataGetAll[0].Status)
	})
}
