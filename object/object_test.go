package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is Johnny"}
	diff2 := &String{Value: "My name is Johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	int1 := &Integer{Value: 1}
	int2 := &Integer{Value: 1}
	diff1 := &Integer{Value: 2}
	diff2 := &Integer{Value: 2}

	if int1.HashKey() != int2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}

	if int1.HashKey() == diff1.HashKey() {
		t.Errorf("integers with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	boolTrue1 := &Boolean{Value: true}
	boolTrue2 := &Boolean{Value: true}
	boolFalse1 := &Boolean{Value: false}
	boolFalse2 := &Boolean{Value: false}

	if boolTrue1.HashKey() != boolTrue2.HashKey() {
		t.Errorf("booleans with same content have different hash keys")
	}

	if boolFalse1.HashKey() != boolFalse2.HashKey() {
		t.Errorf("booleans with same content have different hash keys")
	}

	if boolTrue1.HashKey() == boolFalse1.HashKey() {
		t.Errorf("booleans with different content have same hash keys")
	}
}
