package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
)

func DeCodeBody[T any](r io.Reader, v *validator.Validate) (T, error) {
	jr := json.NewDecoder(r)
	var a T
	err := jr.Decode(&a)
	if err != nil {
		return a, fmt.Errorf("DeCodeBody: %w", err)
	}

	err = v.Struct(a)
	if err != nil {
		return a, fmt.Errorf("DeCodeBody: %w", err)
	}
	return a, nil
}
