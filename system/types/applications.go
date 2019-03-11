package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/rules"
)

type (
	Application struct {
		ID      uint64 `json:"applicationID,string" db:"id"`
		Name    string `json:"name" db:"name"`
		OwnerID uint64 `json:"ownerID" db:"rel_owner"`
		Enabled bool   `json:"enabled" db:"enabled"`

		Unify *ApplicationUnify `json:"unify,omitempty" db:"unify"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	ApplicationUnify struct {
		Name   string `json:"name,omitempty"`
		Listed bool   `json:"listed"`
		Icon   string `json:"icon"`
		Logo   string `json:"logo"`
		Url    string `json:"url"`
		Config string `json:"config"`
		Order  uint   `json:"order"`
	}
)

func (u *Application) Valid() bool {
	return u.ID > 0 && u.DeletedAt == nil
}

func (u *Application) Identity() uint64 {
	return u.ID
}

// Resource returns a system resource ID for this type
func (u *Application) Resource() rules.Resource {
	resource := rules.Resource{
		Service: "system",
		Scope:   "application",
	}
	if u != nil {
		resource.ID = u.ID
		resource.Name = u.Name
	}
	return resource
}

func (au *ApplicationUnify) Scan(value interface{}) error {
	switch value.(type) {
	case nil:
		au = nil
	case []uint8:
		if err := json.Unmarshal(value.([]byte), au); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into ApplicationUnify", value)
		}
	}

	return nil
}

func (au ApplicationUnify) Value() (driver.Value, error) {
	return json.Marshal(au)
}