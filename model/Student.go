package model

type Student struct {
	StudentId       string `json:"StudentId"`
	StudentRedisKey string `json:"StudentRedisKey"`
	StudentName     string `json:"StudentName"`
	StudentEmail    string `json:"StudentEmail"`
}
