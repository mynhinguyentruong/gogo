package main

import "testing"

func TestGreetingHello(t *testing.T) {
  hello := GreetingHello()

  if hello != "Hello" {
    t.Errorf("It did not greet Hello but %s instead", hello)
  }

} 
