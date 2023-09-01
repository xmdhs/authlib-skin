package model

type APIStatus int

const (
	OK APIStatus = iota
	ErrInput
	ErrService
)
