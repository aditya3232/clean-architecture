package echo

import (
	"clean-architecture/internal/adapter/inbound/echo/request"
	"clean-architecture/internal/adapter/inbound/echo/response"
	"clean-architecture/internal/domain/entity"
	"clean-architecture/internal/domain/service"
	"clean-architecture/internal/port/inbound"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type roleHandler struct {
	roleService service.RoleServiceInterface
}

func NewRoleHandler(roleService service.RoleServiceInterface) inbound.RoleHandlerInterface {
	return &roleHandler{roleService: roleService}
}

func (r *roleHandler) Create(c echo.Context) error {
	var (
		req         = request.RoleRequest{}
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not found")
		return response.RespondWithError(c, http.StatusNotFound, "[RoleHandler-1] Create", err)
	}

	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-2] Create", err)
	}

	if jwtUserData.RoleName != "Super Admin" {
		err := errors.New("only Super Admin can access API role")
		return response.RespondWithError(c, http.StatusForbidden, "[RoleHandler-3] Create", err)
	}

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[RoleHandler-4] Create", err)
	}

	if err := c.Validate(req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[RoleHandler-5] Create", err)
	}

	roleEntity := entity.RoleEntity{
		Name: req.Name,
	}

	err = r.roleService.Create(ctx, roleEntity)
	if err != nil {
		return response.RespondWithError(c, http.StatusInternalServerError, "[RoleHandler-6] Create", err)
	}

	resp.Message = "Success"
	resp.Data = nil
	return c.JSON(http.StatusCreated, resp)
}

func (r *roleHandler) Delete(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not found")
		return response.RespondWithError(c, http.StatusNotFound, "[RoleHandler-1] Delete", err)
	}

	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-2] Delete", err)
	}

	if jwtUserData.RoleName != "Super Admin" {
		err := errors.New("only Super Admin can access API role")
		return response.RespondWithError(c, http.StatusForbidden, "[RoleHandler-3] Delete", err)
	}

	roleIDString := c.Param("id")
	if roleIDString == "" {
		err := errors.New("missing or invalid role ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-4] Delete", err)
	}

	roleID, err := strconv.Atoi(roleIDString)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-5] Delete", err)
	}

	err = r.roleService.Delete(ctx, int64(roleID))
	if err != nil {
		log.Errorf("[RoleHandler-6] Delete: %v", err)
		if err.Error() == "404" {
			errNotFound := errors.New("role not found")
			return response.RespondWithError(c, http.StatusNotFound, "[RoleHandler-6] Delete", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[RoleHandler-6] Delete", err)

	}

	resp.Message = "Role deleted successfully"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

func (r *roleHandler) GetAll(c echo.Context) error {
	var (
		respRole []response.RoleResponse
		resp     = response.DefaultResponse{}
		ctx      = c.Request().Context()
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not found")
		return response.RespondWithError(c, http.StatusNotFound, "[RoleHandler-1] GetAll", err)
	}

	search := c.QueryParam("search")

	roles, err := r.roleService.GetAll(ctx, search)
	if err != nil {
		return response.RespondWithError(c, http.StatusInternalServerError, "[RoleHandler-2] GetAll", err)
	}

	for _, role := range roles {
		respRole = append(respRole, response.RoleResponse{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	resp.Message = "success"
	resp.Data = respRole
	return c.JSON(http.StatusOK, resp)
}

func (r *roleHandler) GetByID(c echo.Context) error {
	var (
		respRole    = response.RoleResponse{}
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not found")
		return response.RespondWithError(c, http.StatusNotFound, "[RoleHandler-1] GetByID", err)
	}

	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-2] GetByID", err)
	}

	if jwtUserData.RoleName != "Super Admin" {
		err := errors.New("only Super Admin can access API role")
		return response.RespondWithError(c, http.StatusForbidden, "[RoleHandler-3] GetByID", err)
	}

	roleIDString := c.Param("id")
	if roleIDString == "" {
		err := errors.New("missing or invalid role ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-4] GetByID", err)
	}

	roleID, err := strconv.Atoi(roleIDString)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-5] GetByID", err)
	}

	role, err := r.roleService.GetByID(ctx, int64(roleID))
	if err != nil {
		log.Errorf("[RoleHandler-6] GetByID: %v", err)
		if err.Error() == "404" {
			errNotFound := errors.New("role not found")
			return response.RespondWithError(c, http.StatusNotFound, "[RoleHandler-6] GetByID", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[RoleHandler-6] GetByID", err)
	}

	respRole.ID = role.ID
	respRole.Name = role.Name
	resp.Message = "success"
	resp.Data = respRole
	return c.JSON(http.StatusOK, resp)
}

func (r *roleHandler) Update(c echo.Context) error {
	var (
		req         = request.RoleRequest{}
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not found")
		return response.RespondWithError(c, http.StatusNotFound, "[RoleHandler-1] Update", err)
	}

	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-2] Update", err)
	}

	if jwtUserData.RoleName != "Super Admin" {
		err := errors.New("only Super Admin can access API role")
		return response.RespondWithError(c, http.StatusForbidden, "[RoleHandler-3] Update", err)
	}

	roleIDString := c.Param("id")
	if roleIDString == "" {
		err := errors.New("missing or invalid role ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-4] Update", err)
	}

	roleID, err := strconv.Atoi(roleIDString)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[RoleHandler-5] Update", err)
	}

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[RoleHandler-6] Update", err)
	}

	if err := c.Validate(req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[RoleHandler-7] Update", err)
	}

	reqEntity := entity.RoleEntity{
		ID:   int64(roleID),
		Name: req.Name,
	}

	err = r.roleService.Update(ctx, reqEntity)
	if err != nil {
		log.Errorf("[RoleHandler-8] Update: %v", err)
		if err.Error() == "404" {
			errNotFound := errors.New("role not found")
			return response.RespondWithError(c, http.StatusNotFound, "[RoleHandler-8] Update", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[RoleHandler-8] Update", err)
	}

	resp.Message = "Role updated successfully"
	resp.Data = nil

	return c.JSON(http.StatusOK, resp)
}
