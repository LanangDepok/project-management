package routes

import (
	"github.com/LanangDepok/project-management/controllers"
	_ "github.com/LanangDepok/project-management/docs"
	"github.com/LanangDepok/project-management/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/swaggo/swag"
)

func Setup(app *fiber.App,
	uc *controllers.UserController,
	bc *controllers.BoardController,
) {
	// Swagger UI
	app.Get("/swagger", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(swaggerHTML())
	})
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		doc, err := swag.ReadDoc()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		c.Set("Content-Type", "application/json")
		return c.SendString(doc)
	})

	// Public auth routes
	auth := app.Group("/v1/auth")
	auth.Post("/register", uc.Register)
	auth.Post("/login", uc.Login)

	// Protected API routes
	api := app.Group("/api/v1", middleware.JWTProtected())

	users := api.Group("/users")
	users.Get("/page", uc.GetUserPagination)
	users.Get("/:id", uc.GetUser)
	users.Put("/:id", uc.UpdateUser)
	users.Delete("/:id", uc.DeleteUser)

	boards := api.Group("/boards")
	boards.Post("/", bc.CreateBoard)
	boards.Put("/:id", bc.UpdateBoard)
	boards.Post("/:id/members", bc.AddBoardMembers)
	boards.Delete("/:id/members", bc.RemoveBoardMembers)
}

func swaggerHTML() string {
	return `<!DOCTYPE html>
<html>
<head>
  <title>Project Management API</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/swagger-ui.min.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/swagger-ui-bundle.min.js"></script>
<script>
  SwaggerUIBundle({
    url: "/swagger/doc.json",
    dom_id: '#swagger-ui',
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
    layout: "BaseLayout",
    deepLinking: true,
    persistAuthorization: true,
  })
</script>
</body>
</html>`
}
