package utils

import "errors"

var (
	InvalidUuidError       = errors.New("invalid uuid")
	NoCurrencyError        = errors.New("no currency")
	NoUserError            = errors.New("no user")
	NoAccountError         = errors.New("no account")
	NoRateError            = errors.New("no rate")
	LowBalanceError        = errors.New("low balance")
	NoRatesToUploadError   = errors.New("no rates to upload")
	InvalidDateFormatError = errors.New("invalid date format")
	ForbiddenCurrencyError = errors.New("forbidden currency")
	CannotWriteToFileError = errors.New("cannot write to file")
)
