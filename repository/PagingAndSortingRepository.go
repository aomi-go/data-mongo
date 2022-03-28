package repository

import (
	datacommon "github.com/aomi-go/data-common"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PagingAndSortingRepository interface {
	FindAll(filter interface{}, pageable datacommon.Pageable, opts ...*options.FindOptions) datacommon.Page
}
