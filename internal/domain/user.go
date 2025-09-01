package domain

import "time"

// User represents a user in the system.
// This is a domain entity that maps to the onb-ptf-users collection.
type User struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Email       string    `bson:"email" json:"email" validate:"required,email"`
	Password    string    `bson:"pass" json:"-"`
	DateCreated time.Time `bson:"date_created"`
	DateUpdated time.Time `bson:"updated_created"`
}
