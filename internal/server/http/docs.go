package http

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Docs(swagger string) fiber.Handler {
	var docs = fmt.Sprintf(`
<!doctype html>
<html>
  <head>
    <title>Scalar API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <!-- Need a Custom Header? Check out this example https://codepen.io/scalarorg/pen/VwOXqam -->
    <script 
		id="api-reference"
		type="application/json"
	>%s</script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>
`, swagger)

	return func(ctx *fiber.Ctx) error {
		ctx.Set(fiber.HeaderContentType, "text/html")
		return ctx.SendStream(strings.NewReader(docs), len(docs))
	}
}
