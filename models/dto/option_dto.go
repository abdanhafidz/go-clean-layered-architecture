package dto

import entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"

type OptionsRequest struct {
	OptionName  string   `json:"option_name" binding:"required"`
	OptionValue []string `json:"option_values" binding:"required"`
}

type OptionsResponse struct {
	Options []entity.Options `json:"options"`
}
