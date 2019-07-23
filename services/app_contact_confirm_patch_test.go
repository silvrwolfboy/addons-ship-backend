package services_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppContactConfirmPatchHandler(t *testing.T) {
	httpMethod := "PATCH"
	url := "/confirm_email"
	handler := services.AppContactConfirmPatchHandler

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
						return &models.AppContact{App: &models.App{}}, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						appContact.ConfirmedAt = time.Time{}
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppContactPatchResponse{
				Data: services.AppContactPatchResponseData{
					AppContact: &models.AppContact{},
					App:        &models.App{},
				},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		testApp := models.App{Record: models.Record{ID: uuid.NewV4()}}
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.FromStringOrNil("8a230385-0113-4cf3-a9c6-469a313e587a"),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						require.Equal(t, uuid.FromStringOrNil("8a230385-0113-4cf3-a9c6-469a313e587a"), appContact.ID)
						appContact.App = &testApp
						return appContact, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						require.Equal(t, uuid.FromStringOrNil("8a230385-0113-4cf3-a9c6-469a313e587a"), appContact.ID)
						require.Nil(t, appContact.ConfirmationToken)
						require.Equal(t, []string{"ConfirmedAt", "ConfirmationToken"}, whitelist)
						require.NotEqual(t, time.Time{}, appContact.ConfirmedAt)
						appContact.ConfirmedAt = time.Time{}
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppContactPatchResponse{
				Data: services.AppContactPatchResponseData{
					AppContact: &models.AppContact{Record: models.Record{ID: uuid.FromStringOrNil("8a230385-0113-4cf3-a9c6-469a313e587a")}},
					App:        &testApp,
				},
			},
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
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}
