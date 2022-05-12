package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                    primitive.ObjectID `bson:"_id"`
	User_id               string             `json:"user_id"`
	Email                 *string            `json:"email" validate:"required,email"`
	Password              *string            `json:"password" validate:"required,min=8"`
	Created_at            time.Time          `json:"created_at"`
	Updated_at            time.Time          `json:"updated_at"`
	Token                 *string            `json:"token"`
	Refresh_token         *string            `json:"refresh_token"`
	Verified              bool               `json:"verified"`
	Generated_SignUp_code *string
}

type VerifyEmailPayload struct {
	Otp *string `json:"otp" validate:"required"`
}

// type User struct {
// 	ID           primitive.ObjectID `bson:"_id"`
// 	User_id      string             `json:"user_id"`
// 	Email        *string            `json:"email"`
// 	Password     *string            `json:"password"`
// 	First_name   *string            `json:"first_name"`
// 	Last_name    *string            `json:"last_name"`
// 	Company_name *string            `json:"company_name"`
// 	Phone_number *string            `json:"phone_number"`
// }
