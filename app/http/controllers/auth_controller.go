package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"

	"leave/app/helpers"
	authRequest "leave/app/http/requests/auth"
	"leave/app/models"
	"leave/services/basalam"
)

type AuthController struct {
	basalam *basalam.Client
}

func NewAuthController() *AuthController {
	return &AuthController{
		basalam: basalam.NewClient(),
	}
}

// knownLoginErrors are the error codes BasalamRedirect/BasalamCallback may emit;
// each maps to a key under "login.errors.*" in lang/<locale>.json.
var knownLoginErrors = map[string]bool{
	"not_configured": true,
	"state_mismatch": true,
	"missing_code":   true,
	"exchange":       true,
	"whoami":         true,
	"save_user":      true,
	"credentials":    true,
}

// ShowLogin renders the login page (or redirects in if already authenticated).
func (r *AuthController) ShowLogin(ctx http.Context) http.Response {
	if session := ctx.Request().Session(); session != nil {
		if id, ok := session.Get("user_id").(string); ok && id != "" {
			return ctx.Response().Redirect(http.StatusFound, "/leaves-me")
		}
	}

	var errMessage string
	if code := ctx.Request().Query("error"); knownLoginErrors[code] {
		errMessage = facades.Lang(ctx).Get("login.errors." + code)
	}

	var success string
	if ctx.Request().Query("registered") == "1" {
		success = facades.Lang(ctx).Get("login.registered")
	}

	return ctx.Response().View().Make("login.tmpl", map[string]any{
		"title":   facades.Lang(ctx).Get("app.tagline"),
		"error":   errMessage,
		"success": success,
	})
}

// Login authenticates a user by mobile + password (manual login).
func (r *AuthController) Login(ctx http.Context) http.Response {
	mobile := helpers.EnglishDigits(strings.TrimSpace(ctx.Request().Input("mobile")))
	password := ctx.Request().Input("password")
	if mobile == "" || password == "" {
		return ctx.Response().Redirect(http.StatusFound, "/login?error=credentials")
	}

	var user models.User
	_ = facades.Orm().Query().Where("mobile", mobile).First(&user)
	if user.ID == 0 || user.Password == "" {
		return ctx.Response().Redirect(http.StatusFound, "/login?error=credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/login?error=credentials")
	}

	ctx.Request().Session().Put("user_id", strconv.FormatUint(uint64(user.ID), 10))
	return ctx.Response().Redirect(http.StatusFound, "/leaves-me")
}

// ShowRegister renders the public manual user-creation form.
func (r *AuthController) ShowRegister(ctx http.Context) http.Response {
	if session := ctx.Request().Session(); session != nil {
		if id, ok := session.Get("user_id").(string); ok && id != "" {
			return ctx.Response().Redirect(http.StatusFound, "/leaves-me")
		}
	}
	return ctx.Response().View().Make("register.tmpl", map[string]any{
		"title": facades.Lang(ctx).Get("register.title"),
	})
}

// Register creates a user manually (name + mobile). No login is performed;
// the user later signs in via Basalam (linked by mobile).
func (r *AuthController) Register(ctx http.Context) http.Response {
	var req authRequest.RegisterRequest

	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil {
		return r.renderRegisterError(ctx, ctx.Request().Input("name"), ctx.Request().Input("mobile"), debugMessage(ctx, "register.errors.failed", err))
	}
	if errors != nil {
		return r.renderRegisterError(ctx, req.Name, req.Mobile, errors.One())
	}

	// Mobile is always stored with English (ASCII) digits.
	req.Mobile = helpers.EnglishDigits(strings.TrimSpace(req.Mobile))

	var existing models.User
	_ = facades.Orm().Query().Where("mobile", req.Mobile).First(&existing)
	if existing.ID != 0 {
		return r.renderRegisterError(ctx, req.Name, req.Mobile, facades.Lang(ctx).Get("register.errors.mobile_taken"))
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		facades.Log().Error(err)
		return r.renderRegisterError(ctx, req.Name, req.Mobile, debugMessage(ctx, "register.errors.failed", err))
	}

	user := models.User{Name: req.Name, Mobile: req.Mobile, Password: string(hashed), IsActive: true, IsAdmin: noAdminYet()}
	if err := facades.Orm().Query().Create(&user); err != nil {
		facades.Log().Error(err)
		return r.renderRegisterError(ctx, req.Name, req.Mobile, debugMessage(ctx, "register.errors.failed", err))
	}

	return ctx.Response().Redirect(http.StatusFound, "/login?registered=1")
}

