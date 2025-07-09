package file

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewOwnerDecoder(t *testing.T) {
	want := &OwnerDecoder{
		owners: map[string]owner{
			"AAA": {Company: "my company", City: "my city", Country: "my country"},
			"CUS": {Company: "my custom company", City: "my custom city", Country: "my custom country"},
		},
	}

	dir := t.TempDir()
	remotePath := filepath.Join(dir, "remote-owners.csv")
	customPath := filepath.Join(dir, "custom-owners.csv")

	if _, err := os.Stat(remotePath); err == nil {
		t.Errorf("NewOwnerDecoder() file should not exist: %s", remotePath)
		return
	}

	_ = os.WriteFile(customPath, []byte("CUS;my custom company;my custom city;my custom country"), 0o644)

	got, err := NewOwnerDecoder(remotePath, customPath)
	if err != nil {
		t.Errorf("NewOwnerDecoder() error = %v, want no err", err)
		return
	}
	if _, err := os.Stat(remotePath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("NewOwnerDecoder() file should exist: %s", remotePath)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewOwnerDecoder() got = %v, want %v", got, want)
	}
}
