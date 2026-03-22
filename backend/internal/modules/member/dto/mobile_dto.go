package dto

type MemberDashboardResponse struct {
	Profile       MemberResponse `json:"profile"`
	TotalSavings  float64        `json:"total_savings"`
	TotalLoans    float64        `json:"total_loans"`
	QRCode        string         `json:"qr_code"`
}

type MobileLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
