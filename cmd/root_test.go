package cmd

import (
	"testing"
)

func TestExample(t *testing.T) {
	have := true
	want := true
	get := have
	if get != want {
		t.Fatalf("Execute, want: %v, got: %v", want, get)
	}

}
