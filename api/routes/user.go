package routes

import (
	"net/http"

	"crypto/sha256"
	"encoding/hex"
	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/api/user"
)

const (
	UpdateUserURL = "/user"
)

func UpdateUser(c apicontext.Context) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json: "newPassword"`

	}
	if err := c.Bind(&req); err != nil {
		return err
	}

	tenant := ""
	if v := c.Tenant(); v != nil {
		tenant = v.ID
	}
	if req.OldPassword != "" {
		sum := sha256.Sum256([]byte(req.OldPassword))
		sumByte:= sum[:]
		req.OldPassword = hex.EncodeToString(sumByte)
	}
	if req.NewPassword!= "" {
		sum := sha256.Sum256([]byte(req.NewPassword))
		sumByte:= sum[:]
		req.NewPassword = hex.EncodeToString(sumByte)
	}

	svc := user.NewService(c.Store())

	if err := svc.UpdateDataUser(c.Ctx(), req.Username, req.Email, req.OldPassword, req.NewPassword, tenant); err != nil {
		if err == user.ErrUnauthorized {
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	return c.JSON(http.StatusOK, nil)
}
