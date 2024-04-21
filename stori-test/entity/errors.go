package entity

import "errors"

var (
	ErrorBody    = errors.New("error loadind request body; %s")
	ErrEmail     = errors.New("email field is required: %s")
	ErrTrxnSave  = errors.New("transactions dont saved: %s")
	ErrCSV       = errors.New("error loading csv: %s")
	ErrTrxsValid = errors.New("error vaidation csv transactions: %s")
)
