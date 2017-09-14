package open_exchange_rates

import (
	"fmt"
)

type RateStore struct {
	Rates map[string]float64 `json:"rates"`
	Timestamp int32 `json:"timestamp"`
	Base string `json:"base"`
}


type Rate struct {
	Base string
	Target string
	Rate float64
}

type ErrorResponse struct {
	ErrorCode   int64  `json:"status"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v", r.Description)
}