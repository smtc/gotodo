package models

type TextValue struct {
	Text  string      `json:"text"`
	Value interface{} `json:"value"`
}
