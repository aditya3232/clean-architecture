package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	echoinboundadapter "clean-architecture/internal/adapter/inbound/echo"
	"clean-architecture/internal/domain/entity"
	"clean-architecture/tests"
	"clean-architecture/tests/mock"

	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestGetAllRoles_Success(t *testing.T) {
	c, rec := tests.NewEchoContext(http.MethodGet, "/roles", nil)
	c.Set("search", "")
	c.Set("user", "test-user")

	// Setup mock service pakai testify
	mockService := new(mock.MockRoleService)
	mockService.On("GetAll", testifymock.Anything, "").Return([]entity.RoleEntity{
		{ID: 1, Name: "Admin"},
	}, nil)

	// Inject ke handler
	roleHandler := echoinboundadapter.NewRoleHandler(mockService)

	// Jalankan handler
	err := roleHandler.GetAll(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Validasi body
	var body map[string]any
	err = json.Unmarshal(rec.Body.Bytes(), &body)

	// print
	t.Logf("Response Body: %s", rec.Body.String()) // print isi JSON response
	t.Logf("Parsed: %+v", body)                    // print hasil decode map

	assert.NoError(t, err)
	assert.NotNil(t, body["data"])
	assert.Equal(t, "success", body["message"])

	// Verifikasi mock dipanggil
	mockService.AssertCalled(t, "GetAll", testifymock.Anything, "")
}
