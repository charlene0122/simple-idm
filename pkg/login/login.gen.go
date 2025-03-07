// Package login provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version v0.3.0 DO NOT EDIT.
package login

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/discord-gophers/goapi-gen/runtime"
	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// DeliveryOption defines model for DeliveryOption.
type DeliveryOption struct {
	DisplayValue string `json:"display_value,omitempty"`
	HashedValue  string `json:"hashed_value,omitempty"`
}

// EmailVerifyRequest defines model for EmailVerifyRequest.
type EmailVerifyRequest struct {
	Email string `json:"email"`
}

// FindUsernameRequest defines model for FindUsernameRequest.
type FindUsernameRequest struct {
	// Email address to find username for
	Email openapi_types.Email `json:"email"`
}

// Login defines model for Login.
type Login struct {
	// Token for 2FA verification if required
	LoginToken *string `json:"loginToken,omitempty"`
	Message    string  `json:"message"`

	// Whether 2FA verification is required
	Requires2fA *bool  `json:"requires2FA,omitempty"`
	Status      string `json:"status"`
	User        User   `json:"user"`

	// List of users associated with the login. Usually contains one user, but may contain multiple if same username is shared.
	Users []User `json:"users,omitempty"`
}

// PasswordPolicyResponse defines model for PasswordPolicyResponse.
type PasswordPolicyResponse struct {
	// Whether common passwords are disallowed
	DisallowCommonPwds *bool `json:"disallow_common_pwds,omitempty"`

	// Number of days until password expires
	ExpirationDays *int `json:"expiration_days,omitempty"`

	// Number of previous passwords to check against
	HistoryCheckCount *int `json:"history_check_count,omitempty"`

	// Maximum number of repeated characters allowed
	MaxRepeatedChars *int `json:"max_repeated_chars,omitempty"`

	// Minimum length of the password
	MinLength *int `json:"min_length,omitempty"`

	// Whether the password requires a digit
	RequireDigit *bool `json:"require_digit,omitempty"`

	// Whether the password requires a lowercase letter
	RequireLowercase *bool `json:"require_lowercase,omitempty"`

	// Whether the password requires a special character
	RequireSpecialChar *bool `json:"require_special_char,omitempty"`

	// Whether the password requires an uppercase letter
	RequireUppercase *bool `json:"require_uppercase,omitempty"`
}

// PasswordReset defines model for PasswordReset.
type PasswordReset struct {
	NewPassword string `json:"new_password"`
	Token       string `json:"token"`
}

// PasswordResetInit defines model for PasswordResetInit.
type PasswordResetInit struct {
	// Username of the account to reset password for
	Username string `json:"username"`
}

// RegisterRequest defines model for RegisterRequest.
type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// SelectUserRequiredResponse defines model for SelectUserRequiredResponse.
type SelectUserRequiredResponse struct {
	Message   string `json:"message,omitempty"`
	Status    string `json:"status,omitempty"`
	TempToken string `json:"temp_token,omitempty"`
	Users     []User `json:"users,omitempty"`
}

// Tokens defines model for Tokens.
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// TwoFactorMethod defines model for TwoFactorMethod.
type TwoFactorMethod struct {
	DeliveryOptions []DeliveryOption `json:"delivery_options,omitempty"`
	Type            string           `json:"type,omitempty"`
}

