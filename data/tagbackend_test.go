package data

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/vdbulcke/cert-manager/data/database"
)

func TestTags(t *testing.T) {

	// clean up prev test
	dbInfo, err := os.Stat("/tmp/sqlite.db")
	if dbInfo != nil {
		fmt.Println("clearing db")
		os.Remove("/tmp/sqlite.db")
	}
	// init new db
	db := database.NewSqliteDB("/tmp/sqlite.db")
	err = db.AutoMigrate(&Tag{}, &Certificate{})
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}
	// setup logger and validation
	v := NewValidation()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "cert-manager",
		Level:      hclog.LevelFromString("INFO"),
		JSONFormat: true,
	})

	// create a cert backend
	certBakcend := NewCertBackend(logger, db, v)

	// create some tags
	tag1, err := certBakcend.CreateTag("test_tag1")
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	logger.Debug("tag created", "tag", tag1)

	tag2, err := certBakcend.CreateTag("test_tag2")
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	tag3, err := certBakcend.CreateTagWithDescription("test_tag3", "Description of Tag3")
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	// Creating tag with existing name should return an error
	_, err = certBakcend.CreateTag("test_tag3")
	if err == nil {
		t.Logf("Expecting Error when creating Tag with existing name test_tag3")
		t.FailNow()
	}

	// Testing Getter
	getTag1, err := certBakcend.GetTagByName("test_tag1")
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	// check if equals
	if tag1.ID != getTag1.ID {
		logger.Error("Expecting tag to be the same", "tag1", tag1, "getTag1", getTag1)
		t.Logf("Expecting tag to be the same ")
		t.FailNow()
	}

	// testing setting description
	updatedTag2, err := certBakcend.SetTagDescriptionByID(tag2.ID, "Description tag2")
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	descTag2, err := certBakcend.GetTagByID(tag2.ID)
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	if updatedTag2.Description != descTag2.Description {

		t.Logf("Error setting description Tag '%s' does not match '%s'", updatedTag2.Description, descTag2.Description)
		t.FailNow()

	}

	// check list
	tagList, err := certBakcend.ListTags()
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	checkTag1 := false
	checkTag2 := false
	checkTag3 := false
	for _, t := range tagList {

		if t.ID == tag1.ID {
			checkTag1 = true
		}

		if t.ID == tag2.ID {
			checkTag2 = true
		}

		if t.ID == tag3.ID {
			checkTag3 = true
		}
	}

	if !(checkTag1 && checkTag2 && checkTag3) {
		t.Logf("Error listing Tag")
		t.FailNow()
	}

	err = certBakcend.DeleteTagByID(tag1.ID)
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	// check delete
	deletedTag1, err := certBakcend.GetTagByName("test_tag1")
	if err == nil {
		logger.Error("Found tag after delete", "deletedTag1", deletedTag1)
		t.Logf("Error expeting error looking up tag after delete ")
		t.FailNow()
	}

}
