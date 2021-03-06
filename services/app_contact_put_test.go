package services_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppContactPutHandler(t *testing.T) {
	httpMethod := "PATCH"
	url := "/apps/{app-slug}/contacts/{contact-id}"
	handler := services.AppContactPutHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppContactService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppContactService: &testAppContactService{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppContactID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppContactService: &testAppContactService{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{}, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						return nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppContactPutResponse{
				Data: &models.AppContact{
					NotificationPreferencesData: json.RawMessage(`{"new_version":false,"successful_publish":false,"failed_publish":false}`),
				},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{}, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						notificationPreferences, err := appContact.NotificationPreferences()
						require.NoError(t, err)
						require.Equal(t, models.NotificationPreferences{NewVersion: true}, notificationPreferences)
						return nil
					},
				},
			},
			requestBody:        `{"new_version":true}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppContactPutResponse{
				Data: &models.AppContact{
					NotificationPreferencesData: json.RawMessage(`{"new_version":true,"successful_publish":false,"failed_publish":false}`),
				},
			},
		})
	})

	t.Run("when request body is invalid JSON", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{}, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						return nil
					},
				},
			},
			requestBody:        `invald JSON`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Invalid request body, JSON decode failed"},
		})
	})

	t.Run("when error happens at finding app contact", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return nil, gorm.ErrRecordNotFound
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						appContact.ConfirmedAt = time.Time{}
						return nil
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: record not found",
		})
	})

	t.Run("when error happens at updating app contact", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{}, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						return errors.New("SOME-SQL-ERROR")
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}
