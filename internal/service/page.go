package service

type PaginateReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginateRsp struct {
	PaginateReq

	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
