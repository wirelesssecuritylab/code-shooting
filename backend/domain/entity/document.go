package entity

type DocumentsInfo struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	ParentId  string          `json:"parentId"`
	Type      string          `json:"type"`
	Ducoments []string        `json:"ducoments"`
	Children  []DocumentsInfo `json:"children"`
}

type DocumentsDetailReq struct {
	FilePath  	string  `json:"filePath"`
}
