package utils

import (
	"strconv"
)

type PageQuery struct {
	Page  int
	Limit int
}

func ParsePageQuery(pageStr, limitStr string) PageQuery {
	page := 1
	limit := 20
	if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 100 {
		limit = v
	}
	return PageQuery{Page: page, Limit: limit}
}

func (pq PageQuery) Offset() int { return (pq.Page - 1) * pq.Limit }
