package score

import "io"

type Range struct {
	RangeId  string
	Language string
}

type AnswerPaper struct {
	Range
	Data io.Reader
}

type TargetAnswer struct {
	FileName     string
	StartLineNum int
	EndLineNum   int
	StartColNum  int
	EndColNum    int
	DefectCode   string
	Remark       string
}

type TargetAnswerPaper struct {
	Range   Range
	Answers []TargetAnswer
}
