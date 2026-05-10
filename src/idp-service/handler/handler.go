package handler

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"lab2/src/idp-service/auth"
	"lab2/src/idp-service/models"
	"lab2/src/idp-service/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//go:embed templates/login.html
var loginTemplate string

//go:embed templates/signup.html
var signupTemplate string

type Handler struct {
	db         *storage.PgStorage
	jwtManager *auth.JWTManager
	issuer     string
	baseURL    string
}

func NewHandler(db *storage.PgStorage, jwtManager *auth.JWTManager, baseURL string) *Handler {
	return &Handler{
		db:         db,
		jwtManager: jwtManager,
		issuer:     baseURL,
		baseURL:    baseURL,
	}
}

func (h *Handler) GetJWKS(c *gin.Context) {
	jwks := h.jwtManager.GetJWKS()
	c.JSON(200, jwks)
}

func (h *Handler) Authorize(c *gin.Context) {
	var req models.AuthorizeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, models.ErrorResponse{
			Error:            "invalid_request",
			ErrorDescription: err.Error(),
		})
		return
	}

	fmt.Printf("[Authorize] Received request:\n")
	fmt.Printf("  ClientId: %s\n", req.ClientId)
	fmt.Printf("  RedirectUri: %s\n", req.RedirectUri)
	fmt.Printf("  ResponseType: %s\n", req.ResponseType)
	fmt.Printf("  Scope: %s\n", req.Scope)
	fmt.Printf("  State: %s\n", req.State)
	fmt.Printf("  Nonce: %s\n", req.Nonce)

	if req.ResponseType != "code" {
		c.JSON(400, models.ErrorResponse{
			Error:            "unsupported_response_type",
			ErrorDescription: "Only authorization_code flow is supported",
		})
		return
	}

	if !auth.Contains(req.Scope, "openid") {
		req.Scope = req.Scope + " openid"
	}

	authRequestData := fmt.Sprintf("%s|%s|%s|%s|%s",
		req.ClientId, req.RedirectUri, req.Scope, req.State, req.Nonce)

	fmt.Printf("[Authorize] Saving auth_request cookie: %s\n", authRequestData)

	c.SetCookie("auth_request", authRequestData, 3600, "/", "", false, false)

	loginHTML := loginTemplate

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, loginHTML)
}

func (h *Handler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(400, models.ErrorResponse{Error: "invalid_request"})
		return
	}

	user, err := h.db.GetUserByUsername(c, username)
	if err != nil {
		c.JSON(401, models.ErrorResponse{Error: "invalid_credentials"})
		return
	}

	if !auth.VerifyPassword(password, user.PasswordHash) {
		c.JSON(401, models.ErrorResponse{Error: "invalid_credentials"})
		return
	}

	authRequestCookie, err := c.Cookie("auth_request")
	if err != nil {
		fmt.Printf("[Login] ERROR: Failed to read auth_request cookie: %v\n", err)
		c.JSON(400, models.ErrorResponse{Error: "invalid_request"})
		return
	}

	fmt.Printf("[Login] Raw auth_request cookie: %s\n", authRequestCookie)

	parts := strings.Split(authRequestCookie, "|")
	fmt.Printf("[Login] Cookie parts count: %d\n", len(parts))
	for i, part := range parts {
		fmt.Printf("[Login]   Part[%d]: %s\n", i, part)
	}

	if len(parts) < 5 {
		fmt.Printf("[Login] ERROR: Expected 5 parts, got %d\n", len(parts))
		c.JSON(400, models.ErrorResponse{Error: "invalid_request"})
		return
	}

	clientId := parts[0]
	redirectUri := parts[1]
	scope := parts[2]
	state := parts[3]
	nonce := parts[4]

	fmt.Printf("[Login] Parsed auth_request cookie:\n")
	fmt.Printf("  clientId: %s\n", clientId)
	fmt.Printf("  redirectUri: %s\n", redirectUri)
	fmt.Printf("  scope: %s\n", scope)
	fmt.Printf("  state: %s\n", state)
	fmt.Printf("  nonce: %s\n", nonce)

	code := auth.GenerateAuthorizationCode()
	authCode := &models.AuthCode{
		Code:        code,
		UserUid:     user.UserUid,
		ClientId:    clientId,
		RedirectUri: redirectUri,
		Scope:       scope,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	if err := h.db.SaveAuthorizationCode(c, authCode); err != nil {
		c.JSON(500, models.ErrorResponse{Error: "server_error"})
		return
	}

	redirectURL := redirectUri + "?code=" + code + "&state=" + state

	fmt.Printf("[Login] Final redirect URL: %s\n", redirectURL)

	c.Redirect(302, redirectURL)
}

func (h *Handler) Signup(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, signupTemplate)
		return
	}

	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	fullname := c.PostForm("fullname")

	if username == "" || email == "" || password == "" || len(password) < 8 {
		c.JSON(400, models.ErrorResponse{Error: "invalid_input"})
		return
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		FullName:     fullname,
		PasswordHash: auth.HashPassword(password),
		Role:         "user",
	}

	_, err := h.db.CreateUser(c, user)
	if err != nil {
		c.JSON(400, models.ErrorResponse{Error: "user_already_exists"})
		return
	}

	c.Redirect(302, "/oauth2/authorize")
}

