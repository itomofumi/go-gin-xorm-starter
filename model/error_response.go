package model

import (
	"encoding/json"
	"fmt"
)

// ErrorResponse の定義
type ErrorResponse struct {
	Errors []*ErrorResponseInner `json:"errors"`
}

// ErrorResponseInner の定義
type ErrorResponseInner struct {
	Code     string    `json:"code"`
	Type     ErrorType `json:"type"`
	Messages []string  `json:"messages"`
}

// Append adds an error to ErrorResponse
func (res *ErrorResponse) Append(code string, t ErrorType, messages ...interface{}) {
	errRes := &ErrorResponseInner{
		Messages: msgToStrings(messages),
		Code:     code,
		Type:     t,
	}
	res.Errors = append(res.Errors, errRes)
}

// String is a stringer impl
func (res ErrorResponse) String() string {
	str, _ := json.MarshalIndent(&res, "", "  ")
	return string(str)
}

func msgToStrings(messages []interface{}) []string {
	results := make([]string, 0, len(messages))
	for _, m := range messages {

		switch v := m.(type) {
		case string:
			results = append(results, v)
		case error:
			results = append(results, v.Error())
		case fmt.Stringer:
			results = append(results, v.String())
		default:
			results = append(results, fmt.Sprint(v))
		}
	}
	return results
}

// ErrorType エラータイプ
type ErrorType string

const (
	// ErrorAuth authentication error
	ErrorAuth ErrorType = "AuthError"
	// ErrorUnknown unknown error
	ErrorUnknown ErrorType = "UnknownError"
	// ErrorParam parameter error
	ErrorParam ErrorType = "ParamError"
	// ErrorNotFound not found error
	ErrorNotFound ErrorType = "NotFoundError"
	// ErrorLimitExceeded throttling error
	ErrorLimitExceeded ErrorType = "LimitExceededError"
)

// NewErrorResponse APIエラー時の詳細レスポンスを生成
func NewErrorResponse(code string, t ErrorType, messages ...interface{}) *ErrorResponse {
	res := &ErrorResponse{
		Errors: []*ErrorResponseInner{
			&ErrorResponseInner{
				Messages: msgToStrings(messages),
				Code:     code,
				Type:     t,
			}}}
	return res
}
