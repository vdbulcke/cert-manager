package data

import "fmt"

// DBObjectNotFound error when object not found in DB
type DBObjectNotFound struct {
	Err error
	ID  string
}

func (e *DBObjectNotFound) Error() string {
	return fmt.Sprintf("Object not Found, id=%s", e.ID)
}

// DBObjectAlreadyExist error when object already existing in DB
type DBObjectAlreadyExist struct {
	Err error
	ID  string
}

func (e *DBObjectAlreadyExist) Error() string {
	return fmt.Sprintf("Error ObjectAlreadyExist id='%s'", e.ID)
}

// DBObjectValidationError error when object is not valid
type DBObjectValidationError struct {
	Err error
	Msg string
}

func (e *DBObjectValidationError) Error() string {
	return fmt.Sprintf("ObjectValidationError msg='%s'", e.Msg)
}
