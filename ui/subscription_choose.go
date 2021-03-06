// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/client"
	"miniflux.app/http/response/html"
	"miniflux.app/http/request"
	"miniflux.app/http/route"
	"miniflux.app/ui/form"
	"miniflux.app/ui/session"
	"miniflux.app/ui/view"
)

// ChooseSubscription shows a page to choose a subscription.
func (c *Controller) ChooseSubscription(w http.ResponseWriter, r *http.Request) {
	sess := session.New(c.store, request.SessionID(r))
	view := view.New(c.tpl, r, sess)

	user, err := c.store.UserByID(request.UserID(r))
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	categories, err := c.store.Categories(user.ID)
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	view.Set("categories", categories)
	view.Set("menu", "feeds")
	view.Set("user", user)
	view.Set("countUnread", c.store.CountUnreadEntries(user.ID))
	view.Set("countErrorFeeds", c.store.CountErrorFeeds(user.ID))
	view.Set("defaultUserAgent", client.DefaultUserAgent)

	subscriptionForm := form.NewSubscriptionForm(r)
	if err := subscriptionForm.Validate(); err != nil {
		view.Set("form", subscriptionForm)
		view.Set("errorMessage", err.Error())
		html.OK(w, r, view.Render("add_subscription"))
		return
	}

	feed, err := c.feedHandler.CreateFeed(
		user.ID,
		subscriptionForm.CategoryID,
		subscriptionForm.URL,
		subscriptionForm.Crawler,
		subscriptionForm.UserAgent,
		subscriptionForm.Username,
		subscriptionForm.Password,
	)
	if err != nil {
		view.Set("form", subscriptionForm)
		view.Set("errorMessage", err)
		html.OK(w, r, view.Render("add_subscription"))
		return
	}

	html.Redirect(w, r, route.Path(c.router, "feedEntries", "feedID", feed.ID))
}
