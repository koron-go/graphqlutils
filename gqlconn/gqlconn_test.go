package gqlconn

import "testing"

func TestGenerateNames(t *testing.T) {
	for i, c := range []struct {
		baseName  string
		expType  string
		expSmall string
	}{
		{"foo", "Foo", "foo"},
		{"Foo", "Foo", "foo"},
		{"bar", "Bar", "bar"},
		{"Bar", "Bar", "bar"},
	} {
		actType, actSmall := generateNames(c.baseName)
		if c.expType != actType {
			t.Fatalf("typeName mismatch #%d: want=%q got=%q", i, c.expType, actType)
		}
		if c.expSmall != actSmall {
			t.Fatalf("smallName mismatch #%d: want=%q got=%q", i, c.expSmall, actSmall)
		}
	}
}
