package shared

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func PrettyPrint[A any](a A) {
	bytes, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func DefaultGormHandle() (*gorm.DB, error) {
	dsn := "host=localhost user=user password=pass dbname=dev_db port=6876 sslmode=disable TimeZone=Europe/Warsaw"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

type ValidationError struct {
	Field      string `json:"field"`
	Message    string `json:"message"`
	ValueGiven any    `json:"value_given"`
}

func getJSONTag(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}

	return tag
}

// Counterintuitivly, validation also may return an error of its own.
// In that case 500 status code should be returned.
func RunValidation[A any](a A) ([]ValidationError, error) {
	validate := validator.New()
	validationErr := validate.Struct(a)
	if validationErrors, ok := validationErr.(validator.ValidationErrors); ok {
		structType := reflect.TypeOf(a)
		var errors []ValidationError
		for _, e := range validationErrors {
			field, _ := structType.FieldByName(e.Field())
			jsonTag := getJSONTag(field)
			errors = append(errors, ValidationError{
				Field:      jsonTag,
				Message:    e.Tag(),
				ValueGiven: e.Value(),
			})
		}
		return errors, nil
	}

	return nil, validationErr
}
