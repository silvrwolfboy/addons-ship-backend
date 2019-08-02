package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

func Test_LoginPostHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/login"
	handler := services.LoginPostHandler

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
	})

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppService: &testAppService{},
		},
	})

	t.Run("ok", func(t *testing.T) {
		testAppID := uuid.FromStringOrNil("")
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AddonFrontendHostURL: "http://ship.bitrise.io",
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, testAppID, app.ID)
						app.AppSlug = "test-app-slug"
						app.APIToken = "test-app-api-token"
						return app, nil
					},
				},
			},
			expectedStatusCode:       http.StatusMovedPermanently,
			expectedResponseLocation: "http://ship.bitrise.io/apps/test-app-slug?token=test-app-api-token",
		})
	})

	t.Run("when app not found", func(t *testing.T) {
		testAppID := uuid.FromStringOrNil("")
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AddonFrontendHostURL: "http://ship.bitrise.io",
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return app, gorm.ErrRecordNotFound
					},
				},
			},
			expectedInternalErr: "SQL Error: record not found",
		})
	})
}