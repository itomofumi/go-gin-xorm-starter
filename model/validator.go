package model

import (
	"reflect"
	"sync"

	// "gopkg.in/go-playground/validator.v9"
	"github.com/go-playground/validator"
)

// StructValidator defines validations for models.
type StructValidator struct {
	once     sync.Once
	validate *validator.Validate
}

// ValidateStruct validates struct with tags.
func (v *StructValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

// Engine returns validate engine.
func (v *StructValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *StructValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// add any custom validations etc. here
		v.validate.RegisterStructValidation(FruitBodyStructLevelValidation, FruitBody{})
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
