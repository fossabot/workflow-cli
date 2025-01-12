package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/workflow-cli/pkg/git"
	"github.com/drycc/workflow-cli/settings"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	name, err := ioutil.TempDir("", "drycc-cli-unit-test-load")
	assert.NoError(t, err)
	defer os.RemoveAll(name)

	filename := filepath.Join(name, "test.json")
	host := "drycc.example.com"
	client, err := drycc.New(false, host, "")
	assert.NoError(t, err)

	config := settings.Settings{
		Username: "test",
		Client:   client,
	}

	filename, err = config.Save(filename)
	assert.NoError(t, err)

	_, appID, err := load(filename, "test")
	assert.NoError(t, err)
	assert.Equal(t, appID, "test", "app")

	assert.NoError(t, os.Chdir(name))

	_, appID, err = load(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, appID, filepath.Base(name), "app")

	assert.NoError(t, git.Init(git.DefaultCmd))
	assert.NoError(t, git.CreateRemote(git.DefaultCmd, host, "drycc", "testing"))

	_, appID, err = load(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, appID, "testing", "app")
}

func TestDrinkOfChoice(t *testing.T) {
	os.Setenv("DRYCC_DRINK_OF_CHOICE", "test")
	assert.Equal(t, drinkOfChoice(), "test", "output")
	os.Unsetenv("DRYCC_DRINK_OF_CHOICE")
	assert.Equal(t, drinkOfChoice(), "coffee", "output")
}

func TestLimitsCount(t *testing.T) {
	t.Parallel()
	assert.Equal(t, limitCount(1, 1), "\n", "output")
	assert.Equal(t, limitCount(1, 2), " (1 of 2)\n", "output")
}

func TestAPICompatibility(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	cmdr := DryccCmd{WErr: &b, ConfigFile: ""}
	client := drycc.Client{ControllerAPIVersion: "v1.0"}

	err := cmdr.checkAPICompatibility(&client, drycc.ErrAPIMismatch)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `!    WARNING: Client and server API versions do not match. Please consider upgrading.
!    Client version: 2.3
!    Server version: v1.0
`, "output")

	// After being warned once, the warning should not be printed again.
	b.Reset()
	err = cmdr.checkAPICompatibility(&client, drycc.ErrAPIMismatch)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "", "output")

	b.Reset()
	err = cmdr.checkAPICompatibility(&client, drycc.ErrConflict)
	assert.Error(t, drycc.ErrConflict, err)
	assert.Equal(t, b.String(), "", "output")
}
