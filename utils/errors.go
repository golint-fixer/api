package utils

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

type Error struct {
	Success  bool        `json:"success"`
	Code     int         `json:"code,omitempty"`
	Location string      `json:"location,omitempty"`
	Error    interface{} `json:"error"`
	Severe   bool        `json:"-"`
}

func NewError(code int, input interface{}, severe bool) *Error {
	_, file, line, _ := runtime.Caller(1)

	return &Error{
		Code:     code,
		Error:    input,
		Location: filepath.Base(file) + ":" + strconv.Itoa(line),
		Severe:   severe,
	}
}

func (e *Error) String() string {
	buf := &bytes.Buffer{}

	if e.Code != 0 {
		buf.WriteString("[")
		buf.WriteString(strconv.Itoa(e.Code))
		buf.WriteString("] ")
	}

	if e.Location != "" {
		buf.WriteString(e.Location)
		buf.WriteString(": ")
	}

	switch v := e.Error.(type) {
	case error:
		buf.WriteString(v.Error())
	case string:
		buf.WriteString(v)
	default:
		buf.WriteString(fmt.Sprintf("%+v", v))
	}

	return buf.String()
}

const (
	AccountsListUnknown = 10000 + iota
)

const (
	AccountsCreateUnknown = 10100 + iota
	AccountsCreateInvalidInput
	AccountsCreateUnknownStep
	AccountsCreateInvalidLength
	AccountsCreateUsernameTaken
	AccountsCreateEmailUsed
	AccountsCreateUnableToInsertAccount
	AccountsCreateUserNotFound
	AccountsCreateInvalidToken
	AccountsCreateInvalidTokenOwner
	AccountsCreateInvalidTokenType
	AccountsCreateExpiredToken
	AccountsCreateAlreadyConfigured
	AccountsCreateWeakPassword
	AccountsCreateUnableToHash
	AccountsCreateUnableToPrepareLabels
	AccountsCreateUnableToCreateAddress
	AccountsCreateUnableToUpdateAccount
)

const (
	AccountsGetUnknown = 10200 + iota
	AccountsGetOnlyMe
	AccountsGetUnableToGet
)

const (
	AccountsUpdateUnknown = 10300 + iota
	AccountsUpdateInvalidInput
	AccountsUpdateOnlyMe
	AccountsUpdateUnableToGet
	AccountsUpdateInvalidCurrentPassword
	AccountsUpdateWeakPassword
	AccountsUpdateUnableToHash
	AccountsUpdateInvalidPublicKey
	AccountsUpdateInvalidPublicKeyOwner
	AccountsUpdateUnableToUpdate
)

const (
	AccountsDeleteUnknown = 10400 + iota
	AccountsDeleteOnlyMe
	AccountsDeleteUnableToGet
	AccountsDeleteUnableToDelete
)

const (
	AccountsWipeDataUnknown = 10500 + iota
	AccountsWipeDataOnlyMe
	AccountsWipeDataUnableToGet
	AccountsWipeDataUnableToDelete
)

const (
	AccountsStartOnboardingUnknown = 10600 + iota
	AccountsStartOnboardingOnlyMe
	AccountsStartOnboardingUnableToGet
	AccountsStartOnboardingMisconfigured
	AccountsStartOnboardingUnableToInit
)

const (
	AddressesListUnknown = 11000 + iota
	AddressesListUnableToGet
)

const (
	AvatarsUnknown = 12000 + iota
	AvatarsInvalidWidth
)

const (
	ContactsListUnknown = 13000 + iota
	ContactsListUnableToGet
)

const (
	ContactsCreateUnknown = 13100 + iota
	ContactsCreateInvalidInput
	ContactsCreateUnableToInsert
)

const (
	ContactsGetUnknown = 13200 + iota
	ContactsGetUnableToGet
	ContactsGetNotOwned
)

const (
	ContactsUpdateUnknown = 13300 + iota
	ContactsUpdateInvalidInput
	ContactsUpdateUnableToGet
	ContactsUpdateNotOwned
	ContactsUpdateUnableToUpdate
)

const (
	ContactsDeleteUnknown = 13400 + iota
	ContactsDeleteUnableToGet
	ContactsDeleteNotOwned
	ContactsDeleteUnableToDelete
)

const (
	EmailsListUnknown = 14000 + iota
	EmailsListInvalidOffset
	EmailsListInvalidLimit
	EmailsListUnableToGet
	EmailsListUnableToCount
)

const (
	EmailsCreateUnknown = 14100 + iota
	EmailsCreateInvalidInput
	EmailsCreateUnableToFetchFiles
	EmailsCreateFileNotOwned
	EmailsCreateUnableToFetchAccount
	EmailsCreateUnableToFetchLabel
	EmailsCreateInvalidFromAddress
	EmailsCreateUnableToFetchThread
	EmailsCreateThreadNotOwned
	EmailsCreateUnableToUpdateThread
	EmailsCreateUnableToInsertThread
	EmailsCreateUnableToInsertEmail
	EmailsCreateUnableToQueue
)

const (
	EmailsGetUnknown = 14200 + iota
	EmailsGetUnableToGet
	EmailsGetNotOwned
)

const (
	EmailsDeleteUnknown = 14300 + iota
	EmailsDeleteUnableToGet
	EmailsDeleteNotOwned
	EmailsDeleteUnableToDelete
)

const (
	FilesListUnknown = 15000 + iota
)

const (
	FilesCreateUnknown = 15100 + iota
)

const (
	FilesGetUnknown = 15200 + iota
)

const (
	FilesUpdateUnknown = 15300 + iota
)

const (
	FilesDeleteUnknown = 15400 + iota
)

const (
	KeysListUnknown = 16000 + iota
)

const (
	KeysCreateUnknown = 16100 + iota
)

const (
	KeysGetUnknown = 16200 + iota
)

const (
	KeysVoteUnknown = 16300 + iota
)

const (
	LabelsListUnknown = 17000 + iota
)

const (
	LabelsCreateUnknown = 17100 + iota
)

const (
	LabelsGetUnknown = 17200 + iota
)

const (
	LabelsUpdateUnknown = 17300 + iota
)

const (
	LabelsDeleteUnknown = 17400 + iota
)

const (
	MiddlewareUnknown = 18000 + iota
)

const (
	ThreadsListUnknown = 19000 + iota
)

const (
	ThreadsGetUnknown = 19100 + iota
)

const (
	ThreadsUpdateUnknown = 19200 + iota
)

const (
	ThreadsDeleteUnknown = 19300 + iota
)

const (
	TokensGetUnknown = 20000 + iota
)

const (
	TokensCreateUnknown = 20100 + iota
)

const (
	TokensDeleteUnknown = 20200 + iota
)
