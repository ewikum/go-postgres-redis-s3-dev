package main

import "testing"
func TestHello(t *testing.T){
	expected:="Hello2"
	actual:=hello()
	if actual!=expected{
		t.Errorf("Test failed",expected,actual)
	}
}