package riman

import (
	"github.com/pistolricks/kbeauty-api/internal/data"
)

func calculateMetadata(totalRecords, page, pageSize int) data.Metadata {
	if totalRecords == 0 {

		return data.Metadata{}
	}

	return data.Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
