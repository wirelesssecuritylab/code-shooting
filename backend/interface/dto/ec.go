package dto

import "code-shooting/domain/entity/ec"

type QueryECRsp struct {
	ECs []EC `json:"ecs"`
}

type EC struct {
	ec.Organization

	Id                string `json:"id"`
	Importer          string `json:"importer"`
	ImportTime        string `json:"importTime"`
	ConvertedToTarget bool   `json:"convertedToTarget"`
}
