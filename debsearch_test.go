package debsearch

import (
    "testing"
)

func Test001(t *testing.T) {
    expected := "Hello debsearch v0.1.0\n"
    actual := Hello()
    if actual != expected {
        t.Errorf("expected %q, got %q", expected, actual)
    }
}