func (r *AuthController) renderRegisterError(ctx http.Context, name, mobile, msg string) http.Response {
	return ctx.Response().View().Make("register.tmpl", map[string]any{
		"title":  facades.Lang(ctx).Get("register.title"),
		"error":  msg,
		"name":   name,
		"mobile": mobile,
	})
}

// BasalamRedirect builds the Basalam SSO URL and redirects the user to it.
func (r *AuthController) BasalamRedirect(ctx http.Context) http.Response {
	if !r.basalam.IsConfigured() {
		return ctx.Response().Redirect(http.StatusFound, "/login?error=not_configured")
	}

	state, err := randomState()
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/login?error=state_mismatch")
	}

	ctx.Request().Session().Put("oauth_state", state)

	return ctx.Response().Redirect(http.StatusFound, r.basalam.AuthorizationURL(state))
}

// BasalamCallback handles the OAuth callback: verify state, exchange the code,
// load the identity, upsert the local user and start the session.
func (r *AuthController) BasalamCallback(ctx http.Context) http.Response {
	session := ctx.Request().Session()

	state := ctx.Request().Query("state")
	expected, _ := session.Get("oauth_state").(string)
	session.Forget("oauth_state")
	if state == "" || expected == "" || state != expected {
		return ctx.Response().Redirect(http.StatusFound, "/login?error=state_mismatch")
	}

	code := ctx.Request().Query("code")
	if code == "" {
		return ctx.Response().Redirect(http.StatusFound, "/login?error=missing_code")
	}

	token, err := r.basalam.ExchangeCode(code)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Redirect(http.StatusFound, "/login?error=exchange")
	}

	basalamUser, err := r.basalam.Whoami(token)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Redirect(http.StatusFound, "/login?error=whoami")
	}

	user, err := upsertUser(basalamUser)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Redirect(http.StatusFound, "/login?error=save_user")
	}

	session.Put("user_id", strconv.FormatUint(uint64(user.ID), 10))

	return ctx.Response().Redirect(http.StatusFound, "/leaves-me")
}

// Logout clears the session and returns to the login page.
func (r *AuthController) Logout(ctx http.Context) http.Response {
	if session := ctx.Request().Session(); session != nil {
		session.Forget("user_id")
		_ = session.Invalidate()
	}

	return ctx.Response().Redirect(http.StatusFound, "/login")
}

// upsertUser finds or creates the local user linked to the Basalam account.
func upsertUser(b basalam.BasalamUser) (*models.User, error) {
	var user models.User

	// Primary match: by Basalam user id.
	_ = facades.Orm().Query().Where("basalam_user_id", b.ID).First(&user)
	if user.ID != 0 {
		return &user, nil
	}

	// Secondary match: link an existing user that shares the same mobile.
	if b.Mobile != "" {
		_ = facades.Orm().Query().Where("mobile", b.Mobile).First(&user)
	}

	basalamID := b.ID

	if user.ID != 0 {
		user.BasalamUserId = &basalamID
		if user.Name == "" {
			user.Name = b.Name
		}
		if err := facades.Orm().Query().Save(&user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	user = models.User{
		Name:          b.Name,
		BasalamUserId: &basalamID,
		Mobile:        b.Mobile,
		IsActive:      true,
		IsAdmin:       noAdminYet(),
	}
	if err := facades.Orm().Query().Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// debugMessage returns the real error text when app.debug is on, otherwise the
// translated fallback message. Helps diagnose failures without leaking details
// in production.
func debugMessage(ctx http.Context, langKey string, err error) string {
	if err != nil && facades.Config().GetBool("app.debug") {
		return err.Error()
	}
	return facades.Lang(ctx).Get(langKey)
}

// noAdminYet reports whether the system has no admin user yet (so the first
// user created becomes the admin).
func noAdminYet() bool {
	count, _ := facades.Orm().Query().Model(&models.User{}).Where("is_admin", true).Count()
	return count == 0
}

func randomState() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
