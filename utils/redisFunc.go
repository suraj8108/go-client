package utils

import (
	"encoding/json"
	"fmt"

	"github.com/suraj8108/clientApp/model"
)

func CreateRedisKeys(startKey int, keyRange int) []string {
	keys := make([]string, keyRange)
	for i := 1; i <= keyRange; i++ {
		keyId := startKey + i
		customKey := fmt.Sprintf("record%v", keyId)
		keys = append(keys, customKey)
	}
	return keys
}

func UnMarshalRedisData(data string) model.Student {
	byteData := []byte(data)

	var studentDetails model.Student
	json.Unmarshal(byteData, &studentDetails)
	return studentDetails
}

func UnMarshalRedisBulkData(data []any) []model.Student {
	allStudentDetails := make([]model.Student, 0)

	for _, details := range data {
		studentData, ok := details.(string)
		var studentDetails model.Student
		if !ok {
			continue
		} else {
			studentDetails = UnMarshalRedisData(studentData)
		}
		// fmt.Println(studentDetails)
		allStudentDetails = append(allStudentDetails, studentDetails)

	}
	fmt.Println(allStudentDetails)
	return allStudentDetails
}
