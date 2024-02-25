package requests

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
