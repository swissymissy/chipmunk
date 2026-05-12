package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, ".env.example")
	dst := filepath.Join(dir, ".env")

	sourceContent := `JWT_SECRET=""
PROFESSOR_PASSWORD_HASH=""
PORT=8080
`

	if err := os.WriteFile(src, []byte(sourceContent), 0600); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile returned error: %v", err)
	}

	copiedBytes, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	if string(copiedBytes) != sourceContent {
		t.Errorf("copied content mismatch\nexpected:\n%s\ngot:\n%s", sourceContent, string(copiedBytes))
	}
}

func TestCopyFile_ReturnsErrorWhenSourceMissing(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "missing.env.example")
	dst := filepath.Join(dir, ".env")

	err := CopyFile(src, dst)
	if err == nil {
		t.Fatal("expected error when source file is missing, got nil")
	}
}

func TestReadEnvMap(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := `
# Chipmunk config
JWT_SECRET="abc123"
PROFESSOR_PASSWORD_HASH='$argon2id$v=19$m=65536,t=3,p=12$hashvalue'
PORT=8080
EMPTY_VALUE=
INVALID_LINE
`

	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	env, err := ReadEnvMap(envPath)
	if err != nil {
		t.Fatalf("ReadEnvMap returned error: %v", err)
	}

	if env["JWT_SECRET"] != "abc123" {
		t.Errorf("expected JWT_SECRET to be abc123, got %q", env["JWT_SECRET"])
	}

	expectedHash := "$argon2id$v=19$m=65536,t=3,p=12$hashvalue"
	if env["PROFESSOR_PASSWORD_HASH"] != expectedHash {
		t.Errorf("expected password hash %q, got %q", expectedHash, env["PROFESSOR_PASSWORD_HASH"])
	}

	if env["PORT"] != "8080" {
		t.Errorf("expected PORT to be 8080, got %q", env["PORT"])
	}

	if env["EMPTY_VALUE"] != "" {
		t.Errorf("expected EMPTY_VALUE to be empty, got %q", env["EMPTY_VALUE"])
	}

	if _, ok := env["INVALID_LINE"]; ok {
		t.Error("expected INVALID_LINE to be ignored")
	}
}

func TestReadEnvMap_ReturnsErrorWhenFileMissing(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	_, err := ReadEnvMap(envPath)
	if err == nil {
		t.Fatal("expected error when .env file is missing, got nil")
	}
}

func TestPatchEnvLine_ReplacesExistingKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := `# Chipmunk config
JWT_SECRET=""
PROFESSOR_PASSWORD_HASH=""
PORT=8080
`

	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	if err := PatchEnvLine(envPath, "JWT_SECRET", `"new-secret"`); err != nil {
		t.Fatalf("PatchEnvLine returned error: %v", err)
	}

	updatedBytes, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("failed to read updated .env: %v", err)
	}

	updated := string(updatedBytes)

	if !strings.Contains(updated, `JWT_SECRET="new-secret"`) {
		t.Errorf("expected JWT_SECRET to be replaced, got:\n%s", updated)
	}

	if !strings.Contains(updated, "# Chipmunk config") {
		t.Errorf("expected comment to be preserved, got:\n%s", updated)
	}

	if !strings.Contains(updated, "PROFESSOR_PASSWORD_HASH=\"\"") {
		t.Errorf("expected unrelated key to be preserved, got:\n%s", updated)
	}

	if !strings.Contains(updated, "PORT=8080") {
		t.Errorf("expected PORT to be preserved, got:\n%s", updated)
	}
}

func TestPatchEnvLine_AppendsMissingKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := `JWT_SECRET="abc"
PORT=8080
`

	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	if err := PatchEnvLine(envPath, "PROFESSOR_PASSWORD_HASH", `'$argon2id$hash'`); err != nil {
		t.Fatalf("PatchEnvLine returned error: %v", err)
	}

	updatedBytes, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("failed to read updated .env: %v", err)
	}

	updated := string(updatedBytes)

	if !strings.Contains(updated, `PROFESSOR_PASSWORD_HASH='$argon2id$hash'`) {
		t.Errorf("expected PROFESSOR_PASSWORD_HASH to be appended, got:\n%s", updated)
	}
}

func TestPatchEnvLine_DoesNotReplaceCommentedKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := `# JWT_SECRET="old-commented-secret"
PORT=8080
`

	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	if err := PatchEnvLine(envPath, "JWT_SECRET", `"real-secret"`); err != nil {
		t.Fatalf("PatchEnvLine returned error: %v", err)
	}

	updatedBytes, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("failed to read updated .env: %v", err)
	}

	updated := string(updatedBytes)

	if !strings.Contains(updated, `# JWT_SECRET="old-commented-secret"`) {
		t.Errorf("expected commented JWT_SECRET to remain unchanged, got:\n%s", updated)
	}

	if !strings.Contains(updated, `JWT_SECRET="real-secret"`) {
		t.Errorf("expected real JWT_SECRET to be appended, got:\n%s", updated)
	}
}

func TestPatchEnvLine_ReturnsErrorWhenFileMissing(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	err := PatchEnvLine(envPath, "JWT_SECRET", `"secret"`)
	if err == nil {
		t.Fatal("expected error when .env file is missing, got nil")
	}
}

func TestGenerateSecret(t *testing.T) {
	secret, err := generateSecret()
	if err != nil {
		t.Fatalf("generateSecret returned error: %v", err)
	}

	if secret == "" {
		t.Fatal("expected generated secret to not be empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		t.Fatalf("expected generated secret to be valid base64: %v", err)
	}

	if len(decoded) != 64 {
		t.Errorf("expected decoded secret to be 64 bytes, got %d", len(decoded))
	}
}

func TestGenerateSecret_ReturnsDifferentValues(t *testing.T) {
	secret1, err := generateSecret()
	if err != nil {
		t.Fatalf("generateSecret returned error: %v", err)
	}

	secret2, err := generateSecret()
	if err != nil {
		t.Fatalf("generateSecret returned error: %v", err)
	}

	if secret1 == secret2 {
		t.Fatal("expected two generated secrets to be different")
	}
}