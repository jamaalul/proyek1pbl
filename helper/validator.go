package helper

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	validate           *validator.Validate
	validateOnce       sync.Once
	nimRegex           = regexp.MustCompile(`^[0-9]{12}$`)
	tahunAkademikRegex = regexp.MustCompile(`^[0-9]{4}/[0-9]{4}-(Ganjil|Genap)$`)
)

func getValidator() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New()

		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		_ = validate.RegisterValidation("nim", func(fl validator.FieldLevel) bool {
			return nimRegex.MatchString(fl.Field().String())
		})

		_ = validate.RegisterValidation("angkatan", func(fl validator.FieldLevel) bool {
			year, err := strconv.Atoi(fl.Field().String())
			if err != nil {
				year = int(fl.Field().Int())
			}
			currentYear := time.Now().Year()
			return year >= 1000 && year <= currentYear
		})

		_ = validate.RegisterValidation("tahunakademik", func(fl validator.FieldLevel) bool {
			return tahunAkademikRegex.MatchString(fl.Field().String())
		})
	})
	return validate
}

func ValidateStruct(s any) map[string][]string {
	v := getValidator()
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	errsMap := make(map[string][]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := fieldErr.Field()
			msg := translateValidation(fieldErr)
			errsMap[field] = append(errsMap[field], msg)
		}
	}
	return errsMap
}

func translateValidation(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", field)
	case "email":
		return fmt.Sprintf("%s harus berupa format email yang valid", field)
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", field, fe.Param())
	case "gt":
		return fmt.Sprintf("%s harus lebih besar dari %s", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s minimal bernilai %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("%s harus lebih kecil dari %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s maksimal bernilai %s", field, fe.Param())
	case "alphanum":
		return fmt.Sprintf("%s hanya boleh berisi huruf dan angka", field)
	case "nim":
		return fmt.Sprintf("%s harus berupa 12 digit nomor induk mahasiswa", field)
	case "angkatan":
		return fmt.Sprintf("%s harus berupa tahun valid tidak melebihi tahun berjalan", field)
	case "tahunakademik":
		return fmt.Sprintf("%s format harus 'YYYY/YYYY-Ganjil' atau 'YYYY/YYYY-Genap'", field)
	default:
		return fmt.Sprintf("%s tidak valid (%s)", field, fe.Tag())
	}
}
