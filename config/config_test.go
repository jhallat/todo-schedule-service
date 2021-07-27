package config

import "testing"

func TestTokenize(t *testing.T) {

	actual := tokenizeCommandLine(" -P username=test -P password = something")
	expected := []string{"-P", "username", "=", "test", "-P", "password", "=", "something"}
	if len(actual) != len(expected) {
		t.Fatalf("Expected array length is %d, but received %d", len(expected), len(actual))
	}
	for index, value := range expected {
		if actual[index] != value {
			t.Errorf("Expected value '%s' at index %d but got '%s'", value, index, actual[index])
		}
	}
}

func TestParseTokens(t *testing.T) {

	tokens := []string{"-P", "username", "=", "test", "-P", "password", "=", "something"}
	actual, err := parseTokens(tokens)
	if err != nil {
		t.Fatalf("Received unexpected error %s", err)
	}
	if actual["username"] != "test" {
		t.Errorf("Expected value %s for 'username' but got '%s'", "test", actual["username"])
	}
	if actual["password"] != "something" {
		t.Errorf("Expected value %s for 'password' but got '%s'", "test", actual["password"])
	}

}

func TestParseCommandLine(t *testing.T) {
	args := []string{"config-test.exe","-P","username=test", "-P", "password=something"}
	actual := make(map[string]string)
	err := ParseCommandLine(args, actual)
	if err != nil {
		t.Fatalf("Received unexpected error %s", err)
	}
	if actual["username"] != "test" {
		t.Errorf("Expected value %s for 'username' but got '%s'", "test", actual["username"])
	}
	if actual["password"] != "something" {
		t.Errorf("Expected value %s for 'password' but got '%s'", "test", actual["password"])
	}
}