package data

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

//
// Tag CRUD functions
// Create: CreateTag, CreateTagWithDescription
// Read:   GetTagByName, GetTagByID, ListTags,  getTagByNameWithoutPreload
// Update: SetTagDescriptionByID
// Delete: DeleteTagByID,  DeleteTagPendingRecords
//

// CreateTag create a new Tag
func (certBackend *CertBackend) CreateTag(tagName string) (*Tag, error) {
	certBackend.logger.Debug("CreateTag: Creating Tag", "tagName", tagName)

	tag, err := NewTag(tagName)
	if err != nil {
		return nil, err
	}
	validationErr := certBackend.v.Validate(tag)
	if validationErr != nil {
		return nil, &DBObjectValidationError{Msg: strings.Join(validationErr.Errors(), "\n")}
	}

	// lookup if object already exist
	foundTag, _ := certBackend.getTagByNameWithoutPreload(tagName)
	if foundTag != nil {
		return foundTag, &DBObjectAlreadyExist{ID: foundTag.ID.String()}
	}

	// Create new entry
	result := certBackend.db.Create(&tag)
	if result.Error != nil {
		return nil, result.Error
	}

	return tag, nil
}

// CreateTagWithDescription create a new Tag with a description
func (certBackend *CertBackend) CreateTagWithDescription(tagName string, description string) (*Tag, error) {
	certBackend.logger.Debug("CreateTagWithDescription: Creating Tag", "tagName", tagName, "description", description)

	tag, err := NewTag(tagName)
	if err != nil {
		return nil, err
	}

	// setting description
	tag.Description = description

	validationErr := certBackend.v.Validate(tag)
	if validationErr != nil {
		return nil, &DBObjectValidationError{Msg: strings.Join(validationErr.Errors(), "\n")}
	}

	// lookup if object already exist
	foundTag, _ := certBackend.getTagByNameWithoutPreload(tagName)
	if foundTag != nil {
		return foundTag, &DBObjectAlreadyExist{ID: foundTag.ID.String()}
	}

	// Create new entry
	result := certBackend.db.Create(&tag)
	if result.Error != nil {
		return nil, result.Error
	}

	return tag, nil
}

// GetTagByName return Tag with associated Certs
func (certBackend *CertBackend) GetTagByName(tagName string) (*Tag, error) {
	certBackend.logger.Debug("GetTagByName: Getting Tag", "tagName", tagName)

	var tag Tag

	result := certBackend.db.Preload(clause.Associations).Where("name = ?", tagName).First(&tag)
	if result.Error != nil {
		return nil, errors.New("Tag not found " + tagName)
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return nil, &DBObjectNotFound{ID: tagName}
	}

	return &tag, nil
}

// getTagByNameWithoutPreload return Tag (without associated cert)
func (certBackend *CertBackend) getTagByNameWithoutPreload(tagName string) (*Tag, error) {
	certBackend.logger.Debug("getTagByNameWithoutPreload: Getting Tag...", "tagName", tagName)

	var tag Tag

	result := certBackend.db.Where("name = ?", tagName).First(&tag)
	if result.Error != nil {
		return nil, errors.New("Tag not found " + tagName)
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return nil, &DBObjectNotFound{ID: tagName}
	}

	return &tag, nil
}

// GetTagByID return Tag (without associated cert)
func (certBackend *CertBackend) GetTagByID(uuid uuid.UUID) (*Tag, error) {
	certBackend.logger.Debug("GetTagByID: Getting Tag...", "uuid", uuid)

	var tag Tag

	result := certBackend.db.Preload(clause.Associations).Find(&tag, uuid)
	if result.Error != nil {
		return nil, errors.New("Tag not found " + uuid.String())
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return nil, &DBObjectNotFound{ID: uuid.String()}

	}

	return &tag, nil
}

// ListTags returns all certs without their associated certs
func (certBackend *CertBackend) ListTags() (Tags, error) {
	certBackend.logger.Debug("ListTags: ")
	// List of Tags
	var tagList Tags

	result := certBackend.db.Find(&tagList)
	if result.Error != nil {
		return nil, result.Error
	}

	return tagList, nil
}

// SetTagDescriptionByID Update the Tag description
func (certBackend *CertBackend) SetTagDescriptionByID(uuid uuid.UUID, description string) (*Tag, error) {
	certBackend.logger.Debug("SetTagDescriptionByID: ", "uuid", uuid, "description", description)

	tag, err := certBackend.GetTagByID(uuid)
	if err != nil {
		return nil, err
	}
	// set descrition
	tag.Description = description

	// save change
	updateresult := certBackend.db.Save(&tag)
	if updateresult.Error != nil {
		return nil, updateresult.Error
	}

	return tag, nil
}

// DeleteTagByID delete the tag with uuid
// This is only doing a soft delete (update the DeleteAt field)
func (certBackend *CertBackend) DeleteTagByID(uuid uuid.UUID) error {
	certBackend.logger.Debug("DeleteTagByID: Deleting tag", "uuid", uuid)

	result := certBackend.db.Delete(&Tag{}, uuid)
	if result.Error != nil {
		return result.Error
	}

	// check if result is empty
	if result.RowsAffected == 0 {
		return &DBObjectNotFound{ID: uuid.String()}

	}

	return nil
}

// DeleteTagPendingRecords deletes all Tags that where flagged for deleting
func (certBackend *CertBackend) DeleteTagPendingRecords() error {
	certBackend.logger.Debug("DeleteTagPendingRecords: Deletings cert")

	// lookup all unscoped
	var unscopedTags Tags

	// lookup the entries with a deleted_at not null
	result := certBackend.db.Unscoped().Where("deleted_at IS NOT NULL").Find(&unscopedTags)
	if result.Error != nil {
		return result.Error
	}

	// Iterate over unscoped certificates
	// and delete them permanently
	for _, t := range unscopedTags {
		certBackend.logger.Error("DeleteTagPendingRecords: Deletings tag", "tag", t)

		// clearing all associated Tags with that cert
		certBackend.db.Model(&t).Association("Certificates").Clear()

		result := certBackend.db.Unscoped().Delete(&t)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
