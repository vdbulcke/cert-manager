package data

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Certificate defines a X509 Cert
// swagger:model
type Certificate struct {
	gorm.Model
	// the id for the certifcate
	//
	// required: false
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	// the sha256 of the pem cert
	//
	// required: false
	SHA256 string `json:"sha256"  gorm:"uniqueIndex" validate:"required"`
	// the Subject of the pem cert
	//
	// required: false
	Subject string `json:"subject,omitempty"`
	// the Issuer of the pem cert
	//
	// required: false
	Issuer string `json:"issuer,omitempty"`
	// the Serial Number of the pem cert
	//
	// required: false
	SerialNumber string `json:"serial_number,omitempty"`
	// the List of Subject Alternative Names of the pem cert
	//
	// required: false
	SANs string `json:"sans,omitempty"`
	// the Not Before validity of the pem cert
	//
	// required: false
	NotBefore time.Time `json:"not_before"`
	// the Not After validity of the pem cert
	//
	// required: false
	NotAfter time.Time `json:"not_after"`
	// the signature Algorithm of the pem cert
	//
	// required: false
	SignatureAlgorithm string `json:"sigalg"`
	// the Authority Key ID of the pem cert
	//
	// required: false
	AKI string `json:"authority_key_id"`
	// the Subject Key ID of the pem cert
	//
	// required: false
	SKI string `json:"subject_key_id"`
	// the OCSP server of the pem cert
	//
	// required: false
	OCSP string `json:"ocsp"`
	// the CRL server of the pem cert
	//
	// required: false
	CRL string `json:"crl"`
	// the Issuing CA URL of the pem cert
	//
	// required: false
	IssuingCAUrl string `json:"issuing_ca_url"`
	// if the pem cert  Is CA
	//
	// required: false
	IsCA bool `json:"is_ca"`
	// the Raw PEM string of the pem cert
	//
	// required: false
	RawPEM string `json:"pem"`
	// the List of tags for the pem cert
	//
	// required: false
	Tags []Tag `json:"tags"  gorm:"many2many:tags_ref;"`
	// Tags string `json:"tags" `
	// the CreatedAt timestamp for the Cert
	//
	// required: false
	CreatedAt time.Time `json:"-" `

	// the UpdatedAt timestamp for the Cert
	//
	// required: false
	UpdatedAt time.Time `json:"-" `
}

// BeforeCreate will set a UUID rather than numeric ID.
func (cert *Certificate) BeforeCreate(tx *gorm.DB) (err error) {
	if cert.ID == uuid.Nil {
		uuid := uuid.New()
		cert.ID = uuid
	}

	return
}

// NewCertificate Creates a Certificate from PEM string
func NewCertificate(pem string) (*Certificate, error) {

	// parse PEM string
	x509Cert, err := ParseX509FromPEM(pem)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		SHA256:             GetSHA256FingerprintFromX509Cert(x509Cert),
		Subject:            GetSubjectFromX509Cert(x509Cert),
		Issuer:             GetIssuerFromX509Cert(x509Cert),
		SerialNumber:       GetSerialNumberFromX509Cert(x509Cert),
		SignatureAlgorithm: GetSignatureAlgorithmFromX509Cert(x509Cert),
		AKI:                GetAKIFromX509Cert(x509Cert),
		SKI:                GetSKIFromX509Cert(x509Cert),
		SANs:               GetSNASFromX509Cert(x509Cert),
		OCSP:               GetOCSPFromX509Cert(x509Cert),
		CRL:                GetCRLFromX509Cert(x509Cert),
		IssuingCAUrl:       GetIssuingCAFromX509Cert(x509Cert),
		IsCA:               GetIsCAFromX509Cert(x509Cert),
		RawPEM:             pem,
	}, nil

}

// Tag defines Tag for X509 Cert
// swagger:model
type Tag struct {
	gorm.Model
	// the id for the Tag
	//
	// required: false
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`

	// the Tag name
	//
	// required: false
	Name string `json:"name"  gorm:"uniqueIndex" validate:"required,excludesall= "`

	// the List of Certificates for the tag
	//
	// required: false
	Certificates []Certificate `json:"_"  gorm:"many2many:tags_ref;"`

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
