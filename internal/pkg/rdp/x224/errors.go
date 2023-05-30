package x224

import "errors"

var (
	ErrSmallConnectionConfirmLength = errors.New("small connection confirm length")
	ErrWrongDataLength              = errors.New("wrong data length")
	ErrWrongConnectionConfirmCode   = errors.New("wrong connection confirm code")
)
