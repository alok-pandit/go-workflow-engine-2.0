package initiator

import (
	"log"
	"os"
	"time"

	"github.com/alok-pandit/go-workflow-engine-2.0/src/handlers"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/storage/rueidis"
	"github.com/rs/zerolog"
	"github.com/segmentio/encoding/json"
)

func Initialize() {

	app := fiber.New(fiber.Config{
		Prefork:      false,
		ServerHeader: "Fiber",
		AppName:      "workflow-engine",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	cacheDB := rueidis.New(rueidis.Config{
		InitAddress: []string{os.Getenv("REDIS_URL")},
		Username:    "",
		Password:    "",
		SelectDB:    0,
		Reset:       false,
		TLSConfig:   nil,
		CacheTTL:    30 * time.Minute,
	})

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
		Storage:      cacheDB,
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(etag.New(etag.Config{
		Weak: true,
	}))

	app.Use(requestid.New())

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
		Fields: []string{"ip", "port", "latency", "time", "status", "${locals:requestid}", "method", "url", "error"},
	}))

	app.Use(idempotency.New())

	app.Use(recover.New())

	app.Use(helmet.New())

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_ENC_KEY"),
	}))

	app.Use(healthcheck.New())

	go handlers.ExecuteProcess(handlers.Build("src/workflow_definitions/login_workflow.bpmn"))

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))

}
