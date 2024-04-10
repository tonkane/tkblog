package store

const defaultLimitValue = 20

// 设置默认的查询数
func defaultLimit(limit int) int {
	if limit == 0 {
		limit = defaultLimitValue
	}
	return limit
}