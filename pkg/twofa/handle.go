package twofa

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tendant/simple-idm/auth"
)

const (
	ACCESS_TOKEN_NAME  = "accessToken"
	REFRESH_TOKEN_NAME = "refreshToken"
)

type JwtService interface {
	ParseTokenStr(tokenStr string) (*jwt.Token, error)
	CreateAccessToken(claimData interface{}) (auth.IdmToken, error)
	CreateRefreshToken(claimData interface{}) (auth.IdmToken, error)
}

type Handle struct {
	twoFaService *TwoFaService
	jwtService   JwtService
}

func NewHandle(twoFaService *TwoFaService, jwtService JwtService) Handle {
	return Handle{
		twoFaService: twoFaService,
		jwtService:   jwtService,
	}
}

// setTokenCookie sets a cookie with the given token name, value, and expiration
func (h Handle) setTokenCookie(w http.ResponseWriter, tokenName, tokenValue string, expire time.Time) {
	tokenCookie := &http.Cookie{
		Name:     tokenName,
		Path:     "/",
		Value:    tokenValue,
		Expires:  expire,
		HttpOnly: true,                 // Make the cookie HttpOnly
		Secure:   true,                 // Ensure it's sent over HTTPS
		SameSite: http.SameSiteLaxMode, // Prevent CSRF
	}

	http.SetCookie(w, tokenCookie)
}

// Initiate sending 2fa code
// (POST /2fa/send)
func (h Handle) Post2faSend(w http.ResponseWriter, r *http.Request) *Response {
	var resp SuccessResponse

	data := &Post2faSendJSONRequestBody{}
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		return &Response{
			Code: http.StatusBadRequest,
			body: "unable to parse body",
		}
	}

	// FIXME: read the login id from session cookies
	// Get bearer token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return &Response{
			Code: http.StatusUnauthorized,
			body: "Missing or invalid Authorization header",
		}
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate token
	token, err := h.jwtService.ParseTokenStr(tokenStr)
	if err != nil {
		return &Response{
			Code: http.StatusUnauthorized,
			body: "Invalid access token",
		}
	}

	// Get claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Invalid token claims",
		}
	}

	// Extract login_id from custom_claims
	customClaims, ok := claims["custom_claims"].(map[string]interface{})
	if !ok {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Invalid custom claims format",
		}
	}

	loginIdStr, ok := customClaims["login_id"].(string)
	if !ok {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Missing or invalid login_id in token",
		}
	}

	loginId, err := uuid.Parse(loginIdStr)
	if err != nil {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Invalid login_id format in token",
		}
	}

	err = h.twoFaService.SendTwoFaNotification(r.Context(), loginId, data.TwofaType, data.DeliveryOption)
	if err != nil {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "failed to init 2fa: " + err.Error(),
		}
	}

	return Post2faSendJSON200Response(resp)
}

// Authenticate 2fa passcode
// (POST /2fa/validate)
func (h Handle) Post2faValidate(w http.ResponseWriter, r *http.Request) *Response {
	var resp SuccessResponse

	// Get bearer token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return &Response{
			Code: http.StatusUnauthorized,
			body: "Missing or invalid Authorization header",
		}
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate token
	token, err := h.jwtService.ParseTokenStr(tokenStr)
	if err != nil {
		return &Response{
			Code: http.StatusUnauthorized,
			body: "Invalid access token",
		}
	}

	// Get claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Invalid token claims",
		}
	}

	// Extract login_id from custom_claims
	customClaims, ok := claims["custom_claims"].(map[string]interface{})
	if !ok {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Invalid custom claims format",
		}
	}

	loginIdStr, ok := customClaims["login_id"].(string)
	if !ok {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Missing or invalid login_id in token",
		}
	}

	loginId, err := uuid.Parse(loginIdStr)
	if err != nil {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Invalid login_id format in token",
		}
	}

	data := &Post2faValidateJSONRequestBody{}
	err = render.DecodeJSON(r.Body, &data)
	if err != nil {
		return &Response{
			Code: http.StatusBadRequest,
			body: "unable to parse body",
		}
	}

	valid, err := h.twoFaService.Validate2faPasscode(r.Context(), loginId, data.TwofaType, data.Passcode)
	if err != nil {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "failed to validate 2fa: " + err.Error(),
		}
	}

	if !valid {
		return &Response{
			Code: http.StatusBadRequest,
			body: "2fa validation failed",
		}
	}

	// 2FA validation successful, create access and refresh tokens
	// Extract user data from claims to use for token creation
	userData := customClaims

	// Create access token
	accessToken, err := h.jwtService.CreateAccessToken(userData)
	if err != nil {
		slog.Error("Failed to create access token", "userData", userData, "err", err)
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Failed to create access token",
		}
	}

	// Create refresh token
	refreshToken, err := h.jwtService.CreateRefreshToken(userData)
	if err != nil {
		slog.Error("Failed to create refresh token", "userData", userData, "err", err)
		return &Response{
			Code: http.StatusInternalServerError,
			body: "Failed to create refresh token",
		}
	}

	// Set cookies
	h.setTokenCookie(w, ACCESS_TOKEN_NAME, accessToken.Token, accessToken.Expiry)
	h.setTokenCookie(w, REFRESH_TOKEN_NAME, refreshToken.Token, refreshToken.Expiry)

	// Include tokens in response
	resp.Result = "success"

	return Post2faValidateJSON200Response(resp)
}

// Get all enabled 2fas
// (GET /2fa/enabled)
func (h Handle) Get2faEnabled(w http.ResponseWriter, r *http.Request, loginID string) *Response {
	// Get login ID from path parameter
	loginId, err := uuid.Parse(loginID)
	if err != nil {
		return &Response{
			Code: http.StatusBadRequest,
			body: "invalid login id",
		}
	}

	// Find enabled 2FA methods
	twoFAs, err := h.twoFaService.FindEnabledTwoFAs(r.Context(), loginId)
	if err != nil {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "failed to validate 2fa: " + err.Error(),
		}
	}

	return Get2faEnabledJSON200Response(struct {
		N2faMethods []string `json:"2fa_methods,omitempty"`
	}{
		N2faMethods: twoFAs,
	})
}
