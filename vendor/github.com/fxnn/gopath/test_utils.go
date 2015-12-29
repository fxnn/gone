package gopath

import (
	"reflect"
	"testing"
)

func assertGoPathHasErr(t *testing.T, actual GoPath) {
	if !actual.HasErr() {
		t.Errorf("Expected to have Err(), but was %v", actual)
	}
}

func assertGoPathEqual(t *testing.T, actual GoPath, expected GoPath) {
	if actual.Err() != expected.Err() {
		t.Errorf("expected Err() to be %v, but was %v", expected.Err(), actual.Err())
	}
	if actual.Path() != expected.Path() {
		t.Errorf("expected Path() to be %v, but was %v", expected.Path(), actual.Path())
	}
}

func assertStrEqual(t *testing.T, actual string, expected string) {
	if expected != actual {
		t.Errorf("expected to be equal %s, but was %s", expected, actual)
	}
}

func assertSliceEmpty(t *testing.T, actual interface{}) {
	switch reflect.TypeOf(actual).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(actual)
		if slice.Len() != 0 {
			t.Errorf("expected to be empty, but was %v", actual)
		}
	}
}

func assertStrSliceEqual(t *testing.T, actual []string, expected ...string) {
	if len(actual) != len(expected) {
		t.Errorf("expected to be a [] string of len %d, but was of len %d", len(expected), len(actual))
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("expected %dth element to be %v, but was %v", expected[i], actual[i])
		}
	}
}