func (h *Handler) Token(c *gin.Context) {
	var req models.TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, models.ErrorResponse{Error: "invalid_request"})
		return
	}

	var user *models.User
	var err error

	switch req.GrantType {
	case "authorization_code":
		if req.Code == "" || req.RedirectUri == "" {
			c.JSON(400, models.ErrorResponse{Error: "invalid_request"})
			return
		}

		authCode, err := h.db.GetAuthorizationCode(c, req.Code)
		if err != nil {
			c.JSON(400, models.ErrorResponse{
				Error:            "invalid_grant",
				ErrorDescription: err.Error(),
			})
			return
		}

		if authCode.RedirectUri != req.RedirectUri || authCode.ClientId != req.ClientId {
			c.JSON(400, models.ErrorResponse{Error: "invalid_grant"})
			return
		}

		user, err = h.db.GetUserByUid(c, authCode.UserUid)
		if err != nil {
			c.JSON(400, models.ErrorResponse{Error: "invalid_grant"})
			return
		}

		h.db.DeleteAuthorizationCode(c, req.Code)

		req.Scope = authCode.Scope

	case "refresh_token":
		if req.RefreshToken == "" {
			c.JSON(400, models.ErrorResponse{Error: "invalid_request"})
			return
		}

		refreshToken, err := h.db.GetRefreshToken(c, req.RefreshToken)
		if err != nil {
			c.JSON(400, models.ErrorResponse{Error: "invalid_grant"})
			return
		}

		if refreshToken.ClientId != req.ClientId {
			c.JSON(400, models.ErrorResponse{Error: "invalid_grant"})
			return
		}

		user, err = h.db.GetUserByUid(c, refreshToken.UserUid)
		if err != nil {
			c.JSON(400, models.ErrorResponse{Error: "invalid_grant"})
			return
		}

		req.Scope = refreshToken.Scope

	default:
		c.JSON(400, models.ErrorResponse{Error: "unsupported_grant_type"})
		return
	}

	if user == nil {
		c.JSON(400, models.ErrorResponse{Error: "invalid_grant"})
		return
	}

	if !auth.Contains(req.Scope, "openid") {
		req.Scope = "openid " + req.Scope
	}

	accessToken, err := h.jwtManager.GenerateAccessToken(user, req.ClientId, req.Scope)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: "server_error"})
		return
	}

	idToken, err := h.jwtManager.GenerateIdToken(user, req.ClientId, req.Scope, "")
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: "server_error"})
		return
	}

	refreshTokenStr := auth.GenerateRefreshToken()
	refreshToken := &models.RefreshToken{
		Token:     refreshTokenStr,
		UserUid:   user.UserUid,
		ClientId:  req.ClientId,
		Scope:     req.Scope,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := h.db.SaveRefreshToken(c, refreshToken); err != nil {
		c.JSON(500, models.ErrorResponse{Error: "server_error"})
		return
	}

	response := models.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600 * 24 * 30,
		RefreshToken: refreshTokenStr,
		IdToken:      idToken,
		Scope:        req.Scope,
	}

	c.JSON(200, response)
}

func (h *Handler) UserInfo(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := h.jwtManager.VerifyAccessToken(tokenString)
	if err != nil {
		c.JSON(401, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	userInfo := map[string]interface{}{
		"sub": (*claims)["sub"],
	}

	if scope, ok := (*claims)["scope"]; ok {
		scopeStr := scope.(string)
		if auth.Contains(scopeStr, "profile") {
			userInfo["preferred_username"] = (*claims)["preferred_username"]
			userInfo["name"] = (*claims)["name"]
		}
		if auth.Contains(scopeStr, "email") {
			userInfo["email"] = (*claims)["email"]
		}
	}

	c.JSON(200, userInfo)
}

func (h *Handler) CreateUser(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := h.jwtManager.VerifyAccessToken(tokenString)
	if err != nil {
		c.JSON(401, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	if (*claims)["role"] != "admin" {
		c.JSON(403, models.ErrorResponse{Error: "forbidden"})
		return
	}

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, models.ErrorResponse{Error: "invalid_request"})
		return
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: auth.HashPassword(req.Password),
		Role:         "user",
	}

	createdUser, err := h.db.CreateUser(c, user)
	if err != nil {
		c.JSON(400, models.ErrorResponse{Error: "user_already_exists"})
		return
	}

	response := models.UserResponse{
		UserUid:  createdUser.UserUid,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		FullName: createdUser.FullName,
		Role:     createdUser.Role,
	}

	c.JSON(201, response)
}

func (h *Handler) GetMe(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := h.jwtManager.VerifyAccessToken(tokenString)
	if err != nil {
		c.JSON(401, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	userUid, _ := uuid.Parse((*claims)["sub"].(string))
	user, err := h.db.GetUserByUid(c, userUid)
	if err != nil {
		c.JSON(404, models.ErrorResponse{Error: "user_not_found"})
		return
	}

	response := models.UserResponse{
		UserUid:  user.UserUid,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	}

	c.JSON(200, response)
}

func (h *Handler) GetHealth(c *gin.Context) {
	c.JSON(200, gin.H{"status": "UP"})
}
