package models

var (
	LEVELS = map[int]string{
		0: "最低",
		1: "非常低",
		2: "低",
		3: "中",
		4: "高",
		5: "非常高",
		6: "紧急",
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
