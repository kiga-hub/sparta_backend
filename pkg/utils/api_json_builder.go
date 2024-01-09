package utils

// Pagination pagination
type Pagination struct {
	PageTotal   int `json:"page_total"`   // total page
	TotalCount  int `json:"total_count"`  // total count
	PageCurrent int `json:"page_current"` // current page
	PageSize    int `json:"page_size"`    // page size
}

// JSONResult -
type JSONResult struct {
	Code   Code        `json:"code"`
	Msg    Msg         `json:"msg"`
	Detail string      `json:"detail,omitempty"`
	Data   interface{} `json:"data"`
	Page   *Pagination `json:"page,omitempty"`
}

// SuccessJSONData -
func SuccessJSONData(data interface{}) *JSONResult {
	return &JSONResult{
		Code: Success,
		Msg:  SuccessMsg,
		Data: data,
		Page: nil,
	}
}

// SuccessJSONPageData -
func SuccessJSONPageData(data interface{}, page *Pagination) *JSONResult {
	return &JSONResult{
		Code: Success,
		Msg:  SuccessMsg,
		Data: data,
		Page: page,
	}
}

// FailJSONData -
func FailJSONData(code Code, msg Msg, err error) *JSONResult {
	var detail string
	if err != nil {
		detail = err.Error()
	}
	return &JSONResult{
		Code:   code,
		Msg:    msg,
		Detail: detail,
		Data:   nil,
		Page:   nil,
	}
}

// BuildPagination build pagination
func BuildPagination(page, limit, total int) *Pagination {
	// have no data
	if total == 0 {
		return &Pagination{
			PageCurrent: page,
		}
	}
	// page total
	pageTotal := total / limit
	lastSize := total % limit
	if lastSize != 0 {
		pageTotal++
	}
	pageSize := limit
	pageCurrent := page
	// set page size
	if page == pageTotal && lastSize != 0 {
		pageSize = lastSize
	}
	return &Pagination{
		PageTotal:   pageTotal,
		TotalCount:  total,
		PageCurrent: pageCurrent,
		PageSize:    pageSize,
	}
}
