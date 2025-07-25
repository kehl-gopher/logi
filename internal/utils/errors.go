package utils

import "errors"

var ErrorTableInsertFailed = errors.New("could not insert into table")
var ErrorEmailAlreadyExists = errors.New("email already exists")
