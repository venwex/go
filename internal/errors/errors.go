package errors

import "errors"

var ErrTaskNotFound = errors.New("task not found")
var ErrInvalidID = errors.New("invalid id")
var ErrInvalidTitleName = errors.New("invalid title name")
var ErrMissingId = errors.New("missind id")
var ErrConverting = errors.New("error during converting")

var ErrUserNotFound = errors.New("user not found")
var ErrUserInvalidName = errors.New("invalid user name")
