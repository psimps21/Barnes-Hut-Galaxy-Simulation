package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestComputeCentrOfGravity(t *testing.T) {
	t1Star1 := &Star{mass: 1.0, position: OrderedPair{-4.0, 7.0}}
	t1Star2 := &Star{mass: 1.0, position: OrderedPair{0.0, 5.0}}
	t1Star3 := &Star{mass: 1.0, position: OrderedPair{10.0, 3.0}}
	test1 := []*Star{t1Star1, t1Star2, t1Star3}

	t2Star1 := &Star{mass: 4.0, position: OrderedPair{-4.0, 7.0}}
	t2Star2 := &Star{mass: 2.0, position: OrderedPair{0.0, 5.0}}
	t2Star3 := &Star{mass: 1.0, position: OrderedPair{10.0, 3.0}}
	test2 := []*Star{t2Star1, t2Star2, t2Star3}

	v1 := ComputeCenterOfGravity(test1)
	v2 := ComputeCenterOfGravity(test2)

	if v1.x != 2.0 && v1.y != 5.0 {
		v1x := fmt.Sprintf("%f", v1.x)
		v1y := fmt.Sprintf("%f", v1.y)
		t.Errorf("ComputeCenterOfGravity(test1) = %s", "("+v1x+","+v1y+")")
	}

	if v2.x != -6.0/7.0 && v2.y != 41.0/7.0 {
		v2x := fmt.Sprintf("%f", v2.x)
		v2y := fmt.Sprintf("%f", v2.y)
		t.Errorf("ComputeCenterOfGravity(test1) = %s", "("+v2x+","+v2y+")")
	}
}

func TestFindLeaves(t *testing.T) {
	var t1n1 *Node
	t1n2 := &Node{parent: t1n1}
	t1n3 := &Node{parent: t1n1}
	t1n1.children = []*Node{t1n2, t1n3}

	v1 := t1n1.FindLeaves([]*Node{})

	if reflect.DeepEqual(v1, []*Node{t1n2, t1n3}) {
		t.Error("Failed. ComputeCenterOfGravity(test1)")

	}

}
