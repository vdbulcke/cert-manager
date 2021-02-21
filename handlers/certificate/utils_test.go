package certificate

import "testing"

func TestSha(t *testing.T) {
	validSha := "ca42dd41745fd0b81eb902362cf9d8bf719da1bd1b1efc946f5b4c99f42c1b9e"

	validSHA256 := "9D:96:82:81:48:26:5A:B1:3B:3D:6D:B9:81:DB:11:64:3C:17:FF:6A:E0:91:E2:02:7C:09:92:82:45:B0:6A:D5"

	invalidSha := "9D96828148265AB13B3D6DB981DB11643C17FF6AE091E2027C09928245B06ADZ"

	//sha1
	invalidSha256 := "2E:7E:41:27:0F:E0:D9:A8:E4:5E:68:DC:89:64:5F:A5:D0:FB:47:BF"

	err := validateSHA256(validSha)
	if err != nil {
		t.Logf("Error this is a valid sha %s, err %s", validSha, err.Error())
		t.FailNow()
	}

	err = validateSHA256(invalidSha)
	if err == nil {
		t.Logf("Error this is not a valid sha %s", invalidSha)
		t.FailNow()
	}

	err = validateSHA256(invalidSha256)
	if err == nil {
		t.Logf("Error this is not a valid sha %s", invalidSha256)
		t.FailNow()
	}

	// testing conversion
	err = validateSHA256(convertSHA256(validSHA256))
	if err != nil {
		t.Logf("Error this is a valid sha %s, err %s", validSHA256, err.Error())
		t.FailNow()
	}

}
