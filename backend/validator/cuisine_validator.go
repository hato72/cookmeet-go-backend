package validator

import (
	"backend/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ICuisineValidator interface {
	CuisineValidate(cuisine model.Cuisine) error
}

type cuisineValidator struct{}

func NewCuisineValidator() ICuisineValidator { // cuisineValidatorのインスタンスを生成するためのコンストラクタ
	return &cuisineValidator{}
}

func (tv *cuisineValidator) CuisineValidate(cuisine model.Cuisine) error {
	return validation.ValidateStruct(&cuisine,
		validation.Field(
			&cuisine.Title,
			validation.Required.Error("title is required"), // 料理名に値が存在するか
			// validation.RuneLength(1, 10).Error("limited max 10 char"), //1文字から10文字までの文字数になっているかどうか
		),
	)
}
