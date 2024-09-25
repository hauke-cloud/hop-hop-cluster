package domain

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// CopyValues copies the values of the other instance config to this instance config
// Parameters:
// - to: The destination object
// - from: The source object
// Returns:
// - error: The error if any
func CopyStructs(from, to interface{}) error {
	// Ensure both dst and src are pointers
	if reflect.ValueOf(from).Kind() != reflect.Ptr || reflect.ValueOf(to).Kind() != reflect.Ptr {
		return fmt.Errorf("both arguments must be pointers")
	}
	fromValue := reflect.ValueOf(to).Elem()
	toValue := reflect.ValueOf(from).Elem()

	// Ensure both dst and src are pointers
	if fromValue.Kind() != reflect.Ptr || toValue.Kind() != reflect.Ptr {
		return fmt.Errorf("both arguments must be pointers")
	}

	// Ensure both dst and src are valid and are struct types
	if !fromValue.IsValid() || !toValue.IsValid() {
		return fmt.Errorf("invalid source or destination")
	}

	if fromValue.Kind() != reflect.Struct || toValue.Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be structs")
	}

	for i := 0; i < toValue.NumField(); i++ {
		fieldType := toValue.Type().Field(i)

		// Skip if field is part of gorm.Model
		if fieldType.Anonymous && fieldType.Type == reflect.TypeOf(gorm.Model{}) {
			continue
		}

		// Skip if copy is disabled by tag
		tag := fieldType.Tag.Get("compy")
		if tag != "" && (tag == "compare" || tag == "none") {
			continue
		}

		fromField := fromValue.Field(i)
		toField := toValue.Field(i)

		if toField.CanSet() {
			toField.Set(fromField)
		}
	}

	return nil
}

// CompareStructs compares the values of two structs
// Parameters:
// - one: The first struct
// - two: The second struct
// Returns:
// - bool: The result of the comparison
// - error: The error if any
func CompareStructs(one, two interface{}) (bool, error) {
	// Ensure both dst and src are pointers
	if reflect.ValueOf(one).Kind() != reflect.Ptr || reflect.ValueOf(two).Kind() != reflect.Ptr {
		return false, fmt.Errorf("both arguments must be pointers")
	}

	oneValue := reflect.ValueOf(two).Elem()
	twoValue := reflect.ValueOf(one).Elem()

	// Ensure both dst and src are valid and are struct types
	if !oneValue.IsValid() || !twoValue.IsValid() {
		return false, fmt.Errorf("invalid source or destination")
	}

	if oneValue.Kind() != reflect.Struct || twoValue.Kind() != reflect.Struct {
		return false, fmt.Errorf("both arguments must be structs")
	}

	for i := 0; i < twoValue.NumField(); i++ {
		fieldType := twoValue.Type().Field(i)

		// Skip if field is part of gorm.Model
		if fieldType.Anonymous && fieldType.Type == reflect.TypeOf(gorm.Model{}) {
			continue
		}

		// Skip if copy is disabled by tag
		tag := fieldType.Tag.Get("compy")
		if tag != "" && (tag == "copy" || tag == "none") {
			continue
		}

		interfaceOneField := oneValue.Field(i).Interface()
		interfaceTwoField := twoValue.Field(i).Interface()

		if !reflect.DeepEqual(interfaceOneField, interfaceTwoField) {
			return false, nil
		}
	}

	return true, nil
}
