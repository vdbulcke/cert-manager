package data

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/vdbulcke/cert-manager/data/database"
)

func TestMigrationDB(t *testing.T) {
	db := database.NewSqliteDB("/tmp/sqlite.db")
	err := db.AutoMigrate(&Tag{}, &Certificate{})
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}
	v := NewValidation()
	logger := hclog.Default()
	certBakcend := NewCertBackend(logger, db, v)
	tag, err := certBakcend.CreateTag("test_tag")
	if err != nil {
		t.Logf("Error creating Tag %s", err.Error())
		t.FailNow()
	}

	logger.Info("tag", tag)
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
	logger.Info("cert", cert)

	// set tag to cert
	certTag, err := certBakcend.SetCertTagNameByID(cert.ID, "test_tag")
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}

	logger.Info("cert tag", certTag)
	// Test Query
	foundCert, err := certBakcend.GetCertByID(cert.ID)
	if err != nil {
		t.Logf("Error creating Cert %s", err.Error())
		t.FailNow()
	}
	logger.Info("found Cert ", foundCert)
	t.FailNow()

}
