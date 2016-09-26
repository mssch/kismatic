package inspector

import (
	"io/ioutil"
	"testing"
)

func TestFileContentCheckFileDoesntExist(t *testing.T) {
	c := FileContentCheck{
		File:         "doesntExist",
		SearchString: "foo",
	}
	err := c.Check()
	if err == nil {
		t.Errorf("Expected an error, but didn't get one")
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
	if c.Check() != nil {
		t.Errorf("The check failed when we were expecting success. File content: %s\n Search String: %s", fileConts, c.SearchString)
	}
}
