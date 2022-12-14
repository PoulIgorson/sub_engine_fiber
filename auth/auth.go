// Package auth implements interface for auth.
package auth

import (
	"github.com/gofiber/fiber/v2"

	db "github.com/PoulIgorson/sub_engine_fiber/database"
	. "github.com/PoulIgorson/sub_engine_fiber/define"
)

var IgnoreUrls = []string{
	"/", "/login", "/logout", "/registration",
}
var funcCheckUser func(*db.DB, string) bool

// New return handler for auth.
func New(db_ *db.DB, funcCheckUser_ func(*db.DB, string) bool, ignoreUrls ...[]string) fiber.Handler {
	funcCheckUser = funcCheckUser_
	if len(ignoreUrls) > 0 {
		IgnoreUrls = ignoreUrls[0]
	}
	return myNew(db_)
}

func myNew(db_ *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userStr := c.Cookies("userCookie")
		if Contains(IgnoreUrls, c.Path()) || funcCheckUser(db_, userStr) {
			return c.Next()
		}
		return c.Redirect("/login")
	}
}
