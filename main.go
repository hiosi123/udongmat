package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hiosi123/udongmat/auth"
	"github.com/hiosi123/udongmat/config"
	"github.com/hiosi123/udongmat/db"
	"github.com/hiosi123/udongmat/handlers"
	"github.com/hiosi123/udongmat/storage"
	"go.uber.org/fx"

	_ "github.com/go-sql-driver/mysql"
)

func newFiberServer(lc fx.Lifecycle, userHandlers *handlers.UserHandler) *fiber.App {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// attach the user handlers
	userGroup := app.Group("/users")
	userGroup.Post("/sign-up", userHandlers.SignUpUser)
	userGroup.Post("/sign-in", userHandlers.SignInUser)
	userGroup.Get("/me", userHandlers.GetUserInfo)
	userGroup.Post("/sign-out", userHandlers.SignOutUser)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// TODO: switch the port to an env variable
			fmt.Println("Starting fiber server on port 8080")
			go app.Listen(":8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	return app
}

func main() {
	fx.New(
		fx.Provide(
			// creates: configEnvVars
			config.LoadEnv,
			// create: *sqlx.DB
			db.CreateMySqlConnection,
			// creates: *storage.UserStorage
			storage.NewUserStorage,
			// creates: *handlers.UserHandler
			handlers.NewUserHandler,
			// creats: *redis.Client
			db.CreateRedisConnection,
			// creates: *auth.SessionManager
			auth.NewSessionManager,
		),
		fx.Invoke(newFiberServer),
	).Run()
}

// func Hello(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "hello world!")
// }

// func Udongmat(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "안녕하세요, 우리동네 맛집을 선정하기 위해 만들어진 웹 사이트 입니다.")
// }

// type Address struct {
// }

// func PostAddress(w http.ResponseWriter, r *http.Request) {
// 	var address Address

// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/hello":
// 		if r.Method == http.MethodGet {
// 			Hello(w, r)
// 		} else {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	case "/udongmat":
// 		if r.Method == http.MethodGet {
// 			Udongmat(w, r)
// 		} else {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	case "/address":
// 		if r.Method == http.MethodPost {
// 			PostAddress(w, r)
// 		}
// 	default:
// 		fmt.Fprintf(w, "health checking...")
// 	}
// }

// func main() {
// 	http.HandleFunc("/", handler)
// 	http.ListenAndServe(":3001", nil)
// }