// TwoFactorRequiredResponse defines model for TwoFactorRequiredResponse.
type TwoFactorRequiredResponse struct {
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`

	// Temporary token to use for 2FA verification
	TempToken        string            `json:"temp_token,omitempty"`
	TwoFactorMethods []TwoFactorMethod `json:"two_factor_methods,omitempty"`
}

// TwoFactorVerify defines model for TwoFactorVerify.
type TwoFactorVerify struct {
	// TOTP code
	Code string `json:"code"`

	// Token from initial login response
	LoginToken string `json:"loginToken"`
}

// User defines model for User.
type User struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Name  string `json:"name"`

	// Whether 2FA is enabled for this user
	TwoFactorEnabled bool `json:"twoFactorEnabled"`
}

// Post2faVerifyJSONBody defines parameters for Post2faVerify.
type Post2faVerifyJSONBody TwoFactorVerify

// PostEmailVerifyJSONBody defines parameters for PostEmailVerify.
type PostEmailVerifyJSONBody EmailVerifyRequest

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody struct {
	// 2FA verification code if enabled
	Code     string `json:"code,omitempty"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// PostMobileLoginJSONBody defines parameters for PostMobileLogin.
type PostMobileLoginJSONBody struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// PostPasswordResetJSONBody defines parameters for PostPasswordReset.
type PostPasswordResetJSONBody PasswordReset

// PostPasswordResetInitJSONBody defines parameters for PostPasswordResetInit.
type PostPasswordResetInitJSONBody PasswordResetInit

// GetPasswordResetPolicyParams defines parameters for GetPasswordResetPolicy.
type GetPasswordResetPolicyParams struct {
	// Password reset token
	Token string `json:"token"`
}

// PostRegisterJSONBody defines parameters for PostRegister.
type PostRegisterJSONBody RegisterRequest

// PostUserSwitchJSONBody defines parameters for PostUserSwitch.
type PostUserSwitchJSONBody struct {
	// ID of the user to switch to
	UserID string `json:"user_id"`
}

// PostUsernameFindJSONBody defines parameters for PostUsernameFind.
type PostUsernameFindJSONBody FindUsernameRequest

// Post2faVerifyJSONRequestBody defines body for Post2faVerify for application/json ContentType.
type Post2faVerifyJSONRequestBody Post2faVerifyJSONBody

// Bind implements render.Binder.
func (Post2faVerifyJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostEmailVerifyJSONRequestBody defines body for PostEmailVerify for application/json ContentType.
type PostEmailVerifyJSONRequestBody PostEmailVerifyJSONBody

// Bind implements render.Binder.
func (PostEmailVerifyJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody PostLoginJSONBody

// Bind implements render.Binder.
func (PostLoginJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostMobileLoginJSONRequestBody defines body for PostMobileLogin for application/json ContentType.
type PostMobileLoginJSONRequestBody PostMobileLoginJSONBody

// Bind implements render.Binder.
func (PostMobileLoginJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostPasswordResetJSONRequestBody defines body for PostPasswordReset for application/json ContentType.
type PostPasswordResetJSONRequestBody PostPasswordResetJSONBody

// Bind implements render.Binder.
func (PostPasswordResetJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostPasswordResetInitJSONRequestBody defines body for PostPasswordResetInit for application/json ContentType.
type PostPasswordResetInitJSONRequestBody PostPasswordResetInitJSONBody

// Bind implements render.Binder.
func (PostPasswordResetInitJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostRegisterJSONRequestBody defines body for PostRegister for application/json ContentType.
type PostRegisterJSONRequestBody PostRegisterJSONBody

// Bind implements render.Binder.
func (PostRegisterJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostUserSwitchJSONRequestBody defines body for PostUserSwitch for application/json ContentType.
type PostUserSwitchJSONRequestBody PostUserSwitchJSONBody

// Bind implements render.Binder.
func (PostUserSwitchJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostUsernameFindJSONRequestBody defines body for PostUsernameFind for application/json ContentType.
type PostUsernameFindJSONRequestBody PostUsernameFindJSONBody

// Bind implements render.Binder.
func (PostUsernameFindJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
// It may also be instantiated directly, for the purpose of responding with a single status code.
type Response struct {
	body        interface{}
	Code        int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.Code)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(code int) *Response {
	resp.Code = code
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// Post2faVerifyJSON200Response is a constructor method for a Post2faVerify response.
// A *Response is returned with the configured status code and content type from the spec.
func Post2faVerifyJSON200Response(body Login) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostEmailVerifyJSON200Response is a constructor method for a PostEmailVerify response.
// A *Response is returned with the configured status code and content type from the spec.
func PostEmailVerifyJSON200Response(body struct {
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostLoginJSON200Response is a constructor method for a PostLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func PostLoginJSON200Response(body Login) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostLoginJSON202Response is a constructor method for a PostLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func PostLoginJSON202Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        202,
		contentType: "application/json",
	}
}

// PostLoginJSON400Response is a constructor method for a PostLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func PostLoginJSON400Response(body struct {
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostLogoutJSON200Response is a constructor method for a PostLogout response.
// A *Response is returned with the configured status code and content type from the spec.
func PostLogoutJSON200Response(body struct {
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostMobileLoginJSON200Response is a constructor method for a PostMobileLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func PostMobileLoginJSON200Response(body struct {
	// JWT access token
	AccessToken string `json:"accessToken"`

	// JWT refresh token
	RefreshToken string `json:"refreshToken"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostPasswordResetJSON200Response is a constructor method for a PostPasswordReset response.
// A *Response is returned with the configured status code and content type from the spec.
func PostPasswordResetJSON200Response(body struct {
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostPasswordResetInitJSON200Response is a constructor method for a PostPasswordResetInit response.
// A *Response is returned with the configured status code and content type from the spec.
func PostPasswordResetInitJSON200Response(body struct {
	Code *string `json:"code,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetPasswordResetPolicyJSON200Response is a constructor method for a GetPasswordResetPolicy response.
// A *Response is returned with the configured status code and content type from the spec.
func GetPasswordResetPolicyJSON200Response(body PasswordPolicyResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostRegisterJSON201Response is a constructor method for a PostRegister response.
// A *Response is returned with the configured status code and content type from the spec.
func PostRegisterJSON201Response(body struct {
	Email *openapi_types.Email `json:"email,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostTokenRefreshJSON200Response is a constructor method for a PostTokenRefresh response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTokenRefreshJSON200Response(body Tokens) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostUserSwitchJSON200Response is a constructor method for a PostUserSwitch response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUserSwitchJSON200Response(body Login) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostUserSwitchJSON400Response is a constructor method for a PostUserSwitch response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUserSwitchJSON400Response(body struct {
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostUserSwitchJSON403Response is a constructor method for a PostUserSwitch response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUserSwitchJSON403Response(body struct {
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        403,
		contentType: "application/json",
	}
}

// PostUsernameFindJSON200Response is a constructor method for a PostUsernameFind response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUsernameFindJSON200Response(body struct {
	// A message indicating the request was processed
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Verify 2FA code during login
	// (POST /2fa/verify)
	Post2faVerify(w http.ResponseWriter, r *http.Request) *Response
	// Verify email address
	// (POST /email/verify)
	PostEmailVerify(w http.ResponseWriter, r *http.Request) *Response
	// Login a user
	// (POST /login)
	PostLogin(w http.ResponseWriter, r *http.Request) *Response
	// Logout user
	// (POST /logout)
	PostLogout(w http.ResponseWriter, r *http.Request) *Response
	// Mobile login endpoint
	// (POST /mobile/login)
	PostMobileLogin(w http.ResponseWriter, r *http.Request) *Response
	// Reset password
	// (POST /password/reset)
	PostPasswordReset(w http.ResponseWriter, r *http.Request) *Response
	// Initiate password reset using username
	// (POST /password/reset/init)
	PostPasswordResetInit(w http.ResponseWriter, r *http.Request) *Response
	// Get password reset policy
	// (GET /password/reset/policy)
	GetPasswordResetPolicy(w http.ResponseWriter, r *http.Request, params GetPasswordResetPolicyParams) *Response
	// Register a new user
	// (POST /register)
	PostRegister(w http.ResponseWriter, r *http.Request) *Response
	// Refresh JWT tokens
	// (POST /token/refresh)
	PostTokenRefresh(w http.ResponseWriter, r *http.Request) *Response
	// Switch to a different user when multiple users are available for the same login
	// (POST /user/switch)
	PostUserSwitch(w http.ResponseWriter, r *http.Request) *Response
	// Send username to user's email address
	// (POST /username/find)
	PostUsernameFind(w http.ResponseWriter, r *http.Request) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// Post2faVerify operation middleware
func (siw *ServerInterfaceWrapper) Post2faVerify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.Post2faVerify(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostEmailVerify operation middleware
func (siw *ServerInterfaceWrapper) PostEmailVerify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostEmailVerify(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostLogin operation middleware
func (siw *ServerInterfaceWrapper) PostLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostLogin(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostLogout operation middleware
func (siw *ServerInterfaceWrapper) PostLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostLogout(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostMobileLogin operation middleware
func (siw *ServerInterfaceWrapper) PostMobileLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostMobileLogin(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostPasswordReset operation middleware
func (siw *ServerInterfaceWrapper) PostPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostPasswordReset(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostPasswordResetInit operation middleware
func (siw *ServerInterfaceWrapper) PostPasswordResetInit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostPasswordResetInit(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetPasswordResetPolicy operation middleware
func (siw *ServerInterfaceWrapper) GetPasswordResetPolicy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPasswordResetPolicyParams

	// ------------- Required query parameter "token" -------------

	if err := runtime.BindQueryParameter("form", true, true, "token", r.URL.Query(), &params.Token); err != nil {
		err = fmt.Errorf("invalid format for parameter token: %w", err)
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{err, "token"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetPasswordResetPolicy(w, r, params)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostRegister operation middleware
func (siw *ServerInterfaceWrapper) PostRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostRegister(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostTokenRefresh operation middleware
func (siw *ServerInterfaceWrapper) PostTokenRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostTokenRefresh(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostUserSwitch operation middleware
func (siw *ServerInterfaceWrapper) PostUserSwitch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostUserSwitch(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostUsernameFind operation middleware
func (siw *ServerInterfaceWrapper) PostUsernameFind(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostUsernameFind(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter %s: %v", err.paramName, err.err)
}

func (err UnescapedCookieParamError) Unwrap() error { return err.err }

type UnmarshalingParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnmarshalingParamError) Error() string {
	return fmt.Sprintf("error unmarshaling parameter %s as JSON: %v", err.paramName, err.err)
}

func (err UnmarshalingParamError) Unwrap() error { return err.err }

type RequiredParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err RequiredParamError) Error() string {
	if err.err == nil {
		return fmt.Sprintf("query parameter %s is required, but not found", err.paramName)
	} else {
		return fmt.Sprintf("query parameter %s is required, but errored: %s", err.paramName, err.err)
	}
}

func (err RequiredParamError) Unwrap() error { return err.err }

type RequiredHeaderError struct {
	paramName string
}

// Error implements error.
func (err RequiredHeaderError) Error() string {
	return fmt.Sprintf("header parameter %s is required, but not found", err.paramName)
}

type InvalidParamFormatError struct {
	err       error
	paramName string
}

// Error implements error.
func (err InvalidParamFormatError) Error() string {
	return fmt.Sprintf("invalid format for parameter %s: %v", err.paramName, err.err)
}

func (err InvalidParamFormatError) Unwrap() error { return err.err }

type TooManyValuesForParamError struct {
	NumValues int
	paramName string
}

// Error implements error.
func (err TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("expected one value for %s, got %d", err.paramName, err.NumValues)
}

// ParameterName is an interface that is implemented by error types that are
// relevant to a specific parameter.
type ParameterError interface {
	error
	// ParamName is the name of the parameter that the error is referring to.
	ParamName() string
}

func (err UnescapedCookieParamError) ParamName() string  { return err.paramName }
func (err UnmarshalingParamError) ParamName() string     { return err.paramName }
func (err RequiredParamError) ParamName() string         { return err.paramName }
func (err RequiredHeaderError) ParamName() string        { return err.paramName }
func (err InvalidParamFormatError) ParamName() string    { return err.paramName }
func (err TooManyValuesForParamError) ParamName() string { return err.paramName }

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:    "/",
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Post("/2fa/verify", wrapper.Post2faVerify)
		r.Post("/email/verify", wrapper.PostEmailVerify)
		r.Post("/login", wrapper.PostLogin)
		r.Post("/logout", wrapper.PostLogout)
		r.Post("/mobile/login", wrapper.PostMobileLogin)
		r.Post("/password/reset", wrapper.PostPasswordReset)
		r.Post("/password/reset/init", wrapper.PostPasswordResetInit)
		r.Get("/password/reset/policy", wrapper.GetPasswordResetPolicy)
		r.Post("/register", wrapper.PostRegister)
		r.Post("/token/refresh", wrapper.PostTokenRefresh)
		r.Post("/user/switch", wrapper.PostUserSwitch)
		r.Post("/username/find", wrapper.PostUsernameFind)
	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"/j8AAP//6exzPfEnAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
