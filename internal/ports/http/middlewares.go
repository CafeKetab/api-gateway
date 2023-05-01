package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// authenticate will extract the token and put the user-information
// to the request header (before it being redirected)
func (s *Server) authenticate(c *fiber.Ctx) error {
	header := c.Request().Header.Peek("Authorization")
	if len(header) == 0 {
		s.logger.Error("Missing authorization header")
		response := "please provide your authentication information"
		return c.Status(http.StatusUnauthorized).SendString(response)
	}

	id, err := s.auth.Authenticate(c.Context(), string(header))
	if err != nil {
		s.logger.Error("Invalid token header", zap.Error(err))
		response := "invalid token header, please login again"
		return c.Status(http.StatusUnauthorized).SendString(response)
	}

	c.Request().Header.Del("Authorization")
	c.Request().Header.Add("X-User-Id", strconv.FormatUint(id, 10))

	return c.Next()
}

// rate adds a new rate for a ride
func (s *Server) redirect(c *fiber.Ctx) error {
	path := strings.TrimPrefix(string(c.Request().URI().Path()), "/v1/")

	constructRedirectURL := func(endpoint, base string) string {
		path = strings.TrimPrefix(path, endpoint)

		if len(path) > 1 {
			path = strings.TrimSuffix(path, "/")
			base += path
		}

		if query := string(c.Request().URI().QueryString()); len(query) > 0 {
			base += fmt.Sprintf("?%s", query)
		}

		return base
	}

	if endpoint := "users"; strings.HasPrefix(path, endpoint) {
		location := constructRedirectURL(endpoint, s.config.TargetUrls.Users)
		return c.Redirect(location, http.StatusMovedPermanently)
	} else if endpoint = "books"; strings.HasPrefix(path, endpoint) {
		location := constructRedirectURL(endpoint, s.config.TargetUrls.Books)
		return c.Redirect(location, http.StatusMovedPermanently)
	}

	s.logger.Error("Invalid endpoint", zap.ByteString("URI", c.Request().URI().FullURI()))
	return c.Status(http.StatusNotFound).SendString("The requested endpoint doesn't found")
}
