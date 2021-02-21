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

//
// Cert CRUD functions
// Create: CreateCertificate, CreateCertificateWithTags
// Read:   GetCertByID, GetCertByFingerprint, ListCerts
// Update: SetCertTagNameByID, SetCertTagsNameByID
// Delete: DeleteCertByID, DeleteCertPendingRecords
//

// CreateCertificate create a new Certificate from pem and insert it into DB
func (certBackend *CertBackend) CreateCertificate(pem string) (*Certificate, error) {
	certBackend.logger.Debug("CreateCertificate: Creating Cert...")

	cert, err := NewCertificate(pem)
	if err != nil {
		return nil, err
	}

	validationErr := certBackend.v.Validate(cert)
	if validationErr != nil {
		return nil, &DBObjectValidationError{Msg: strings.Join(validationErr.Errors(), "\n")}
	}

	// lookup if cert already exist
	foundCert, _ := certBackend.GetCertByFingerprint(cert.SHA256)
	if foundCert != nil {
		return foundCert, &DBObjectAlreadyExist{ID: foundCert.ID.String()}
	}

	// Create new entry
	result := certBackend.db.Create(&cert)
	if result.Error != nil {
		// log.Fatal(result.Error)
		return nil, result.Error
	}

	return cert, nil
}

// CreateCertificateWithTags create a new Certificate from pem and assigned some tags
func (certBackend *CertBackend) CreateCertificateWithTags(pem string, tags []string) (*Certificate, error) {
	certBackend.logger.Debug("CreateCertificateWithTags: Creating Cert...")

	// Create cert in db
	cert, err := certBackend.CreateCertificate(pem)
	if err != nil {
		return nil, err
	}

	// for simplicity use update method SetCertTagsNameByID on
	// newly created cert
	return certBackend.SetCertTagsNameByID(cert.ID, tags)
}

// GetCertByID return Tag (without associated cert)
func (certBackend *CertBackend) GetCertByID(uuid uuid.UUID) (*Certificate, error) {
	certBackend.logger.Debug("GetCertByID: Getting Certificate...")

	var cert Certificate
	// preload the tags
	result := certBackend.db.Preload(clause.Associations).Find(&cert, uuid)
	if result.Error != nil {
		return nil, errors.New("Cert not found " + uuid.String())
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return nil, &DBObjectNotFound{ID: uuid.String()}
	}

	return &cert, nil
}

// GetCertByFingerprint return Tag (without associated cert)
func (certBackend *CertBackend) GetCertByFingerprint(sha256 string) (*Certificate, error) {
	certBackend.logger.Debug("GetCertByFingerprint: Getting Certificate...")

	var cert Certificate

	result := certBackend.db.Preload("Tags").Where("sha256 = ?", sha256).First(&cert)
	if result.Error != nil {
		return nil, errors.New("Cert not found " + sha256)
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return nil, &DBObjectNotFound{ID: sha256}
	}

	return &cert, nil
}

// ListCerts returns all certs with their assosicated tags
func (certBackend *CertBackend) ListCerts() (Certificates, error) {

	// List of Certificates
	var certList Certificates

	result := certBackend.db.Preload(clause.Associations).Find(&certList)
	if result.Error != nil {
		return nil, result.Error
	}

	return certList, nil
}

// SetCertTagNameByID return Tag (without associated cert)
func (certBackend *CertBackend) SetCertTagNameByID(uuid uuid.UUID, tagName string) (*Certificate, error) {
	certBackend.logger.Debug("SetCertTagNameByID: Setting tag for Certificate...")

	// call SetCertTagsNameByID with list containing only a single tag
	return certBackend.SetCertTagsNameByID(uuid, []string{tagName})
}

