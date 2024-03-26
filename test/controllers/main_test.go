package controllers

import (
	"github.com/NeptuneYeh/simplebank/init/auth"
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/NeptuneYeh/simplebank/init/logger"
	"github.com/gin-gonic/gin"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	configModule := config.NewModule("../../")
	configModule.AccessTokenDuration = time.Minute
	_ = logger.NewModule()
	_ = auth.NewModule()

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

//func newMockStore(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockStore := mockdb.NewMockStore(ctrl)
//	_ = store.NewModuleForTest(mockStore)
//	_ = myGin.NewModule()
//}
