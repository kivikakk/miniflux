package ui  // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/cookie"
	"miniflux.app/http/request"
	"miniflux.app/http/response/html"
	"miniflux.app/http/route"
	"miniflux.app/logger"
	"miniflux.app/ui/form"
	"miniflux.app/ui/session"
	"miniflux.app/ui/view"
)

// CheckLogin validates the username/password and redirects the user to the unread page.
func (c *Controller) CheckLogin(w http.ResponseWriter, r *http.Request) {
	clientIP := request.ClientIP(r)
	sess := session.New(c.store, request.SessionID(r))
	authForm := form.NewAuthForm(r)

	view := view.New(c.tpl, r, sess)
	view.Set("errorMessage", "error.bad_credentials")
	view.Set("form", authForm)

	if err := authForm.Validate(); err != nil {
		logger.Error("[Controller:CheckLogin] %v", err)
		html.OK(w, r, view.Render("login"))
		return
	}

	if err := c.store.CheckPassword(authForm.Username, authForm.Password); err != nil {
		logger.Error("[Controller:CheckLogin] [ClientIP=%s] %v", clientIP, err)
		html.OK(w, r, view.Render("login"))
		return
	}

	sessionToken, userID, err := c.store.CreateUserSession(authForm.Username, r.UserAgent(), clientIP)
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	logger.Info("[Controller:CheckLogin] username=%s just logged in", authForm.Username)
	c.store.SetLastLogin(userID)

	user, err := c.store.UserByID(userID)
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	sess.SetLanguage(user.Language)
	sess.SetTheme(user.Theme)

	http.SetCookie(w, cookie.New(
		cookie.CookieUserSessionID,
		sessionToken,
		c.cfg.IsHTTPS,
		c.cfg.BasePath(),
	))

	html.Redirect(w, r, route.Path(c.router, "unread"))
}
