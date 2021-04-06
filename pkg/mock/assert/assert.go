package assert

import (
	"errors"
	"reflect"
	"testing"
)

type Shell struct{}

func Equal(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Expected %#v but got %#v\n", a, b)
	}
}

func Bool(t *testing.T, a *bool, b *bool) {
	if a == nil {
		Equal(t, true, b == nil)
	} else {
		Equal(t, *a, *b)
	}
}

func True(t *testing.T, a bool) {
	if a != true {
		t.Errorf("Expected %#v to be true\n", a)
	}
}

func Error(t *testing.T, err error, target error) {
	if !errors.Is(err, target) {
		t.Errorf("Expected error %#v to be %#v\n", err, target)
	}
}
