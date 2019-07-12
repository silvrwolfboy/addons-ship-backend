package services

import (
	"net/http"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// AppContactPatchResponse ...
type AppContactPatchResponse struct {
	Data *models.AppContact `json:"data"`
}

// AppContactConfirmPatchHandler ...
func AppContactConfirmPatchHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppContactID, err := GetAuthorizedAppContactIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppContactService == nil {
		return errors.New("No App Contact Service provided")
	}
	appContact, err := env.AppContactService.Find(&models.AppContact{Record: models.Record{ID: authorizedAppContactID}})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}
	appContact.ConfirmedAt = time.Now()
	appContact.ConfirmationToken = nil
	err = env.AppContactService.Update(appContact, []string{"ConfirmedAt", "ConfirmationToken"})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	return httpresponse.RespondWithSuccess(w, AppContactPatchResponse{Data: appContact})
}