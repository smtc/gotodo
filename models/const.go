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

	TASK_STATUS_CREATED  = "created"
	TASK_STATUS_PROGRESS = "progress"
	TASK_STATUS_TESTING  = "testing"
	TASK_STATUS_FINISHED = "finished"
	TASK_STATUS_CANCELED = "canceled"

	TASK_STATUS = map[string]string{
		TASK_STATUS_CREATED:  "已创建",
		TASK_STATUS_PROGRESS: "进行中",
		TASK_STATUS_TESTING:  "发布测试",
		TASK_STATUS_FINISHED: "已完成",
		TASK_STATUS_CANCELED: "已撤销",
	}
)

func GetLevel() []TextValue {
	if len(levels) == 0 {
		for i := 0; i < len(LEVELS); i++ {
			levels = append(levels, TextValue{Value: i, Text: LEVELS[i]})
		}
	}

	return levels
}
