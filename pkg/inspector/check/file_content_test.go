package check

import (
	"io/ioutil"
	"testing"
)

func TestFileContentCheckFileDoesntExist(t *testing.T) {
	c := FileContentCheck{
		File:         "doesntExist",
		SearchString: "foo",
	}
	ok, err := c.Check()
	if err == nil {
		t.Errorf("expected an error but didn't get one")
	}
	if ok {
		t.Errorf("check returned true for a non-existent file")
	}
}

func TestFileContentCheck(t *testing.T) {
	f, err := ioutil.TempFile("", "file-regex-check")
	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}
	fileConts := "hello world\n"
	f.WriteString(fileConts)
	c := FileContentCheck{
		File:         f.Name(),
		SearchString: "^hello w.*",
	}
	ok, err := c.Check()
	if err != nil {
		t.Errorf("Unexpected error when running check: %v", err)
	}
	if !ok {
		t.Errorf("Expected check OK, but check failed. Search string was: %q\nFile contents: \n%s\n", c.SearchString, fileConts)
	}
}
