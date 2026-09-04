package enrollment

import (
	"errors"
	"fmt"
)

var ErrUserIDRequired = errors.New("User ID is required")
var ErrCourseIDRequired = errors.New("Course ID is required")
var ErrStatusRequired = errors.New("status is required")

type ErrNotFound struct {
	enrollmentID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("enrollment '%s' doesn't exist", e.enrollmentID)
}

type ErrInvalidStatus struct {
	Status string
}

func (e ErrInvalidStatus) Error() string {
	return fmt.Sprintf("invalid '%s' status", e.Status)
}
