package apperror

// ======================
// Validation
// ======================
const (
	ValidationError        = "VALIDATION_ERROR"
	ValidationRequired     = "VALIDATION_REQUIRED"
	ValidationMinLength    = "VALIDATION_MIN_LENGTH"
	ValidationMaxLength    = "VALIDATION_MAX_LENGTH"
	ValidationInvalidEmail = "VALIDATION_INVALID_EMAIL"
	ValidationInvalidUUID  = "VALIDATION_INVALID_UUID"
)

// HTTP Errors
const (
	BadRequestError = "BAD_REQUEST"
	NotFoundError   = "NOT_FOUND"
	ConflictError   = "CONFLICT"
)

// ======================
// Authentication
// ======================
const (
	AuthUnauthorized    = "AUTH_UNAUTHORIZED"
	AuthInvalidToken    = "AUTH_INVALID_TOKEN"
	AuthExpiredToken    = "AUTH_EXPIRED_TOKEN"
	AuthForbidden       = "AUTH_FORBIDDEN"
	AuthInvalidPassword = "AUTH_INVALID_PASSWORD"
	InvalidCredentials  = "INVALID_CREDENTIALS"
)

// ======================
// User Domain
// ======================
const (
	UserNotFound         = "USER_NOT_FOUND"
	UserAlreadyExists    = "USER_ALREADY_EXISTS"
	UserInactive         = "USER_INACTIVE"
	UserPasswordMismatch = "USER_PASSWORD_MISMATCH"
)

// ======================
// Permission / Role
// ======================
const (
	PermissionDenied = "PERMISSION_DENIED"
)

// ======================
// Data / Repository
// ======================
const (
	DataConflict     = "DATA_CONFLICT"
	DataNotFound     = "DATA_NOT_FOUND"
	DataDuplicate    = "DATA_DUPLICATE"
	DataConstraint   = "DATA_CONSTRAINT"
	ResourceNotFound = "RESOURCE_NOT_FOUND"
)

// ======================
// External / Integration
// ======================
const (
	ExternalServiceError   = "EXTERNAL_SERVICE_ERROR"
	ExternalTimeout        = "EXTERNAL_TIMEOUT"
	ExternalInvalidRequest = "EXTERNAL_INVALID_REQUEST"
)

// ======================
// System
// ======================
const (
	SystemInternalError     = "INTERNAL_SERVER_ERROR"
	SystemTimeout           = "SYSTEM_TIMEOUT"
	SystemUnavailable       = "SYSTEM_UNAVAILABLE"
	UnknownError            = "UNKNOWN_ERROR"
	RateLimitExceeded       = "RATE_LIMIT_EXCEEDED"
	DatabaseError           = "DATABASE_ERROR"
	DatabaseQueryError      = "DATABASE_QUERY_ERROR"
	DatabaseConnectionError = "DATABASE_CONNECTION_ERROR"
	GeneralError            = "GENERAL_ERROR"
)

// ======================
// Password Errors
// ======================
const (
	PasswordTooShort = "PASSWORD_TOO_SHORT"
	PasswordWeak     = "PASSWORD_WEAK"
)
