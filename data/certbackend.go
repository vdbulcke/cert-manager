package data

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/hashicorp/go-hclog"
)

// CertBackend Struct for interaction with Cert DB
type CertBackend struct {
	logger hclog.Logger
	db     *gorm.DB
	v      *Validation
}

// NewCertBackend creates a new CertBackend
func NewCertBackend(logger hclog.Logger, db *gorm.DB, v *Validation) *CertBackend {

	return &CertBackend{
		logger: logger,
		db:     db,
		v:      v,
	}
}

// Tags a list of Tag
type Tags []*Tag

// Certificates a list of Certificate
type Certificates []*Certificate

//
// Cert functions
//

// CreateCertificate create a new Certificate from pem and insert it into DB
func (certBackend *CertBackend) CreateCertificate(pem string) (*Certificate, error) {
	certBackend.logger.Debug("Creating Cert...")

	cert, err := NewCertificate(pem)
	if err != nil {
		return nil, err
	}

	validationErr := certBackend.v.Validate(cert)
	if validationErr != nil {
		return nil, errors.New(strings.Join(validationErr.Errors(), "\n"))
	}

	// Create new entry
	result := certBackend.db.Create(&cert)
	if result.Error != nil {
		// log.Fatal(result.Error)
		return nil, result.Error
	}

	return cert, nil
}

// GetCertByID return Tag (without associated cert)
func (certBackend *CertBackend) GetCertByID(uuid uuid.UUID) (*Certificate, error) {
	certBackend.logger.Debug("Getting Certificate...")

	var cert Certificate
	// preload the tags
	result := certBackend.db.Preload(clause.Associations).Find(&cert, uuid)
	if result.Error != nil {
		return nil, errors.New("Cert not found " + uuid.String())
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return nil, errors.New("Cert not found " + uuid.String())
	}

	// certBackend.logger.Info("Found cert", cert)

	return &cert, nil
}

// GetCertByFingerprint return Tag (without associated cert)
func (certBackend *CertBackend) GetCertByFingerprint(sha256 string) (*Certificate, error) {
	certBackend.logger.Debug("Getting Certificate...")

	var cert Certificate

	result := certBackend.db.Preload("Tags").Where("sha256 = ?", sha256).First(&cert)
	if result.Error != nil {
		return nil, errors.New("Cert not found " + sha256)
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return nil, errors.New("Cert not found " + sha256)
	}

	return &cert, nil
}

// SetCertTagNameByID return Tag (without associated cert)
func (certBackend *CertBackend) SetCertTagNameByID(uuid uuid.UUID, tagName string) (*Certificate, error) {
	certBackend.logger.Debug("Setting tag for Certificate...")

	// getting tag
	tag, err := certBackend.GetTagByName(tagName)
	if err != nil {
		return nil, err
	}

	certBackend.logger.Info("getTabGyName", tag)

	// get cert by uuid
	cert, cerr := certBackend.GetCertByID(uuid)
	if cerr != nil {
		return nil, cerr
	}

	// check if tag already in list
	for _, t := range cert.Tags {
		if t.ID == tag.ID {
			return cert, nil
		}
	}

	// add tag in tags list
	newTags := append(cert.Tags, *tag)

	tagErr := certBackend.db.Model(&cert).Association("Tags").Append(newTags)
	if tagErr != nil {
		return nil, errors.New("error updating cert " + uuid.String() + " with tag " + tagName)
	}

	// updateresult := certBackend.db.Save(&cert)
	// if updateresult.Error != nil {
	// 	return nil, errors.New("error updating cert " + uuid.String() + " with tag " + tagName)
	// }

	return cert, nil
}

//
// Tag functions
//

// CreateTag create a new Tag
func (certBackend *CertBackend) CreateTag(tagName string) (*Tag, error) {
	certBackend.logger.Debug("Creating Tag...")

	tag, err := NewTag(tagName)
	if err != nil {
		return nil, err
	}
	validationErr := certBackend.v.Validate(tag)
	if validationErr != nil {
		return nil, errors.New(strings.Join(validationErr.Errors(), "\n"))
	}

	// Create new entry
	result := certBackend.db.Create(&tag)
	if result.Error != nil {
		// log.Fatal(result.Error)
		return nil, result.Error
	}

	return tag, nil
}

// GetTagByName return Tag (without associated cert)
func (certBackend *CertBackend) GetTagByName(tagName string) (*Tag, error) {
	certBackend.logger.Debug("Getting Tag...")

	var tag Tag

	result := certBackend.db.Where("name = ?", tagName).First(&tag)
	if result.Error != nil {
		return nil, errors.New("Tag not found " + tagName)
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return nil, errors.New("Tag not found " + tagName)
	}

	return &tag, nil
}

// GetTagByID return Tag (without associated cert)
func (certBackend *CertBackend) GetTagByID(uuid uuid.UUID) (*Tag, error) {
	certBackend.logger.Debug("Getting Tag...")

	var tag Tag

	result := certBackend.db.Find(&tag, uuid)
	if result.Error != nil {
		return nil, errors.New("Tag not found " + uuid.String())
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return nil, errors.New("Tag not found " + uuid.String())
	}

	return &tag, nil
}
