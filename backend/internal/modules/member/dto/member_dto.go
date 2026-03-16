package dto

type MemberCreateRequest struct {
	Name         string `json:"name" validate:"required"`
	NIK          string `json:"nik" validate:"required,len=16"`
	Address      string `json:"address" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	Email        string `json:"email" validate:"omitempty,email"`
	CreateSystem bool   `json:"create_system_user"` // Whether to create login credentials
}

type MemberUpdateRequest struct {
	Name    string `json:"name" validate:"omitempty"`
	Address string `json:"address" validate:"omitempty"`
	Phone   string `json:"phone" validate:"omitempty"`
	Status  string `json:"status" validate:"omitempty,oneof=pending active inactive"`
}

type MemberResponse struct {
	ID           uint                     `json:"id"`
	UserID       uint                     `json:"user_id,omitempty"`
	MemberNumber string                   `json:"member_number"`
	Name         string                   `json:"name"`
	NIK          string                   `json:"nik"`
	Address      string                   `json:"address"`
	Phone        string                   `json:"phone"`
	Status       string                   `json:"status"`
	CreatedAt    string                   `json:"created_at"`
	Documents    []MemberDocumentResponse `json:"documents,omitempty"`
	Cards        []MemberCardResponse     `json:"cards,omitempty"`
}

type MemberDocumentResponse struct {
	ID      uint   `json:"id"`
	Type    string `json:"type"`
	FileURL string `json:"file_url"`
}

type MemberCardResponse struct {
	ID         uint   `json:"id"`
	CardNumber string `json:"card_number"`
	Status     string `json:"status"`
}

type DocumentUploadRequest struct {
	Type string `form:"type" validate:"required,oneof=ktp kk selfie signature"`
}
