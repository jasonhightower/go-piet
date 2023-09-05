package main

import (
    "testing"
    "image"
    "image/color"
)

var colToColor = map[Col]color.Color {
    Black: color.RGBA{A:high},
    White: color.RGBA{A:high, R: high, G: high, B: high},

    LightRed: color.RGBA{A:high, R:high, G:mid, B:mid},
    MediumRed: color.RGBA{A:high, R:high},
    DarkRed: color.RGBA{A:high, R:mid},

    LightGreen: color.RGBA{A:high, G:high, B:mid, R:mid},
    MediumGreen: color.RGBA{A:high, G:high},
    DarkGreen: color.RGBA{A:high, G:mid},

    LightBlue: color.RGBA{A:high, B:high, R:mid, G:mid},
    MediumBlue: color.RGBA{A:high, B:high},
    DarkBlue: color.RGBA{A:high, B:mid},

    LightMagenta: color.RGBA{A:high, R:high, G:mid, B:high},
    MediumMagenta: color.RGBA{A:high, R:high, B:high},
    DarkMagenta: color.RGBA{A:high, R:mid, B:mid},

    LightYellow: color.RGBA{A:high, R:high, G:high, B:mid},
    MediumYellow: color.RGBA{A:high, R:high, G:high},
    DarkYellow: color.RGBA{A:high, R:mid, G:mid},

    LightCyan: color.RGBA{A:high, R:mid, G:high, B:high},
    MediumCyan: color.RGBA{A:high, G:high, B:high},
    DarkCyan: color.RGBA{A:high, G:mid, B:mid},
}

type TestImage struct {
    bounds image.Rectangle
    data []color.Color
}
func NewTestImage(width int, height int) TestImage {
    return TestImage{
        bounds: image.Rectangle{
                    Min: image.Point{X:0, Y:0},
                    Max: image.Point{X:int(width), Y:int(height)},
        },
        data: make([]color.Color, width * height),
    }
}

func (t TestImage) ColorModel() color.Model {
    return color.RGBAModel
}
func (t TestImage) Bounds() image.Rectangle {
    return t.bounds
}
func (t TestImage) At(x int, y int) color.Color {
    return t.data[y * t.bounds.Max.X + x]    
}
func (t *TestImage) Set(x int, y int, clr color.Color) {
    t.data[y * t.bounds.Max.X + x] = clr
}

func (t *TestImage) SetRect(rect image.Rectangle, clr color.Color) {
    for x := rect.Min.X; x < rect.Max.X; x++ {
        for y := rect.Min.Y; y < rect.Max.Y; y++ {
            t.Set(x, y, clr)
        }
    }
}


func TestTokenize(t *testing.T) {
    tImg := NewTestImage(10, 1)

    tImg.SetRect(
        image.Rectangle{image.Point{X:0,Y:0}, image.Point{X:5, Y:1}}, 
        colToColor[LightBlue])
    tImg.SetRect(
        image.Rectangle{image.Point{X:5,Y:0}, image.Point{X:10, Y:1}}, 
        colToColor[MediumBlue])

    for x := 0; x < 5; x++ {
        tImg.Set(x, 0, colToColor[LightBlue])
    }

    var img image.Image = tImg
    
    tokens := Tokenize(img)
    if tokens.At(0, 0) == nil {
        t.Errorf("Error starting token was nil")
        return 
    }
    if tokens.At(0, 0).Size != 5 {
        t.Errorf("Expected 0,0 to have shape of size 5, got %d", tokens.At(0, 0).Size)
    }
    for i := 1; i < 5; i++ {
        if tokens.At(0, 0) != tokens.At(i, 0) {
            t.Errorf("Shape was not filled correctly at %d, %d", i, 0)
        }
    }
    if tokens.At(0, 0) == tokens.At(5, 0) {
        t.Errorf("Shape should have been different at %d, %d", 5, 0)
    }
    for i := 6; i < len(tokens.data); i++ {
        if tokens.At(5, 0) != tokens.At(i, 0) {
            t.Errorf("Shape was not filled correctly at %d, %d", i, 0)
        }
    }

    // create a bunch of shapes

    // tokenize

    // verify all shapes are defined correctly
}
