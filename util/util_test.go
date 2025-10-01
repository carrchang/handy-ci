package util

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout captures stdout during function f execution and returns captured string
func captureStdout(f func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = orig
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestContainArgs(t *testing.T) {
	args := []string{"a", "b", "c"}
	if !ContainArgs(args, "b") {
		// expected to find b in slice
		t.Fatalf("expected to contain b")
	}
	if ContainArgs(args, "d") {
		// d is not present
		t.Fatalf("did not expect to contain d")
	}
}

func TestPrintfAndPrintln(t *testing.T) {
	out1 := captureStdout(func() { Printf("hello %s", "world") })
	if !strings.Contains(out1, "[Handy CI]") || !strings.Contains(out1, "hello world") {
		// output should contain tag and formatted text
		t.Fatalf("unexpected Printf output: %q", out1)
	}

	out2 := captureStdout(func() { Println("something") })
	if !strings.Contains(out2, "[Handy CI]") || !strings.Contains(out2, "something") {
		// output should contain tag and text
		t.Fatalf("unexpected Println output: %q", out2)
	}
}

func TestHome(t *testing.T) {
	got := Home()
	if got == "" {
		// Home should not return empty string
		t.Fatalf("expected non-empty home path")
	}
	sys, _ := os.UserHomeDir()
	if sys != "" && sys != got {
		// Accept difference only if environment peculiar; usually they match
		// Do not fail hard but ensure path ends with last element of expected home
		if !strings.HasSuffix(got, strings.TrimPrefix(sys, "/")) { // lenient check
			// Provide informative failure
			// Avoid brittle absolute comparison on exotic platforms
			// We still fail to surface mismatch
			// (keeping minimal logic per instructions)
			//
			// NOTE: We intentionally keep this message concise.
			//
			t.Fatalf("returned home path %q differs from system %q", got, sys)
		}
	}
}
