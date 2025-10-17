package models

import "time"

type User struct {
	ID              int       `json:"id"`
	MembershipID    string    `json:"membership_id"`    // รหัสสมาชิก เช่น LBK001234
	FirstName       string    `json:"first_name"`       // ชื่อ
	LastName        string    `json:"last_name"`        // นามสกุล
	PhoneNumber     string    `json:"phone_number"`     // เบอร์โทรศัพท์
	Email           string    `json:"email"`            // อีเมล
	MembershipLevel string    `json:"membership_level"` // ระดับสมาชิก (Gold, Silver, Bronze)
	Points          int       `json:"points"`           // แต้มคงเหลือ
	JoinedDate      time.Time `json:"joined_date"`      // วันที่สมัครสมาชิก
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	MembershipLevel string `json:"membership_level"` // Default: Bronze
}

type UpdateUserRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	PhoneNumber     string `json:"phone_number"`
	Email           string `json:"email"`
	MembershipLevel string `json:"membership_level"`
	Points          *int   `json:"points"` // pointer เพื่อให้แยกระหว่าง 0 กับ null
}
