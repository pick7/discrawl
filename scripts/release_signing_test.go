package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCodesignReleaseBinarySkipsCredentialFreeBuilds(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("codesign-release-binary.sh is a POSIX shell script")
	}

	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash is not available")
	}

	scriptDir := signingScriptDir(t)
	cmd := exec.CommandContext(t.Context(), bash, "./codesign-release-binary.sh", "darwin_arm64", "/does/not/exist")
	cmd.Dir = scriptDir
	cmd.Env = signingTestEnv("DISCRAWL_CODESIGN_REQUIRED=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("credential-free signing hook failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "skipping Developer ID signing") {
		t.Fatalf("missing snapshot skip notice: %s", output)
	}
}

func TestCodesignReleaseBinaryFailsClosedForOfficialBuilds(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("codesign-release-binary.sh is a POSIX shell script")
	}

	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash is not available")
	}

	scriptDir := signingScriptDir(t)
	cmd := exec.CommandContext(t.Context(), bash, "./codesign-release-binary.sh", "darwin_arm64", "/does/not/exist")
	cmd.Dir = scriptDir
	cmd.Env = signingTestEnv("DISCRAWL_CODESIGN_REQUIRED=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("official signing hook unexpectedly succeeded: %s", output)
	}
	want := "CODESIGN_IDENTITY is required"
	if runtime.GOOS != "darwin" {
		want = "official macOS release signing must run on macOS"
	}
	if !strings.Contains(string(output), want) {
		t.Fatalf("unexpected official signing failure: %s", output)
	}
}

func signingScriptDir(t *testing.T) string {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate test file")
	}
	return filepath.Dir(testFile)
}

func signingTestEnv(extra string) []string {
	env := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "CODESIGN_IDENTITY=") || strings.HasPrefix(entry, "DISCRAWL_CODESIGN_REQUIRED=") {
			continue
		}
		env = append(env, entry)
	}
	return append(env, extra)
}