// SetCertTagsNameByID update cert (uuid) with list of tags name
func (certBackend *CertBackend) SetCertTagsNameByID(uuid uuid.UUID, tagList []string) (*Certificate, error) {
	certBackend.logger.Debug("SetCertTagsNameByID: Setting tags for Certificate...", "tagList", tagList)

	// get cert by uuid
	cert, cerr := certBackend.GetCertByID(uuid)
	if cerr != nil {
		return nil, cerr
	}

	// getting exitsing  tags
	updatedTags := cert.Tags
	for _, t := range tagList {
		// lookup tag
		tag, err := certBackend.getTagByNameWithoutPreload(t)
		if err != nil {
			// if one tag does not exist abort update
			return nil, err
		}

		certBackend.logger.Debug("SetCertTagsNameByID: go tag", "tagName", t, "tag", tag)
		// if not in tags list append to list
		if !certBackend.isTagAlreadyInList(tag, cert.Tags) {
			updatedTags = append(updatedTags, *tag)
			certBackend.logger.Debug("SetCertTagsNameByID: ", "updatedTags", updatedTags, "cert.Tags", cert.Tags)
		}
	}

	// update the cert with tag new set of tags
	tagErr := certBackend.db.Model(&cert).Association("Tags").Replace(updatedTags)
	if tagErr != nil {
		return nil, errors.New("error updating cert " + uuid.String() + " with tags " + strings.Join(tagList, ","))
	}

	return cert, nil
}

// DeleteCertificateTagsByID remove  list of tags name from cert (uuid)
func (certBackend *CertBackend) DeleteCertificateTagsByID(uuid uuid.UUID, tagList []string) (*Certificate, error) {
	certBackend.logger.Debug("DeleteCertificateTagByID: removing tags from Certificate...", "tagList", tagList)

	// get cert by uuid
	cert, cerr := certBackend.GetCertByID(uuid)
	if cerr != nil {
		return nil, cerr
	}

	// lookup  tags from list
	var resolvedTags []Tag
	for _, t := range tagList {
		// lookup tag
		tag, err := certBackend.getTagByNameWithoutPreload(t)
		if err != nil {
			// if one tag does not exist abort update
			return nil, err
		}

		resolvedTags = append(resolvedTags, *tag)

	}

	// // getting exitsing  tags
	// var updatedTags []Tag

	// for _, existingTag := range cert.Tags {

	// 	// check if existing tag is not in the list of tags to be removed
	// 	if !certBackend.isTagAlreadyInList(&existingTag, resolvedTags) {
	// 		updatedTags = append(updatedTags)
	// 	}

	// }

	// update the cert with tag new set of tags
	// tagErr := certBackend.db.Model(&cert).Association("Tags").Replace(updatedTags)
	tagErr := certBackend.db.Model(&cert).Association("Tags").Delete(resolvedTags)
	if tagErr != nil {
		return nil, errors.New("error updating cert " + uuid.String() + " with tags " + strings.Join(tagList, ","))
	}

	return cert, nil
}

// DeleteCertByID delete the certificate with uuid
// This is only doing a soft delete (update the DeleteAt field)
func (certBackend *CertBackend) DeleteCertByID(uuid uuid.UUID) error {
	certBackend.logger.Debug("DeleteCertByID: Deleting cert", "uuid", uuid)

	result := certBackend.db.Delete(&Certificate{}, uuid)
	if result.Error != nil {
		return result.Error
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return &DBObjectNotFound{ID: uuid.String()}
	}

	return nil
}

// DeleteCertPendingRecords deletes all Certificates that where flagged for deleting
func (certBackend *CertBackend) DeleteCertPendingRecords() error {
	certBackend.logger.Debug("DeleteCertPendingRecords: Deletings cert")

	// lookup all unscoped
	var unscopedCertificates Certificates

	// lookup the entries with a deleted_at not null
	result := certBackend.db.Unscoped().Where("deleted_at IS NOT NULL").Find(&unscopedCertificates)
	if result.Error != nil {
		return result.Error
	}

	// Iterate over unscoped certificates
	// and delete them permanently
	for _, c := range unscopedCertificates {
		certBackend.logger.Debug("DeleteCertPendingRecords: Deletings cert", "cert", c)

		// clearing all associated Tags with that cert
		certBackend.db.Model(&c).Association("Tags").Clear()

		result := certBackend.db.Unscoped().Delete(&c)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

//
// Helper function
//

// isTagAlreadyInList return true if tag in the tagList
func (certBackend *CertBackend) isTagAlreadyInList(tag *Tag, tagList []Tag) bool {

	// check if tag already in list
	for _, t := range tagList {
		if t.ID == tag.ID {
			return true
		}
	}

	return false
}
