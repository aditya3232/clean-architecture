package echo

import (
	"clean-architecture/config"
	"clean-architecture/internal/adapter/inbound/echo/response"
	"clean-architecture/internal/domain/entity"
	"clean-architecture/internal/domain/service"
	"clean-architecture/internal/port/inbound"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

type middlewareAdapter struct {
	cfg        *config.Config
	redis      *redis.Client
	jwtService service.JwtServiceInterface
}

func NewMiddlewareAdapter(cfg *config.Config, redis *redis.Client, jwtService service.JwtServiceInterface) inbound.MiddlewareAdapterInterface {
	return &middlewareAdapter{
		cfg:        cfg,
		redis:      redis,
		jwtService: jwtService,
	}
}

func (m *middlewareAdapter) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				err := errors.New("no header authorization found")
				return response.RespondWithError(c, http.StatusUnauthorized, "[MiddlewareAdapter-1] CheckToken", err)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			_, err := m.jwtService.ValidateToken(tokenString)
			if err != nil {
				return response.RespondWithError(c, http.StatusUnauthorized, "[MiddlewareAdapter-2] CheckToken", err)
			}

			getSession, err := m.redis.Get(c.Request().Context(), tokenString).Result()
			if err != nil || len(getSession) == 0 {
				log.Errorf("[MiddlewareAdapter-3] CheckToken: %v", err)
				errSessionNotFound := errors.New("session not found")
				return response.RespondWithError(c, http.StatusUnauthorized, "[MiddlewareAdapter-3] CheckToken", errSessionNotFound)
			}

			jwtUserData := entity.JwtUserData{}
			err = json.Unmarshal([]byte(getSession), &jwtUserData)
			if err != nil {
				return response.RespondWithError(c, http.StatusInternalServerError, "[MiddlewareAdapter-4] CheckToken", err)
			}

			path := c.Request().URL.Path
			segments := strings.Split(strings.Trim(path, "/"), "/")
			if jwtUserData.RoleName == "Customer" && segments[0] == "admin" {
				err := errors.New("customer cannot access admin routes")
				return response.RespondWithError(c, http.StatusForbidden, "[MiddlewareAdapter-5] CheckToken", err)
			}

			c.Set("user", getSession)
			return next(c)
		}
	}
}
