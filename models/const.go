package models

var (
	LEVELS = []string{
		"最低",
		"非常低",
		"低",
		"中",
		"高",
		"非常高",
		"紧急",
	}
	levels []TextValue
)

func GetLevel() *[]TextValue {
	if len(levels) == 0 {
		for i := 0; i < len(LEVELS); i++ {
			levels = append(levels, TextValue{Value: i, Text: LEVELS[i]})
		}
	}

	return &levels
}
