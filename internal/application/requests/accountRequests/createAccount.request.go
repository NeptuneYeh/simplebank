package accountRequests

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}
