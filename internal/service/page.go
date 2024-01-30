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

// NewPaginateRsp get response pageinate data
func NewPaginateRsp(list interface{}, total int64, req ...PaginateReq) *PaginateRsp {
	rsp := &PaginateRsp{
		List:  list,
		Total: total,
	}

	if len(req) > 0 {
		rsp.Limit = req[0].Limit
		rsp.Page = req[0].Page
	}

	return rsp
}
