package assembler

import (
	"code-shooting/interface/dto"
	"code-shooting/infra/shooting-result"
)

func TargetAnswer2Dto(ts []shootingresult.TargetAnswer) []*dto.TargetAnswer {
	results := make([]*dto.TargetAnswer, 0, len(ts))
	for _, t := range ts {
		results = append(results, &dto.TargetAnswer{
			FileName:     t.FileName,
			StartLineNum: t.StartLineNum,
			EndLineNum:   t.EndLineNum,
			DefectClass:  t.DefectClass,
			DefectSubClass: t.DefectSubClass,
			DefectDescribe: t.DefectDescribe,
		})
	}
	return results
}
