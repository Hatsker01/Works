package utils

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

func ParseFilter(s string) []string {
	return strings.Split(s, ",")
}

func StringSliceToInterfaceSlice(ss []string) []interface{} {
	is := make([]interface{}, len(ss))
	for i, v := range ss {
		is[i] = v
	}
	return is
}

type QueryParams struct {
	Filters  map[string]string
	Page     int64
	Limit    int64
	Ordering []string
	Search   string
}

func ParseQueryParams(queryParams map[string][]string) (*QueryParams, []string) {
	params := QueryParams{
		Filters:  make(map[string]string),
		Page:     1,
		Limit:    10,
		Ordering: []string{},
		Search:   "",
	}
	var errStr []string
	var err error

	for key, value := range queryParams {
		if key == "page" {
			params.Page, err = strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				errStr = append(errStr, "Invalid `page` param")
			}
			continue
		}

		if key == "limit" {
			params.Limit, err = strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				errStr = append(errStr, "Invalid `limit` param")
			}
			continue
		}

		if key == "search" {
			params.Search = value[0]
			continue
		}

		if key == "ordering" {
			params.Ordering = strings.Split(value[0], ",")
			continue
		}
		params.Filters[key] = value[0]
	}

	return &params, errStr
}

func StringToNullString(s string) (ns sql.NullString) {
	if s != "" {
		ns.Valid = true
		ns.String = s
		return ns
	}

	return ns
}

func StringToNullTime(s string) (nt sql.NullTime) {
	if s != "" {
		nt.Valid = true
		nt.Time, _ = time.Parse("2006-01-02", s)

		return nt
	}

	return nt
}

func IntToNullInt64(i int64) (ni sql.NullInt64) {
	if i != 0 {
		ni.Valid = true
		ni.Int64 = i

		return ni
	}

	return ni
}
