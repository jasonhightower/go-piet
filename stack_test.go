package main

import "testing"

func TestEmptyPop(t *testing.T) {
    s := NewStack(3)

    _, ok := s.Pop()
    if ok {
        t.Errorf("Pop on empty stack should have returned false")
    }
}

func TestPop2NotEnoughItems(t *testing.T) {
    s := NewStack(3)

    _, _, ok := s.Pop2()
    if ok {
        t.Errorf("Pop2 on empty stack should have returned false")
    }

    s.Push(1)
    _, _, ok = s.Pop2()
    if ok {
        t.Errorf("Pop2 on single item stack should have returned false")
    }

}

func TestPop2(t *testing.T) {
    s := NewStack(3)
    s.Push(1)
    s.Push(2)

    first, second, ok := s.Pop2()
    if !ok {
        t.Errorf("Pop2 on two item stack returned false")
    }
    if first != 2 {
        t.Errorf("Expected %d for first item. Got %d", 2, first)
    }
    if second != 1 {
        t.Errorf("Expected %d for second item. Got %d", 1, second)
    }
}

func TestPushAndPop (t *testing.T) {
    const v1 = 10
    const v2 = 20
    s := NewStack(3) 
    s.Push(v1)
    s.Push(v2)

    got, ok := s.Pop()
    if !ok {
        t.Errorf("Pop was unsuccessful")
    } else if got != v2 {
        t.Errorf("Got %d expected %d", got, v2)
    }

    got, ok = s.Pop()
    if !ok {
        t.Errorf("Pop was unsuccessful")
    } else if got != v1 {
        t.Errorf("Got %d expected %d", got, v1)
    }

}

