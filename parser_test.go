package main

import (
	"fmt"
	"testing"
)

func TestParsing(t *testing.T) {
    testImage := NewTestImage(2, 2)
    testImage.Set(0,0, colToColor[LightBlue])
    testImage.Set(1, 0, colToColor[MediumBlue])
    testImage.Set(1, 1, colToColor[LightBlue])
    testImage.Set(0, 1, colToColor[Black])

    tokens:= Tokenize(testImage)

    pg := ParseTokens(tokens)

    _, ok := pg.GetEdge(0, DpRight, CcLeft)
    if !ok {
        t.Errorf("No edge found for starting position")
        return
    }
    for i := 0; i < len(pg.adjList); i++ {
        fmt.Printf("Edges (%d)\n", i)
        for j := 0; j < len(pg.adjList[i]); j++ {
            edge := pg.adjList[i][j]
            fmt.Printf("  Edge %s %s\n", edge.Dp, edge.Cc)
            fmt.Printf("    Target: %d\n", edge.Target)
            if edge.Op == Push {
                fmt.Printf("    Push %d\n", edge.Data)
            } else {
                fmt.Printf("    %s\n", edge.Op)
            }
        }
    }
}
