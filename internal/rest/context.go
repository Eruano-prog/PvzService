package rest

type RequestContextKey string

var (
	IdKey = RequestContextKey("userID")
)
