package enrollment

import (
	"errors"
)

var ErrUserIDRequired = errors.New("User ID is required")
var ErrCourseIDRequired = errors.New("Course ID is required")
