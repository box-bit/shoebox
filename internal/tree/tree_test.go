package tree

import (
	"fmt"
	"os"
	"testing"

	"github.com/box-bit/shoebox/internal/testfs"
)

func TestEncodeFile(t *testing.T) {
	data := []struct {
		name    string
		file    string
		content string
		want    string
	}{
		{"empty file", "empty.txt", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"file with content", "file.txt", "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			f := testfs.CreateFile(t, d.file, d.content)
			file, err := os.Open(f)
			if err != nil {
				t.Fatalf("can't open file %q: %v", f, err)
			}
			defer file.Close()
			encoded, err := encodeFile(file)
			if err != nil {
				t.Fatalf("error occurred: %v", err)
			}

			if fmt.Sprintf("%x", encoded) != d.want {
				t.Errorf("encodeFile(%q) = %x, wants %x", f, d.want, fmt.Sprintf("%x", encoded))
			}
		})
	}
}
