// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/request"
	"miniflux.app/http/response/html"
	"miniflux.app/http/route"
	"miniflux.app/logger"
)

// RefreshFeed refresh a subscription and redirect to the feed entries page.
func (c *Controller) RefreshFeed(w http.ResponseWriter, r *http.Request) {
	feedID := request.RouteInt64Param(r, "feedID")
	if err := c.feedHandler.RefreshFeed(request.UserID(r), feedID); err != nil {
		logger.Error("[Controller:RefreshFeed] %v", err)
	}

	html.Redirect(w, r, route.Path(c.router, "feedEntries", "feedID", feedID))
}

// RefreshAllFeeds refresh all feeds in the background for the current user.
func (c *Controller) RefreshAllFeeds(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)
	jobs, err := c.store.NewUserBatch(userID, c.store.CountFeeds(userID))
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	go func() {
		c.pool.Push(jobs)
	}()

	html.Redirect(w, r, route.Path(c.router, "feeds"))
}
