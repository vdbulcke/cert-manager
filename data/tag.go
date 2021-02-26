package data

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag defines Tag for X509 Cert
// swagger:model
type Tag struct {
	// gorm.Model
	// the id for the Tag
	//
	// required: false
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`

	// the Tag name
	//
	// required: false
	Name string `json:"name"  gorm:"uniqueIndex" validate:"required,excludesall= "`

	// the Tag Description
	//
	// required: false
	Description string `json:"description" `

	// the List of Certificates for the tag
	//
	// required: false
	Certificates []Certificate `json:"certificates"  gorm:"many2many:tags_ref;"`

	// the CreatedAt timestamp for the tag
	//
	// required: false
	CreatedAt time.Time `json:"-" `

	// the UpdatedAt timestamp for the tag
	//
	// required: false
	UpdatedAt time.Time `json:"-" `
}

// BeforeCreate will set a UUID rather than numeric ID.
func (tag *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	if tag.ID == uuid.Nil {
		uuid := uuid.New()
		tag.ID = uuid
	}

	return
}

// NewTag create a new tag struct
func NewTag(name string) (*Tag, error) {
	return &Tag{
		Name: name,
	}, nil
}

// Tags a list of Tag
type Tags []*Tag
