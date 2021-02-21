package data

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/vdbulcke/cert-manager/data/database"
)

func TestMigrationDB(t *testing.T) {

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
		Name:       "cert-monitor",
		Level:      hclog.LevelFromString("INFO"),
		JSONFormat: true,
	})

	// create a cert backend
	certBakcend := NewCertBackend(logger, db, v)

	// create some tags
	tag, err := certBakcend.CreateTag("test_tag1")
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	logger.Debug("tag created", "tag", tag)
	_, err = certBakcend.CreateTag("test_tag2")
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	_, err = certBakcend.CreateTag("test_tag3")
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	// create cert
	pem := `
-----BEGIN CERTIFICATE-----
MIIDujCCAqKgAwIBAgILBAAAAAABD4Ym5g0wDQYJKoZIhvcNAQEFBQAwTDEgMB4G
A1UECxMXR2xvYmFsU2lnbiBSb290IENBIC0gUjIxEzARBgNVBAoTCkdsb2JhbFNp
Z24xEzARBgNVBAMTCkdsb2JhbFNpZ24wHhcNMDYxMjE1MDgwMDAwWhcNMjExMjE1
MDgwMDAwWjBMMSAwHgYDVQQLExdHbG9iYWxTaWduIFJvb3QgQ0EgLSBSMjETMBEG
A1UEChMKR2xvYmFsU2lnbjETMBEGA1UEAxMKR2xvYmFsU2lnbjCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAKbPJA6+Lm8omUVCxKs+IVSbC9N/hHD6ErPL
v4dfxn+G07IwXNb9rfF73OX4YJYJkhD10FPe+3t+c4isUoh7SqbKSaZeqKeMWhG8
eoLrvozps6yWJQeXSpkqBy+0Hne/ig+1AnwblrjFuTosvNYSuetZfeLQBoZfXklq
tTleiDTsvHgMCJiEbKjNS7SgfQx5TfC4LcshytVsW33hoCmEofnTlEnLJGKRILzd
C9XZzPnqJworc5HGnRusyMvo4KD0L5CLTfuwNhv2GXqF4G3yYROIXJ/gkwpRl4pa
zq+r1feqCapgvdzZX99yqWATXgAByUr6P6TqBwMhAo6CygPCm48CAwEAAaOBnDCB
mTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUm+IH
V2ccHsBqBt5ZtJot39wZhi4wNgYDVR0fBC8wLTAroCmgJ4YlaHR0cDovL2NybC5n
bG9iYWxzaWduLm5ldC9yb290LXIyLmNybDAfBgNVHSMEGDAWgBSb4gdXZxwewGoG
3lm0mi3f3BmGLjANBgkqhkiG9w0BAQUFAAOCAQEAmYFThxxol4aR7OBKuEQLq4Gs
J0/WwbgcQ3izDJr86iw8bmEbTUsp9Z8FHSbBuOmDAGJFtqkIk7mpM0sYmsL4h4hO
291xNBrBVNpGP+DTKqttVCL1OmLNIG+6KYnX3ZHu01yiPqFbQfXf5WRDLenVOavS
ot+3i9DAgBkcRcAtjOj4LaR0VknFBbVPFd5uRHg5h6h+u/N5GJG79G+dwfCMNYxd
AfvDbbnvRG15RjF+Cv6pgsH/76tuIMRQyV+dTZsXjAzlAcmgQWpzU/qlULRuJQ/7
TBj0/VLZjmmx6BEP3ojY+x1J96relc8geMJgEtslQIxq/H5COEBkEveegeGTLg==
-----END CERTIFICATE-----`

	for _, initT := range []string{"init_tag2", "init_tag3"} {
		_, err = certBakcend.CreateTag(initT)
		if err != nil {
			t.Logf("Expecting Error when creating Tag with existing name")
			t.FailNow()
		}
	}

	// cert, err := certBakcend.CreateCertificate(pem)
	cert, err := certBakcend.CreateCertificateWithTags(pem, []string{"init_tag2", "init_tag3"})
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	logger.Debug("cert created", "cert", cert)

	// set tag to cert
	cert, err = certBakcend.SetCertTagNameByID(cert.ID, "test_tag1")
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	logger.Debug("assigned 'test_tag1'", "cert", cert)

	// set tag list to cert
	tagList := []string{"test_tag2", "test_tag3"}
	cert, err = certBakcend.SetCertTagsNameByID(cert.ID, tagList)
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	logger.Info("assigned ['test_tag2', 'test_tag3']", "cert.Tags", cert.Tags)

	// Test Query
	foundCert, err := certBakcend.GetCertByID(cert.ID)
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}
	logger.Debug("found Cert", "cert", foundCert)

	tagRes, err := certBakcend.GetTagByName("test_tag1")
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}
	logger.Debug("found tag", "tag", tagRes)
	certs, err := certBakcend.ListCerts()
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	for _, c := range certs {
		logger.Error("List cert", "c", c)
	}

	// list
	tags, err := certBakcend.ListTags()
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	for _, c := range tags {
		logger.Error("List tag", "c", c)
	}
	// deleting cert

	// err = certBakcend.DeleteCertByID(foundCert.ID)
	// if err != nil {
	// 	t.Logf("Error creating Cert %s", err.Error())
	// 	t.FailNow()
	// }

	// err = certBakcend.DeletePendingRecords()
	// if err != nil {
	// 	t.Logf("Error creating Cert %s", err.Error())
	// 	t.FailNow()
	// }

	err = certBakcend.DeleteTagByID(tag.ID)
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	err = certBakcend.DeleteTagPendingRecords()
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	// list
	tags, err = certBakcend.ListTags()
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	for _, c := range tags {
		logger.Error("List tag", "c", c)
	}

	//  lookup cert

	foundCert, err = certBakcend.GetCertByID(cert.ID)
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}
	logger.Error("found Cert", "cert", foundCert)

	// just used to display logger
	// t.FailNow()

}

