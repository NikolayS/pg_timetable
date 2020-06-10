package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCommander struct{}

// overwrite CombinedOutput function of os/exec so only parameter syntax and return codes are checked...
func (c testCommander) CombinedOutput(ctx context.Context, command string, args ...string) ([]byte, error) {
	if strings.HasPrefix(command, "ping") {
		return []byte(fmt.Sprint(command, args)), nil
	}
	return []byte(fmt.Sprintf("Command %s not found", command)), &exec.Error{Name: command, Err: exec.ErrNotFound}
}

func TestShellCommand(t *testing.T) {
	cmd = testCommander{}
	var err error
	var out string
	var retCode int

	ctx := context.Background()

	_, _, err = executeShellCommand(ctx, "", []string{""})
	assert.EqualError(t, err, "Shell command cannot be empty", "Empty command should out, fail")

	_, out, err = executeShellCommand(ctx, "ping0", nil)
	assert.NoError(t, err, "Command with nil param is out, OK")
	assert.True(t, strings.HasPrefix(string(out), "ping0"), "Output should containt only command ")

	_, _, err = executeShellCommand(ctx, "ping1", []string{})
	assert.NoError(t, err, "Command with empty array param is OK")

	_, _, err = executeShellCommand(ctx, "ping2", []string{""})
	assert.NoError(t, err, "Command with empty string param is OK")

	_, _, err = executeShellCommand(ctx, "ping3", []string{"[]"})
	assert.NoError(t, err, "Command with empty json array param is OK")

	_, _, err = executeShellCommand(ctx, "ping3", []string{"[null]"})
	assert.NoError(t, err, "Command with nil array param is OK")

	_, _, err = executeShellCommand(ctx, "ping4", []string{`["localhost"]`})
	assert.NoError(t, err, "Command with one param is OK")

	_, _, err = executeShellCommand(ctx, "ping5", []string{`["localhost", "-4"]`})
	assert.NoError(t, err, "Command with many params is OK")

	_, _, err = executeShellCommand(ctx, "pong", nil)
	assert.IsType(t, (*exec.Error)(nil), err, "Uknown command should produce error")

	retCode, _, err = executeShellCommand(ctx, "ping5", []string{`{"param1": "localhost"}`})
	assert.IsType(t, (*json.UnmarshalTypeError)(nil), err, "Command should fail with mailformed json parameter")
	assert.NotEqual(t, 0, retCode, "return code should indicate failure.")
}
