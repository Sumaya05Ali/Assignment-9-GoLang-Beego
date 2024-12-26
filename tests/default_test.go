package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	//"strings"

	_ "beeproject/routers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

// Test /cat route
func TestCatController_Get(t *testing.T) {
	req, err := http.NewRequest("GET", "/cat", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) 
	assert.Contains(t, w.Body.String(), "CatImages") 
}

