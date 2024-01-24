package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestAddToFileAfterMagicString(t *testing.T) {
	type TestCase struct {
		content  string
		magic    string
		toAdd    string
		expected string
	}
	testCases := []TestCase{
		{
			content: `#Hello_World
def hello_world():
	print("Hello World")

# MAGIC_CALL_HELLO_WORLD
`,
			magic: "MAGIC_CALL_HELLO_WORLD\n",
			toAdd: "hello_world()\n",
			expected: `#Hello_World
def hello_world():
	print("Hello World")

# MAGIC_CALL_HELLO_WORLD
hello_world()
`,
		},
	}
	for i, testCase := range testCases {
		name := fmt.Sprintf("TEST_AddToFileAfterMagicString_%v.txt", rand.Int63())
		// Create file
		err := os.WriteFile(name, []byte(testCase.content), 0o644)
		if err != nil {
			t.Errorf("error creating file (testcase index %v). error %v", i, err)
		}
		err = AddToFileAfterMagicString(name, testCase.magic, testCase.toAdd)
		if err != nil {
			t.Errorf("got error %v adding to file after magic string (testcase index %v)", i, err)
		}
		content, err := os.ReadFile(name)
		strContent := string(content)
		if err != nil {
			t.Errorf("got error reading file content %v", err)
		}
		if strContent != testCase.expected {
			t.Errorf("expected %v got %v (testcase index %v )", testCase.expected, strContent, i)
		}
		// Remove file
		err = os.Remove(name)
		if err != nil {
			t.Errorf("error removing file (testcase index %v). error %v", i, err)
		}
	}
}
