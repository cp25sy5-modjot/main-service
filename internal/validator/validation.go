package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Global singleton (หรือจะ inject ก็ได้)
var v *validator.Validate

func Validator() *validator.Validate {
	if v != nil {
		return v
	}
	v = validator.New()

	// ให้ใช้ชื่อ field จาก tag json แทนชื่อ struct field
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			return field.Name
		}
		// ตัด ,omitempty ออก
		if idx := strings.Index(name, ","); idx != -1 {
			return name[:idx]
		}
		return name
	})

	// ===== ตัวอย่าง custom tags =====
	// ตัวอย่าง: ต้องเป็น oneof (male female other)
	// _ = v.RegisterValidation("gender", func(fl validator.FieldLevel) bool {
	// 	val := strings.ToLower(fl.Field().String())
	// 	return val == "male" || val == "female" || val == "other"
	// })

	return v
}

type FieldError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

// MapValidationErrors แปลง validator.ValidationErrors → []FieldError (สำหรับ response)
func MapValidationErrors(err error) []FieldError {
	if err == nil {
		return nil
	}
	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	out := make([]FieldError, 0, len(verrs))
	for _, fe := range verrs {
		out = append(out, FieldError{
			Field: fe.Field(),   // จาก RegisterTagNameFunc จะเป็นชื่อ json
			Msg:   humanize(fe), // สร้างข้อความอ่านง่าย
		})
	}
	return out
}

func humanize(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", fe.Param())
	case "len":
		return fmt.Sprintf("length must be %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	// case "gender":
	// 	return "must be one of: male, female, other"
	default:
		// fallback เป็น "tag(param)" เช่น "gte(0)"
		if p := fe.Param(); p != "" {
			return fe.Tag() + "(" + p + ")"
		}
		return fe.Tag()
	}
}
