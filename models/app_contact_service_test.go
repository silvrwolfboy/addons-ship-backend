// +build database

package models_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppContactService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	t.Run("ok", func(t *testing.T) {
		appContactService := models.AppContactService{DB: dataservices.GetDB()}
		testAppContact := &models.AppContact{Email: "an-email@addr.ess"}

		createdAppContact, verrs, err := appContactService.Create(testAppContact)
		require.NoError(t, err)
		require.Empty(t, verrs)
		require.False(t, createdAppContact.ID.String() == "")
		require.False(t, createdAppContact.CreatedAt.String() == "")
		require.False(t, createdAppContact.UpdatedAt.String() == "")
	})

	t.Run("when email is not valid", func(t *testing.T) {
		appContactService := models.AppContactService{DB: dataservices.GetDB()}
		testAppContact := &models.AppContact{Email: "not a valid email"}

		createdAppContact, verrs, err := appContactService.Create(testAppContact)
		require.NoError(t, err)
		require.Len(t, verrs, 1)
		require.EqualError(t, verrs[0], "email: Wrong format")
		require.Nil(t, createdAppContact)
	})

	t.Run("when email is too long", func(t *testing.T) {
		appContactService := models.AppContactService{DB: dataservices.GetDB()}
		testAppContact := &models.AppContact{Email: "123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-123456789-1234@bitrise.io"}

		createdAppContact, verrs, err := appContactService.Create(testAppContact)
		require.NoError(t, err)
		require.Len(t, verrs, 1)
		require.EqualError(t, verrs[0], "email: Too long")
		require.Nil(t, createdAppContact)
	})
}

func Test_AppContactService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appContactService := models.AppContactService{DB: dataservices.GetDB()}

	testApp := createTestApp(t, &models.App{AppSlug: "test-app-slug"})
	testAppContact := createTestAppContact(t, &models.AppContact{App: testApp, Email: "an-email@addr.ess"})

	t.Run("when querying an app contact that belongs to an app", func(t *testing.T) {
		foundAppContact, err := appContactService.Find(&models.AppContact{Record: models.Record{ID: testAppContact.ID}, AppID: testApp.ID})
		require.NoError(t, err)
		compareAppContacts(t, *testAppContact, *foundAppContact)
	})

	t.Run("error - when app contact is not found", func(t *testing.T) {
		otherTestApp := createTestApp(t, &models.App{AppSlug: "test-app-slug-2"})

		foundAppContact, err := appContactService.Find(&models.AppContact{Record: models.Record{ID: testAppContact.ID}, AppID: otherTestApp.ID})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundAppContact)
	})
}

func Test_AppContactService_FindAll(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appContactService := models.AppContactService{DB: dataservices.GetDB()}
	testApp := createTestApp(t, &models.App{AppSlug: "test-app-slug"})
	anotherTestApp := createTestApp(t, &models.App{AppSlug: "test-app-slug-2"})
	testAppContact1 := createTestAppContact(t, &models.AppContact{
		App:   testApp,
		Email: "someones@email.addr",
	})
	testAppContact2 := createTestAppContact(t, &models.AppContact{
		App:   testApp,
		Email: "someoneelses@email.addr",
	})
	createTestAppContact(t, &models.AppContact{
		App:   anotherTestApp,
		Email: "andanother@email.addr",
	})

	t.Run("when query all app contacts of test app", func(t *testing.T) {
		foundAppContacts, err := appContactService.FindAll(testApp)
		require.NoError(t, err)
		reflect.DeepEqual([]models.AppContact{*testAppContact2, *testAppContact1}, foundAppContacts)
	})
}

func Test_AppContactService_Update(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appContactService := models.AppContactService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testAppContacts := []*models.AppContact{
			createTestAppContact(t, &models.AppContact{Email: "an-email@addr.ess"}),
			createTestAppContact(t, &models.AppContact{Email: "other-email@addr.ess"}),
		}

		testAppContacts[0].NotificationPreferencesData = json.RawMessage(`{"new_version": true}`)

		err := appContactService.Update(testAppContacts[0], []string{"NotificationPreferencesData"})
		require.NoError(t, err)

		t.Log("Check if app contact got updated")
		foundAppContact, err := appContactService.Find(&models.AppContact{Record: models.Record{ID: testAppContacts[0].ID}})
		require.NoError(t, err)

		notificationPrefs, err := foundAppContact.NotificationPreferences()
		require.NoError(t, err)
		require.Equal(t, notificationPrefs.NewVersion, true)

		t.Log("check if no other app contact was updated")
		foundAppContact, err = appContactService.Find(&models.AppContact{Record: models.Record{ID: testAppContacts[1].ID}})
		require.NoError(t, err)
		compareAppContacts(t, *testAppContacts[1], *foundAppContact)
	})
}

func Test_AppContactService_Delete(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appContactService := models.AppContactService{DB: dataservices.GetDB()}
	appContact := createTestAppContact(t, &models.AppContact{Email: "an-email@addr.ess"})

	t.Run("ok", func(t *testing.T) {
		err := appContactService.Delete(&models.AppContact{Record: models.Record{ID: appContact.ID}})
		require.NoError(t, err)
	})

	t.Run("error - when app contact is not found", func(t *testing.T) {
		err := appContactService.Delete(&models.AppContact{Record: models.Record{ID: uuid.NewV4()}})
		require.Equal(t, err, gorm.ErrRecordNotFound)
	})
}
