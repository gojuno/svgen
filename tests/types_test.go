package tests

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestStringType_Scan(t *testing.T) {
	v := StringType("")

	testImplements(t, &v)

	if err := v.Scan("1"); err != nil {
		t.Fatal(err)
	}
	if v != StringValue1 {
		t.Fatalf("invalid value: %v", v)
	}

	if err := v.Scan([]byte("2")); err != nil {
		t.Fatal(err)
	}
	if v != StringValue2 {
		t.Fatalf("invalid value: %v", v)
	}

	if err := v.Scan(nil); err != nil {
		t.Fatal(err)
	}

	if err := v.Scan("undefined"); err == nil {
		t.Fatal("Scan must return error because of invalid value")
	}

	if err := v.Scan(false); err == nil {
		t.Fatal("Scan must return error because of invalid type")
	}
}

func TestStringType_Value(t *testing.T) {
	v := StringType("")

	testImplements(t, &v)

	var value driver.Value
	var err error

	value, err = StringValue1.Value()
	if err != nil {
		t.Fatal(err)
	}
	if value != "1" {
		t.Fatalf("invalid value: %v", value)
	}

	value, err = StringType("3").Value()
	if err == nil {
		t.Fatal("Value must return error because of invalid value")
	}
	if value != nil {
		t.Fatalf("Value must return nil: %v", value)
	}

	value, err = StringType("").Value()
	if err != nil {
		t.Fatal(err)
	}
	if value != nil {
		t.Fatalf("Value must return nil: %v", value)
	}
}

func TestIntType_Scan(t *testing.T) {
	v := IntType(0)

	testImplements(t, &v)

	if err := v.Scan(int64(1)); err != nil {
		t.Fatal(err)
	}
	if v != IntValue1 {
		t.Fatalf("invalid value: %v", v)
	}

	if err := v.Scan(nil); err != nil {
		t.Fatal(err)
	}

	if err := v.Scan(int64(3)); err == nil {
		t.Fatalf("Scan must return error because of invalid value")
	}

	if err := v.Scan(false); err == nil {
		t.Fatalf("Scan must return error because of invalid type")
	}
}

func TestIntType_Value(t *testing.T) {
	v := IntType(0)

	testImplements(t, &v)

	var value driver.Value
	var err error

	value, err = IntValue1.Value()
	if err != nil {
		t.Fatal(err)
	}
	if value != int64(1) {
		t.Fatalf("invalid value: %v", value)
	}

	value, err = IntType(3).Value()
	if err == nil {
		t.Fatal("Value must return error because of invalid value")
	}
	if value != nil {
		t.Fatalf("Value must return nil: %v", value)
	}
}

func TestFloatType(t *testing.T) {
	v := FloatType(0.0)

	testNotImplements(t, &v)
}

func TestString2Type(t *testing.T) {
	v := String2Type("")

	testNotImplements(t, &v)
}

func testImplements(t *testing.T, v interface{}) {
	if !checkScannerImplementation(v) {
		t.Fatalf("%T doesn't implement sql.Scanner interface", v)
	}
	if !checkValuerImplementation(v) {
		t.Fatalf("%T doesn't implement driver.Valuer interface", v)
	}
}

func testNotImplements(t *testing.T, v interface{}) {
	if checkScannerImplementation(v) {
		t.Fatalf("%T implements sql.Scanner interface", v)
	}
	if checkValuerImplementation(v) {
		t.Fatalf("%T implements driver.Valuer interface", v)
	}
}

func checkScannerImplementation(v interface{}) bool {
	return reflect.TypeOf(v).Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem())
}

func checkValuerImplementation(v interface{}) bool {
	return reflect.TypeOf(v).Implements(reflect.TypeOf((*driver.Valuer)(nil)).Elem())
}
