// +build database

package models_test

import (
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func compareApp(t *testing.T, expected, actual models.App) {
	expected.CreatedAt = time.Time{}
	expected.UpdatedAt = time.Time{}
	expected.AppVersions = nil
	actual.CreatedAt = time.Time{}
	actual.UpdatedAt = time.Time{}
	actual.AppVersions = nil
	require.Equal(t, expected, actual)
}

func Test_AppService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appService := models.AppService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testApp := &models.App{
			AppSlug: "test-app-slug",
		}
		createdApp, err := appService.Create(testApp)
		require.NoError(t, err)
		require.False(t, createdApp.ID.String() == "")
		require.False(t, createdApp.CreatedAt == time.Time{})
		require.False(t, createdApp.UpdatedAt == time.Time{})
		require.Equal(t, createdApp.ID, createdApp.AppSettings.AppID)
	})

	t.Run("ok - when encrypted secret IV is filled, no secret will be generated", func(t *testing.T) {
		testApp := &models.App{
			AppSlug:           "test-app-slug",
			EncryptedSecretIV: []byte("somerandombytes"),
		}
		createdApp, err := appService.Create(testApp)
		require.NoError(t, err)
		require.False(t, createdApp.ID.String() == "")
		require.False(t, createdApp.CreatedAt == time.Time{})
		require.False(t, createdApp.UpdatedAt == time.Time{})
		require.Nil(t, createdApp.EncryptedSecret)
	})
}

func Test_AppService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appService := models.AppService{DB: dataservices.GetDB()}

	t.Run("ok - when searching based on app slug", func(t *testing.T) {
		testApp := createTestApp(t, &models.App{
			AppSlug: "test-app-slug",
		})

		foundApp, err := appService.Find(testApp)
		require.NoError(t, err)
		require.Equal(t, testApp, foundApp)
	})

	t.Run("ok - when searching based on app slug an api token", func(t *testing.T) {
		testApp := createTestApp(t, &models.App{
			AppSlug:  "test-app-slug-2",
			APIToken: "test-api-token",
		})

		foundApp, err := appService.Find(testApp)
		require.NoError(t, err)
		require.Equal(t, testApp, foundApp)
	})

	t.Run("error - when searching based on app slug an api token, but there's no such app", func(t *testing.T) {
		createTestApp(t, &models.App{
			AppSlug: "test-app-slug-3",
		})

		foundApp, err := appService.Find(&models.App{AppSlug: "test-app-slug-3", APIToken: "test-api-token"})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundApp)
	})
}

func Test_AppService_Update(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appService := models.AppService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testApps := []*models.App{
			createTestApp(t, &models.App{AppSlug: "test-app-1"}),
			createTestApp(t, &models.App{AppSlug: "test-app-2"}),
		}

		testApps[0].HeaderColor1 = "#FFFFFF"
		verrs, err := appService.Update(testApps[0], []string{"HeaderColor1"})
		require.Empty(t, verrs)
		require.NoError(t, err)

		t.Log("check if app got updated")
		foundApp, err := appService.Find(&models.App{Record: models.Record{ID: testApps[0].ID}})
		require.NoError(t, err)
		require.Equal(t, "#FFFFFF", foundApp.HeaderColor1)

		t.Log("check if no other app were updated")
		foundApp, err = appService.Find(&models.App{Record: models.Record{ID: testApps[1].ID}})
		require.NoError(t, err)
		compareApp(t, *testApps[1], *foundApp)
	})

	t.Run("when trying to update non-existing field", func(t *testing.T) {
		testApp := createTestApp(t, &models.App{AppSlug: "test-app-1"})
		verrs, err := appService.Update(testApp, []string{"NonExistingField"})
		require.EqualError(t, err, "Attribute name doesn't exist in the model")
		require.Equal(t, 0, len(verrs))
	})
}

func Test_AppService_Delete(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appService := models.AppService{DB: dataservices.GetDB()}

	testApp := createTestApp(t, &models.App{
		AppSlug:  "test-app-slug-2",
		APIToken: "test-api-token",
	})

	t.Run("when deleting an app", func(t *testing.T) {
		err := appService.Delete(&models.App{Record: models.Record{ID: testApp.ID}})
		require.NoError(t, err)
	})

	t.Run("error - when app is not found", func(t *testing.T) {
		err := appService.Delete(&models.App{Record: models.Record{ID: uuid.NewV4()}})

		require.Equal(t, err, gorm.ErrRecordNotFound)
	})
}
