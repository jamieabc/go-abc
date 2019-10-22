package report

import "encoding/json"

type ABC struct {
	Assignments int `json:"Assignments"`
	Branches    int `json:"Branches"`
	Conditions  int `json:"Conditions"`
}

type Report struct {
	Path  string `json:"Path"`
	Score int    `json:"Score"`
	ABC   ABC    `json:"Abc"`
}

func (r Report) String() string {
	b, _ := json.MarshalIndent(r, "", "\t")
	return string(b)
}
