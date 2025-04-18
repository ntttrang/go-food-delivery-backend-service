package usermodel

import "errors"

type IntrospectReq struct {
	Token string `json:"token"`
}

func (c *IntrospectReq) Validate() error {
	if c.Token == "" {
		return errors.New("token is required")
	}

	return nil
}
