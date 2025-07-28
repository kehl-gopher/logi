package utils

import "errors"

var ErrorTableInsertFailed = errors.New("could not insert into table")
var ErrorEmailAlreadyExists = errors.New("email already exists")

var ErrorNotFound = errors.New("not found")
var ErrPasswordNotMatch = errors.New("password does not match")
