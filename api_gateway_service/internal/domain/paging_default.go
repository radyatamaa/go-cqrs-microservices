package domain

import (
	"strconv"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/response"
)

const (
	DEFAULT_PAGESIZE = 10
	MAX_PAGESIZE     = 100
	DEFAULT_PAGE     = 0
)

func PaginationQueryParamValidation(pageSizeStr, pageStr string) (int, int, error) {
	// default
	var pageSizeDefault = DEFAULT_PAGESIZE
	var pageDefault = DEFAULT_PAGE

	pageSizeInt, err := strconv.Atoi(pageSizeStr)
	if err != nil && pageSizeStr != "" {
		return 0, 0, response.ErrQueryParamInvalid
	}
	if err == nil {
		pageSizeDefault = pageSizeInt
		if pageSizeInt == 0 {
			pageSizeDefault = DEFAULT_PAGESIZE
		}
		// dafault value maximum pageSize = 100
		if pageSizeDefault > MAX_PAGESIZE {
			pageSizeDefault = MAX_PAGESIZE
		}
		if pageSizeInt < 0 {
			return 0, 0, response.ErrQueryParamInvalid
		}
	}

	pageInt, err := strconv.Atoi(pageStr)
	if err != nil && pageStr != "" {
		return 0, 0, response.ErrQueryParamInvalid
	}
	if err == nil {
		pageDefault = pageInt
		if pageInt == 0 {
			pageDefault = DEFAULT_PAGE
		}
		if pageInt < 0 {
			return 0, 0, response.ErrQueryParamInvalid
		}
	}

	return pageSizeDefault, pageDefault, nil
}
