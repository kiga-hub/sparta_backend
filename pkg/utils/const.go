package utils

// Code -
type Code int

// Msg -
type Msg string

const ErrorMsg = "error: %v"

// Error -
const (
	// Success -
	Success Code = 200
	// SuccessMsg -
	SuccessMsg Msg = "OK"
	// ErrInvalidRequestParamsCode -
	ErrInvalidRequestParamsCode Code = 100100
	// ErrInvalidRequestErrMsg -
	ErrInvalidRequestErrMsg Msg = "invalid request params"
	// ErrInternalServerCode -
	ErrInternalServerCode Code = 100200
	// ErrInternalServerMsg -
	ErrInternalServerMsg Msg = "internal server error"
	// ErrGetDataCode -
	ErrGetDataCode Code = 100300
	// ErrGetDataMsg -
	ErrGetDataMsg Msg = "get data error"
	// ErrEmptyDataCode -
	ErrEmptyDataCode Code = 100400
	// ErrEmptyDataMsg -
	ErrEmptyDataMsg Msg = "empty data"
	// ErrFileOperationCode -
	ErrFileOperationCode Code = 100500
	// ErrFileOperationMsg -
	ErrFileOperationMsg Msg = "file operation error"
	// ErrIOCode io error
	ErrIOCode Code = 100600
	// ErrIOMsg -
	ErrIOMsg Msg = "io error"
	// ErrParseCode -
	ErrParseCode Code = 100700
	// ErrCodeParseMsg -
	ErrCodeParseMsg Msg = "parse error"
	// ErrImportExportCode import error
	ErrImportExportCode Code = 300600
	// ErrImportExportMsg -
	ErrImportExportMsg Msg = "import error"
)

// socket - error
const (
	// ErrSocketConnectFailCode -
	ErrSocketConnectFailCode Code = 400100
	// ErrSocketConnectFailMsg -
	ErrSocketConnectFailMsg Msg = "socket connect error"
	// ErrSocketRWFailCode Socket
	ErrSocketRWFailCode Code = 400200
	// ErrSocketRWFailMsg -
	ErrSocketRWFailMsg Msg = "socket rw error"
)

// storage -
const (
	// UploadDirName uiload dir
	UploadDirName = "upload"
)
