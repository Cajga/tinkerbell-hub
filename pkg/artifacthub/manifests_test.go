package artifacthub_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/tinkerbell/actions/pkg/artifacthub"
)

func TestPopulateFromActionMarkdown(t *testing.T) {
	table := []struct {
		BaseManifest     *artifacthub.Manifest
		MarkdownFilePath string
		Expect           func(t *testing.T, m *artifacthub.Manifest, err error)
	}{
		{
			BaseManifest:     &artifacthub.Manifest{},
			MarkdownFilePath: "./testdata/happy-path-disk-wipe.md",
			Expect: func(t *testing.T, m *artifacthub.Manifest, err error) {
				t.Helper()
				if err != nil {
					t.Error(errors.Wrap(err, "unexpected error, this test should not return errors"))
				}
			},
		},
	}
	for _, s := range table {
		t.Run(s.MarkdownFilePath, func(t *testing.T) {
			md, err := ioutil.ReadFile(s.MarkdownFilePath)
			if err != nil {
				t.Error(err)
			}
			buf := bytes.NewBuffer(md)
			err = artifacthub.PopulateFromActionMarkdown(buf, s.BaseManifest)
			s.Expect(t, s.BaseManifest, err)
		})
	}
}

func TestWriteToFile(t *testing.T) {
	exp := `# this file is generated by ./cmd/gen/main.go
version: 0.1.0
name: test-slug
displayName: testslug
createdAt: ""
description: ""
logoPath: ""
provider:
    name: ""
`

	tmpDir := os.TempDir()
	err := artifacthub.WriteToFile(&artifacthub.Manifest{
		Name:        "test-slug",
		DisplayName: "testslug",
		Version:     "0.1.0",
	}, tmpDir)
	if err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadFile(path.Join(tmpDir, "test-slug", "0.1.0", "artifacthub-pkg.yml"))
	if err != nil {
		t.Error(err)
	}
	if dif := cmp.Diff(string(b), exp); dif != "" {
		t.Error(dif)
	}
}
