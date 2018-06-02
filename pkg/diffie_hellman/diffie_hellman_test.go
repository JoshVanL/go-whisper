package diffie_hellman

import (
	"testing"
)

func Test_New(t *testing.T) {
	if _, err := New(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_SameShared(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	b, err := New()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	ak := a.CalcSharedSecret(b.Intermediate())
	bk := b.CalcSharedSecret(a.Intermediate())

	if ak.Cmp(bk) != 0 {
		t.Errorf("two shared keys are different: \n%v\n%v", ak.String(), bk.String())
	}
}
