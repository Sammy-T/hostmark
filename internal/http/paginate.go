package http

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := r.URL.Query()

		page, err := strconv.Atoi(query.Get("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(query.Get("page_size"))
		switch {
		case err != nil:
			pageSize = 5
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 5
		}

		offset := (page - 1) * pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}
