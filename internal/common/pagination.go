package common

type Pagination struct {
	Page     uint32 `validate:"required,min=1,max=10000000"`
	PageSize uint32 `validate:"required,min=1,max=50"`
}
