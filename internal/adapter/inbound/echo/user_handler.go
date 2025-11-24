package echo

import (
	"clean-architecture/internal/adapter/inbound/echo/request"
	"clean-architecture/internal/adapter/inbound/echo/response"
	"clean-architecture/internal/domain/entity"
	"clean-architecture/internal/domain/service"
	"clean-architecture/internal/port/inbound"
	"clean-architecture/utils/conv"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type userHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) inbound.UserHandlerInterface {
	return &userHandler{userService: userService}
}

func (u *userHandler) DeleteCustomer(c echo.Context) error {
	var (
		resp = response.DefaultResponseWithPaginations{}
		ctx  = c.Request().Context()
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] DeleteCustomer", err)
	}

	idParamStr := c.Param("id")
	if idParamStr == "" {
		err := errors.New("missing or invalid customer ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] DeleteCustomer", err)
	}

	id, err := conv.StringToInt64(idParamStr)
	if err != nil {
		err := errors.New("invalid customer ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-3] DeleteCustomer", err)
	}

	err = u.userService.DeleteCustomer(ctx, id)
	if err != nil {
		log.Infof("[UserHandler-4] DeleteCustomer: %v", err)
		if err.Error() == "404" {
			errNotFOund := errors.New("customer not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-4] DeleteCustomer", errNotFOund)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-4] DeleteCustomer", err)
	}

	resp.Message = "Customer deleted successfully"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

func (u *userHandler) UpdateCustomer(c echo.Context) error {
	var (
		resp = response.DefaultResponseWithPaginations{}
		ctx  = c.Request().Context()
		req  = request.UpdateCustomerRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] UpdateCustomer", err)
	}

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] UpdateCustomer", err)
	}

	if err := c.Validate(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-3] UpdateCustomer", err)
	}

	latString := ""
	lngString := ""
	if req.Lat != 0 {
		latString = strconv.FormatFloat(req.Lat, 'g', -1, 64)
	}

	if req.Lng != 0 {
		lngString = strconv.FormatFloat(req.Lng, 'g', -1, 64)
	}

	idParamStr := c.Param("id")
	if idParamStr == "" {
		err := errors.New("missing or invalid customer ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-4] UpdateCustomer", err)
	}

	id, err := conv.StringToInt64(idParamStr)
	if err != nil {
		err := errors.New("invalid customer ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-5] UpdateCustomer", err)
	}

	reqEntity := entity.UserEntity{
		ID:      id,
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
		Lat:     latString,
		Lng:     lngString,
		Photo:   req.Photo,
	}

	err = u.userService.UpdateDataUser(ctx, reqEntity)
	if err != nil {
		log.Errorf("[UserHandler-6] UpdateCustomer: %v", err)
		if err.Error() == "404" {
			errNotFound := errors.New("customer not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-6] UpdateCustomer", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-6] UpdateCustomer", err)

	}

	resp.Message = "Success"
	resp.Data = nil

	return c.JSON(http.StatusOK, resp)
}

func (u *userHandler) CreateCustomer(c echo.Context) error {
	var (
		resp = response.DefaultResponseWithPaginations{}
		ctx  = c.Request().Context()
		req  = request.CustomerRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] CreateCustomer", err)
	}

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] CreateCustomer", err)
	}

	if err := c.Validate(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-3] CreateCustomer", err)
	}

	if req.Password != req.PasswordConfirmation {
		err := errors.New("password and confirm password does not match")
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-4] CreateCustomer", err)
	}

	latString := strconv.FormatFloat(req.Lat, 'g', -1, 64)
	lngString := strconv.FormatFloat(req.Lng, 'g', -1, 64)

	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Address:  req.Address,
		Lat:      latString,
		Lng:      lngString,
		Photo:    req.Photo,
		RoleID:   req.RoleID,
	}

	err := u.userService.CreateCustomer(ctx, reqEntity)
	if err != nil {
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-5] CreateCustomer", err)
	}

	resp.Message = "success"
	resp.Data = nil
	resp.Pagination = nil

	return c.JSON(http.StatusCreated, resp)
}