// Actual Test Cases

func TestCertificatesTag(t *testing.T) {

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

	// create cert
	pem := `
-----BEGIN CERTIFICATE-----
MIIDujCCAqKgAwIBAgILBAAAAAABD4Ym5g0wDQYJKoZIhvcNAQEFBQAwTDEgMB4G
A1UECxMXR2xvYmFsU2lnbiBSb290IENBIC0gUjIxEzARBgNVBAoTCkdsb2JhbFNp
Z24xEzARBgNVBAMTCkdsb2JhbFNpZ24wHhcNMDYxMjE1MDgwMDAwWhcNMjExMjE1
MDgwMDAwWjBMMSAwHgYDVQQLExdHbG9iYWxTaWduIFJvb3QgQ0EgLSBSMjETMBEG
A1UEChMKR2xvYmFsU2lnbjETMBEGA1UEAxMKR2xvYmFsU2lnbjCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAKbPJA6+Lm8omUVCxKs+IVSbC9N/hHD6ErPL
v4dfxn+G07IwXNb9rfF73OX4YJYJkhD10FPe+3t+c4isUoh7SqbKSaZeqKeMWhG8
eoLrvozps6yWJQeXSpkqBy+0Hne/ig+1AnwblrjFuTosvNYSuetZfeLQBoZfXklq
tTleiDTsvHgMCJiEbKjNS7SgfQx5TfC4LcshytVsW33hoCmEofnTlEnLJGKRILzd
C9XZzPnqJworc5HGnRusyMvo4KD0L5CLTfuwNhv2GXqF4G3yYROIXJ/gkwpRl4pa
zq+r1feqCapgvdzZX99yqWATXgAByUr6P6TqBwMhAo6CygPCm48CAwEAAaOBnDCB
mTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUm+IH
V2ccHsBqBt5ZtJot39wZhi4wNgYDVR0fBC8wLTAroCmgJ4YlaHR0cDovL2NybC5n
bG9iYWxzaWduLm5ldC9yb290LXIyLmNybDAfBgNVHSMEGDAWgBSb4gdXZxwewGoG
3lm0mi3f3BmGLjANBgkqhkiG9w0BAQUFAAOCAQEAmYFThxxol4aR7OBKuEQLq4Gs
J0/WwbgcQ3izDJr86iw8bmEbTUsp9Z8FHSbBuOmDAGJFtqkIk7mpM0sYmsL4h4hO
291xNBrBVNpGP+DTKqttVCL1OmLNIG+6KYnX3ZHu01yiPqFbQfXf5WRDLenVOavS
ot+3i9DAgBkcRcAtjOj4LaR0VknFBbVPFd5uRHg5h6h+u/N5GJG79G+dwfCMNYxd
AfvDbbnvRG15RjF+Cv6pgsH/76tuIMRQyV+dTZsXjAzlAcmgQWpzU/qlULRuJQ/7
TBj0/VLZjmmx6BEP3ojY+x1J96relc8geMJgEtslQIxq/H5COEBkEveegeGTLg==
-----END CERTIFICATE-----`
	cert, err := certBakcend.CreateCertificate(pem)
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	logger.Debug("cert created", "cert", cert)
	tagList := []string{"test_tag2", "test_tag3"}
	cert, err = certBakcend.SetCertTagsNameByID(cert.ID, tagList)
	if err == nil {
		t.Logf("Expecting Error creating Cert with non existing tag")
		t.FailNow()
	}

}

