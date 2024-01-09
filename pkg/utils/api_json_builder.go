package utils

// Pagination 分页结构
type Pagination struct {
	PageTotal   int `json:"page_total"`   // 总页数
	TotalCount  int `json:"total_count"`  // 总数量
	PageCurrent int `json:"page_current"` // 当前页码
	PageSize    int `json:"page_size"`    // 页面大小
}

// JSONResult 用于返回ajax请求的基类
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

// BuildPagination 构建分页信息
func BuildPagination(page, limit, total int) *Pagination {
	// 没有数据就返回空分页
	if total == 0 {
		return &Pagination{
			PageCurrent: page,
		}
	}
	// 总页数
	pageTotal := total / limit
	lastSize := total % limit
	if lastSize != 0 {
		pageTotal++
	}
	pageSize := limit
	pageCurrent := page
	// 最后一页余数数量处理
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
