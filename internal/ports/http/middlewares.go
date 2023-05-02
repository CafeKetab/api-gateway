package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"go.uber.org/zap"
)

func (s *Server) optionalAuthentication(c *fiber.Ctx) error {
	headerBytes := c.Request().Header.Peek("Authorization")
	header := strings.TrimPrefix(string(headerBytes), "Bearer ")

	if len(header) == 0 {
		return c.Next()
	}

	id, err := s.auth.Authenticate(c.Context(), header)
	if err != nil {
		s.logger.Error("Invalid token header", zap.Error(err))
		return c.Next()
	}

	c.Request().Header.Add("X-User-Id", strconv.FormatUint(id, 10))
	c.Request().Header.Del("Authorization")

	return c.Next()
}

// requiredAuthentication will extract the token and put the user-information
// to the request header (before it being redirected)
func (s *Server) requiredAuthentication(c *fiber.Ctx) error {
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

	c.Request().Header.Add("X-User-Id", strconv.FormatUint(id, 10))
	c.Request().Header.Del("Authorization")

	return c.Next()
}

func (s *Server) proxy(c *fiber.Ctx) error {
	path := strings.TrimPrefix(string(c.Request().URI().Path()), "/v1/")

	constructProxyURL := func(endpoint, base string) string {
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

	var location string

	if endpoint := "users"; strings.HasPrefix(path, endpoint) {
		location = constructProxyURL(endpoint, s.config.TargetUrls.Users)
	} else if endpoint = "books"; strings.HasPrefix(path, endpoint) {
		location = constructProxyURL(endpoint, s.config.TargetUrls.Books)
	} else {
		s.logger.Error("Invalid endpoint", zap.ByteString("URI", c.Request().URI().FullURI()))
		return c.Status(http.StatusNotFound).SendString("The requested endpoint doesn't found")
	}

	return proxy.Do(c, location)
}
