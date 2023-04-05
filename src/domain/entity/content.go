package entity

import (
	"strconv"
	"strings"
)

type Content struct {
	ID        int `json:"id"`
	ServiceID int `json:"service_id"`
	Service   *Service
	Name      string `json:"name"`
	Value     string `json:"value"`
	Tid       string `json:"tid"`
}

func (c *Content) GetName() string {
	return c.Name
}

func (c *Content) GetValue() string {
	return c.Value
}

func (c *Content) GetTid() string {
	return c.Tid
}

func (c *Content) SetPIN(pin int) {
	replacer := strings.NewReplacer("@pin", strconv.Itoa(pin))
	c.Value = replacer.Replace(c.Value)
}
