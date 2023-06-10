package score

import "code-shooting/domain/entity"

type Calculator interface {
	CalculateShootingResult(stdAnswer string, tp *TargetAnswerPaper) (*entity.TargetResult, error)
}

type DefectCoder interface {
	EncodeDefect(defectClass, defectSubclass, defectDescribe string) string
	DecodeDefect(defectId string) (string, string, string)
}
