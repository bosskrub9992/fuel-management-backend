package models

import (
	"errors"

	"github.com/jinleejun-corp/corelib/validators"
	"gopkg.in/guregu/null.v4"
)

type CreateCustomerRequest struct {
	ShortName string      `json:"shortName" validate:"required,max=150"`
	FirstName null.String `json:"firstName"`
	LastName  null.String `json:"lastName"`
}

type CreateCustomerResponse struct {
	ID int64 `json:"id"`
}

func (r CreateCustomerRequest) Vaildate() error {
	err := validators.Validate(r)
	if r.FirstName.Valid && len(r.FirstName.String) > 100 {
		err = errors.Join(err, errors.New("firstName is too long"))
	}
	if r.LastName.Valid && len(r.LastName.String) > 100 {
		err = errors.Join(err, errors.New("lastName is too long"))
	}
	return err
}
