package httpctl

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegister(t *testing.T) {
	w := httptest.NewRecorder()
	gctx, _ := gin.CreateTestContext(w)

}
