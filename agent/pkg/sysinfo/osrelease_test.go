package sysinfo

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOSReleaseWithPrettyName(t *testing.T) {
	f, err := ioutil.TempFile("", "os-release")
	assert.NoError(t, err)

	defer os.Remove(f.Name())

	_, err = f.Write([]byte("ID=ubuntu\nNAME=Ubuntu\nPRETTY_NAME=\"Ubuntu 16.04 LTS\""))
	assert.NoError(t, err)

	err = f.Close()
	assert.NoError(t, err)

	DefaultOSReleaseFilename = f.Name()

	osrelease, err := GetOSRelease()
	assert.NoError(t, err)

	assert.Equal(t, &OSRelease{"ubuntu", "Ubuntu 16.04 LTS"}, osrelease)
}

func TestGetOSReleaseWithoutPrettyName(t *testing.T) {
	f, err := ioutil.TempFile("", "os-release")
	assert.NoError(t, err)

	defer os.Remove(f.Name())

	_, err = f.Write([]byte("ID=ubuntu\nNAME=Ubuntu"))
	assert.NoError(t, err)

	err = f.Close()
	assert.NoError(t, err)

	DefaultOSReleaseFilename = f.Name()

	osrelease, err := GetOSRelease()
	assert.NoError(t, err)

	assert.Equal(t, &OSRelease{"ubuntu", "Ubuntu"}, osrelease)
}
