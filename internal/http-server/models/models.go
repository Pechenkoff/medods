package models

type AccessRequest struct {
	GUID  string `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email string `json:"email" example:"user@example.com"`
}

type Response struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzZDU2NyIsImV4cCI6MTYzMjU3MjAwMCwiaWF0IjoxNjMyNTY4NDAwfQ.d6l5N1dlv8HHSX3aWcyhnSY7OLqZ_Px4dUlYVt5yXy_o1s0hD-Gy2GqqP6lVRcCqKbl-dfxRmXhyyQHqt8RThQ"`
	RefreshToken string `json:"refresh_token" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
}

type RefreshRequest struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzZDU2NyIsImV4cCI6MTYzMjU3MjAwMCwiaWF0IjoxNjMyNTY4NDAwfQ.d6l5N1dlv8HHSX3aWcyhnSY7OLqZ_Px4dUlYVt5yXy_o1s0hD-Gy2GqqP6lVRcCqKbl-dfxRmXhyyQHqt8RThQ"`
	RefreshToken string `json:"refresh_token" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
}

type BadRequestResponse struct {
	Error string `json:"error" example:"bad request"`
}

type ServerErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}