func TestError(t *testing.T) {

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
		Name:       "cert-monitor",
		Level:      hclog.LevelFromString("INFO"),
		JSONFormat: true,
	})

	// create a cert backend
	certBakcend := NewCertBackend(logger, db, v)

	// create cert
	nonValidPem := `
-----BEGIN CERTIFICATE-----
MIIDujCCAqKgAwIBAgILBAAAAAABD4Ym5g0wDQYJKoZIhvcNAQEFBQAwTDEgMB4G
A1UECxMXR2xvYmFsU2lnbiBSb290IENBIC0gUjIxEzARBgNVBAoTCkdsb2JhbFNp
Z24xEzARBgNVBAMTCkdsb2JhbFNpZ24wHhcNMDYxMjE1MDgwMDAwWhcNMjExMjE1
MDgwMDAwWjBMMSAwHgYDVQQLExdHbG9iYWxTaWduIFJvb3QgQ0EgLSBSMjETMBEG
A1UEChMKR2xvYmFsU2lnbjETMBEGA1UEAxMKR2xvYmFsU2lnbjCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAKbPJA6+Lm8omUVCxKs+IVSbC9N/hHD6ErPL
v4dfxn+G07IwXNb9rfF73OX4YJYJkhD10FPe+3t+c4isUoh7SqbKSaZeqKeMWhG8
C9XZzPnqJworc5HGnRusyMvo4KD0L5CLTfuwNhv2GXqF4G3yYROIXJ/gkwpRl4pa
3lm0mi3f3BmGLjANBgkqhkiG9w0BAQUFAAOCAQEAmYFThxxol4aR7OBKuEQLq4Gs
J0/WwbgcQ3izDJr86iw8bmEbTUsp9Z8FHSbBuOmDAGJFtqkIk7mpM0sYmsL4h4hO
291xNBrBVNpGP+DTKqttVCL1OmLNIG+6KYnX3ZHu01yiPqFbQfXf5WRDLenVOavS
ot+3i9DAgBkcRcAtjOj4LaR0VknFBbVPFd5uRHg5h6h+u/N5GJG79G+dwfCMNYxd
AfvDbbnvRG15RjF+Cv6pgsH/76tuIMRQyV+dTZsXjAzlAcmgQWpzU/qlULRuJQ/7
TBj0/VLZjmmx6BEP3ojY+x1J96relc8geMJgEtslQIxq/H5COEBkEveegeGTLg==
-----END CERTIFICATE-----`

	// cert, err := certBakcend.CreateCertificate(pem)
	cert, err := certBakcend.CreateCertificateWithTags(nonValidPem, []string{"init_tag2", "init_tag3"})
	if err != nil {
		if _, ok := err.(*DBObjectValidationError); ok {
			logger.Debug("cert created", "validationError", err)
		} else {
			logger.Error("cert created", "validationError", err, "type", err)
			t.Logf("Epxecting validation error creating got %s", err.Error())
			t.FailNow()
		}

	}

	logger.Debug("cert created", "cert", cert)

	// t.FailNow()

}
