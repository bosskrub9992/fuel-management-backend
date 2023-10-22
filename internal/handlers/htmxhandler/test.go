package htmxhandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *HTMXHandler) Test(c echo.Context) error {
	return c.Render(http.StatusOK, "test", "")
}

func (h *HTMXHandler) Example(c echo.Context) error {
	return c.HTML(
		http.StatusOK,
		fmt.Sprintf(
			`<input id="input" type="text" value="%s" required>

			<div id="inputInDiv">
				<label for="inputInDiv">inputInDiv</label>
				<input type="text" value="%s" required>
			</div>
			`,
			uuid.NewString(),
			uuid.NewString(),
		),
	)
}

func (h *HTMXHandler) TestNav(c echo.Context) error {
	return c.Render(http.StatusOK, "test_nav", nil)
}