func (u *userHandler) GetCustomerByID(c echo.Context) error {
	var (
		resp     = response.DefaultResponseWithPaginations{}
		ctx      = c.Request().Context()
		respUser = response.CustomerResponse{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] GetCustomerByID", err)
	}

	idParam := c.Param("id")
	if idParam == "" {
		err := errors.New("id invalid")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] GetCustomerByID", err)
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-3] GetCustomerByID", err)
	}

	result, err := u.userService.GetCustomerByID(ctx, id)
	if err != nil {
		log.Errorf("[UserHandler-4] GetCustomerByID: %v", err)
		if err.Error() == "404" {
			errNotFound := errors.New("customer not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-4] GetCustomerByID", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-4] GetCustomerByID", err)
	}

	resp.Message = "success get customer by id"
	respUser.ID = result.ID
	respUser.RoleID = result.RoleID
	respUser.Name = result.Name
	respUser.Email = result.Email
	respUser.Phone = result.Phone
	respUser.Address = result.Address
	respUser.Photo = result.Photo
	respUser.Lat = result.Lat
	respUser.Lng = result.Lng

	resp.Data = respUser
	resp.Pagination = nil

	return c.JSON(http.StatusOK, resp)
}

func (u *userHandler) GetCustomerAll(c echo.Context) error {
	var (
		resp     = response.DefaultResponseWithPaginations{}
		ctx      = c.Request().Context()
		respUser = []response.CustomerListResponse{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not found")
		return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-1] GetCustomerAll", err)

	}

	search := c.QueryParam("search")
	orderBy := "created_at"
	if c.QueryParam("order_by") != "" {
		orderBy = c.QueryParam("order_by")
	}

	orderType := c.QueryParam("order_type")
	if orderType != "asc" && orderType != "desc" {
		orderType = "desc"
	}

	pageStr := c.QueryParam("page")
	var page int64 = 1
	if pageStr != "" {
		page, _ = conv.StringToInt64(pageStr)
		if page <= 0 {
			page = 1
		}
	}

	limitStr := c.QueryParam("limit")
	var limit int64 = 10
	if limitStr != "" {
		limit, _ = conv.StringToInt64(limitStr)
		if limit <= 0 {
			limit = 10
		}
	}

	reqEntity := entity.QueryStringEntity{
		Search:    search,
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		OrderType: orderType,
	}

	results, countData, totalPages, err := u.userService.GetCustomerAll(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-2] GetCustomerAll", err)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-2] GetCustomerAll", err)
	}

	for _, val := range results {
		respUser = append(respUser, response.CustomerListResponse{
			ID:    val.ID,
			Name:  val.Name,
			Email: val.Email,
			Photo: val.Photo,
			Phone: val.Phone,
		})
	}

	resp.Message = "Data retrieved successfully"
	resp.Data = respUser
	resp.Pagination = &response.Pagination{
		Page:       page,
		TotalCount: countData,
		Limit:      limit, // PerPage (sebelumnya)
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, resp)
}

func (u *userHandler) UpdateDataUser(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		req         = request.UpdateDataUserRequest{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not found")
		return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-1] UpdateDataUser", err)
	}

	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] UpdateDataUser", err)
	}

	userID := jwtUserData.UserID

	if err = c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-3] UpdateDataUser", err)
	}

	if err = c.Validate(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-4] UpdateDataUser", err)
	}

	reqEntity := entity.UserEntity{
		ID:      userID,
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		Lat:     req.Lat,
		Lng:     req.Lng,
		Phone:   req.Phone,
		Photo:   req.Photo,
	}

	err = u.userService.UpdateDataUser(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			errNotFound := errors.New("user not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-5] UpdateDataUser", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-5] UpdateDataUser", err)
	}

	resp.Message = "Success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

func (u *userHandler) GetProfileUser(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		respProfile = response.ProfileResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		err := errors.New("data token not found")
		return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-1] GetProfileUser", err)
	}

	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] GetProfileUser", err)
	}

	userID := jwtUserData.UserID

	dataUser, err := u.userService.GetProfileUser(ctx, userID)
	if err != nil {
		if err.Error() == "404" {
			errNotFound := errors.New("user not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-3] GetProfileUser", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-3] GetProfileUser", err)
	}

	respProfile.Address = dataUser.Address
	respProfile.Name = dataUser.Name
	respProfile.Email = dataUser.Email
	respProfile.ID = dataUser.ID
	respProfile.Lat = dataUser.Lat
	respProfile.Lng = dataUser.Lng
	respProfile.Phone = dataUser.Phone
	respProfile.Photo = dataUser.Photo
	respProfile.RoleName = dataUser.RoleName

	resp.Message = "success"
	resp.Data = respProfile

	return c.JSON(http.StatusOK, resp)
}

