package apperr

import (
	"fmt"
	"net/http"
)

// AppError는 서비스 전반에서 사용하는 커스텀 에러 타입입니다.
type AppError struct {
	Code       string // 내부 에러 코드 (예: "USER_NOT_FOUND")
	Message    string // 클라이언트 전달용 메시지
	StatusCode int    // 매핑할 HTTP Status Code
	Err        error  // 로깅용 원본 에러 (wrapping)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError는 필드명을 지정하지 않고 인자 값만 전달하여 AppError 객체를 생성합니다.
func NewAppError(code, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// 미리 정의된 에러 (Sentinel Errors)
var (
	ErrNotFound     = NewAppError("NOT_FOUND", "요청한 리소스를 찾을 수 없습니다.", http.StatusNotFound, nil)
	ErrAlreadyExist = NewAppError("ALREADY_EXISTS", "이미 존재하는 리소스입니다.", http.StatusConflict, nil)
	ErrInvalidInput = NewAppError("INVALID_INPUT", "잘못된 요청 파라미터입니다.", http.StatusBadRequest, nil)
	ErrUnauthorized = NewAppError("UNAUTHORIZED", "인증이 필요합니다.", http.StatusUnauthorized, nil)
	ErrInternal     = NewAppError("INTERNAL_SERVER_ERROR", "서버 내부 오류가 발생했습니다.", http.StatusInternalServerError, nil)
)
