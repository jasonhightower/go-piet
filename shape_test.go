package main

import (
	"image"
	"testing"
)

func TestShape(t *testing.T) {
    shape := Shape{}
    
    shape.AddPoint(image.Point{X:0, Y:0})
    shape.AddPoint(image.Point{X:1, Y: 1})
    shape.AddPoint(image.Point{X:2, Y:0})
    shape.AddPoint(image.Point{X:2, Y: 2})
    shape.AddPoint(image.Point{X:2, Y: 1})
    shape.AddPoint(image.Point{X:1, Y:0})
    shape.AddPoint(image.Point{X:0, Y: 1})
    shape.AddPoint(image.Point{X:2, Y: 3})

    minX := shape.xEdges.MinNode().Key
    maxX := shape.xEdges.MaxNode().Key
    minY := shape.yEdges.MinNode().Key
    maxY := shape.yEdges.MaxNode().Key

    
    if minX != 0 {
        t.Errorf("Got min X %d expected 0", minX)
    }
    if maxX != 2 {
        t.Errorf("Got max X %d expected 2", maxX)
    }
    if minY != 0 {
        t.Errorf("Got min Y %d expected 0", minY)
    }
    if maxY != 3 {
        t.Errorf("Got max Y %d expected 3", maxY)
    }

    minY_Left := shape.yEdges.MinNode().Min
    minY_Right := shape.yEdges.MinNode().Max

    if minY_Left != 0 {
        t.Errorf("Left value of minimum Y was %d expected 0", minY_Left)
    }
    if minY_Right != 2 {
        t.Errorf("Right value of minimum Y was %d expected 2", minY_Right)
    }

    minX_Left := shape.xEdges.MinNode().Min
    minX_Right := shape.xEdges.MinNode().Max
    
    if minX_Left != 0 {
        t.Errorf("Left value of minimum X was %d expected 0", minX_Left)
    }
    if minX_Right != 1 {
        t.Errorf("Right value of minimum X was %d expected 1", minX_Right)
    }

    maxY_Left := shape.yEdges.MaxNode().Min
    maxY_Right := shape.yEdges.MaxNode().Max

    if maxY_Left != 2 {
        t.Errorf("Left value of maximum Y was %d expected 2", maxY_Left)
    }
    if maxY_Right != 2 {
        t.Errorf("Right value of maximum Y was %d expected 2", maxY_Right)
    }

    maxX_Left := shape.xEdges.MaxNode().Min
    maxX_Right := shape.xEdges.MaxNode().Max
    
    if maxX_Left != 0 {
        t.Errorf("Left value of maximum X was %d expected 0", maxX_Left)
    }
    if maxX_Right != 3 {
        t.Errorf("Right value of maximum X was %d expected 3", maxX_Right)
    }






}
