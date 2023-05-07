package main

import (
  "encoding/json"
	"net/http"
	"net/http/httptest"
  "regexp"
	"testing"

  "Api/responses"

	"github.com/stretchr/testify/assert"
)

func TestUidRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/uid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
  var res responses.Response
  json.Unmarshal(w.Body.Bytes(), &res)
  assert.Regexp(t, regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"), res.Data["uid"])

}

