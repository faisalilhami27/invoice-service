package template

import "errors"

var (
	ErrTemplateAlreadyExist = errors.New(
		"this template with category and service already exist, please change the category or service")
)
