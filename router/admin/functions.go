// Package functions implements handlers for admin pages.
package functions

import (
	"github.com/gofiber/fiber/v2"

	db "github.com/PoulIgorson/sub_engine_fiber/database"
	user "github.com/PoulIgorson/sub_engine_fiber/database/buckets/user"
	"github.com/PoulIgorson/sub_engine_fiber/types"
)

// IndexPage returns handler for admin index page.
func IndexPage(db_ *db.DB, urls ...interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cuser := c.Context().UserValue("user").(*user.User)
		context := fiber.Map{
			"pagename":   "Админ",
			"menu":       urls[0],
			"admin_menu": urls[1],
			"user":       cuser,
		}
		if c.Method() == "GET" && cuser != nil {
			context["notifies"] = types.Notifies(cuser.ID, true)
		}
		return c.Render("admin/index", context)
	}
}