func (u *userHandler) UpdatePassword(c echo.Context) error {
	var (
		resp = response.DefaultResponse{}
		req  = request.UpdatePasswordRequest{}
		ctx  = c.Request().Context()
	)

	tokenString := c.QueryParam("token")
	if tokenString == "" {
		err := errors.New("missing or invalid token")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] UpdatePassword", err)
	}

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] UpdatePassword", err)
	}

	if err := c.Validate(req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-3] UpdatePassword", err)
	}

	if req.NewPassword != req.ConfirmPassword {
		err := errors.New("new password and confirm password does not match")
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-4] UpdatePassword", err)
	}

	reqEntity := entity.UserEntity{
		Password: req.NewPassword,
		Token:    tokenString,
	}

	err := u.userService.UpdatePassword(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			errNotFound := errors.New("user not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-5] UpdatePassword", errNotFound)
		}

		if err.Error() == "401" {
			errUnauthorized := errors.New("token expired or invalid")
			return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-5] UpdatePassword", errUnauthorized)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-5] UpdatePassword", err)
	}

	resp.Data = nil
	resp.Message = "Password updated successfully"

	return c.JSON(http.StatusOK, resp)
}

func (u *userHandler) VerifyAccount(c echo.Context) error {
	var (
		resp       = response.DefaultResponse{}
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	tokenString := c.QueryParam("token")
	if tokenString == "" {
		err := errors.New("missing or invalid token")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] VerifyAccount", err)
	}

	user, err := u.userService.VerifyToken(ctx, tokenString)
	if err != nil {
		if err.Error() == "404" {
			errNotFound := errors.New("user not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-2] VerifyAccount", errNotFound)
		}

		if err.Error() == "401" {
			errUnauthorized := errors.New("token expired or invalid")
			return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-2] VerifyAccount", errUnauthorized)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-2] VerifyAccount", err)
	}

	respSignIn.ID = user.ID
	respSignIn.Name = user.Name
	respSignIn.Email = user.Email
	respSignIn.Role = user.RoleName
	respSignIn.Lat = user.Lat
	respSignIn.Lng = user.Lng
	respSignIn.Phone = user.Phone
	respSignIn.AccessToken = user.Token

	resp.Message = "Success"
	resp.Data = respSignIn

	return c.JSON(http.StatusOK, resp)
}

func (u *userHandler) ForgotPassword(c echo.Context) error {
	var (
		req  = request.ForgotPasswordRequest{}
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-1] ForgotPassword", err)
	}

	if err := c.Validate(req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-2] ForgotPassword", err)
	}

	reqEntity := entity.UserEntity{
		Email: req.Email,
	}

	err := u.userService.ForgotPassword(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			errNotFound := errors.New("user not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-3] ForgotPassword", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-3] ForgotPassword", err)
	}

	resp.Message = "Success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

func (u *userHandler) CreateUserAccount(c echo.Context) error {
	var (
		req  = request.SignUpRequest{}
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-1] CreateUserAccount", err)
	}

	if err := c.Validate(req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-2] CreateUserAccount", err)
	}

	if req.Password != req.PasswordConfirmation {
		err := errors.New("passwords do not match")
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-3] CreateUserAccount", err)
	}

	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := u.userService.CreateUserAccount(ctx, reqEntity)
	if err != nil {
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-4] CreateUserAccount", err)
	}

	resp.Message = "Success"
	return c.JSON(http.StatusCreated, resp)
}

func (u *userHandler) SignIn(c echo.Context) error {
	var (
		req        = request.SignInRequest{}
		resp       = response.DefaultResponse{}
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-1] SignIn", err)
	}

	if err := c.Validate(req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-2] SignIn", err)
	}

	reqEntity := entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}
	user, token, err := u.userService.SignIn(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-3] SignIn", err)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-4] SignIn", err)
	}

	respSignIn.ID = user.ID
	respSignIn.Name = user.Name
	respSignIn.Email = user.Email
	respSignIn.Role = user.RoleName
	respSignIn.Lat = user.Lat
	respSignIn.Lng = user.Lng
	respSignIn.Phone = user.Phone
	respSignIn.AccessToken = token

	resp.Message = "Success"
	resp.Data = respSignIn

	return c.JSON(http.StatusOK, resp)
}
