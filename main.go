package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
    "flag"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type CodelImage struct {
    img image.Image
    csize int
    bounds image.Rectangle
}
func NewCodelImage(img image.Image, codelSize int) *CodelImage {
    rect := image.Rectangle{
        Min: img.Bounds().Min.Div(codelSize),
        Max: img.Bounds().Max.Div(codelSize),
    }
    return &CodelImage{
        img: img,
        csize: codelSize,
        bounds: rect,
    }
}
func (c CodelImage) At(x int, y int) color.Color {
   return c.img.At(x * c.csize, y * c.csize ) 
}
func (c CodelImage) Bounds() image.Rectangle {
    return c.bounds;
}
func (c CodelImage) ColorModel() color.Model {
    return c.img.ColorModel()
}

type Dp byte
const (
    DpRight Dp = 0
    DpDown Dp = 1
    DpLeft Dp = 2
    DpUp Dp = 3
)
func (d Dp) String() string {
    switch d {
    case DpRight:
        return "right"
    case DpLeft:
        return "left"
    case DpUp:
        return "up"
    case DpDown:
        return "down"
    default:
        return "unknown"
    }
}
func (d Dp) Rotate(times int32) Dp {
    return Dp((int32(d) + times) % 4)
}

type Cc byte
const (
    CcLeft Cc = 0
    CcRight Cc = 1
)
func (c Cc) String() string {
    if c == CcLeft {
        return "left"
    } else if c == CcRight {
        return "right"
    } else {
        return "unknown"
    }
}
func (c Cc) Toggle() Cc {
    if c == CcLeft {
        return CcRight
    }
    return CcLeft
}

type Op byte
const (
    Push Op = 1
    Pop Op = 2
    Add Op = 3
    Sub Op = 4
    Mult Op = 5
    Div Op = 6
    Mod Op = 7
    Not Op = 8
    Greater Op = 9
    Pointer Op = 10
    Switch Op = 11
    Dup Op = 12
    Roll Op = 13
    NumIn Op = 14
    CharIn Op = 15
    NumOut Op = 16
    CharOut Op = 17
    Goto Op = 18
    Noop = 19
    Exit Op = 20
)
func (o Op) String() string {
    switch o {
    case Pop:
        return "pop"
    case Push:
        return "push"
    case Noop:
        return "noop"
    case Exit: 
        return "exit"
    case Add:
        return "add"
    case Sub:
        return "sub"
    case Mult:
        return "mult"
    case Div:
        return "div"
    case Mod:
        return "mod"
    case Not:
        return "not"
    case Greater:
        return "greater"
    case Pointer:
        return "pointer"
    case Switch:
        return "switch"
    case Dup:
        return "dup"
    case Roll:
        return "roll"
    case NumIn:
        return "num_in"
    case CharIn:
        return "char_in"
    case CharOut:
        return "char_out"
    case NumOut:
        return "num_out"
    default:
        return fmt.Sprintf("Unknown operator %d", byte(o))
    }
}

type Col byte 
const (
    LightRed Col = 0
    MediumRed Col = 1
    DarkRed Col = 2
    LightYellow Col = 3
    MediumYellow Col = 4
    DarkYellow Col = 5
    LightGreen Col = 6
    MediumGreen Col = 7 
    DarkGreen Col = 8
    LightCyan Col = 9
    MediumCyan Col = 10
    DarkCyan Col = 11
    LightBlue Col = 12
    MediumBlue Col = 13
    DarkBlue Col = 14
    LightMagenta Col = 15
    MediumMagenta Col = 16
    DarkMagenta Col = 17
    White Col = 18
    Black Col = 19
    Unrecoganized Col = 20
)
func (c Col) ToOp(o Col) Op {
    if c == Black || o == Black || c == White || o == White {
        return Noop
    }
    c_hue := int(c) / 3
    o_hue := int(o) / 3
    hue_diff := DiffAndWrap(c_hue, o_hue, 6)

    c_light := int(c) % 3
    o_light := int(o) % 3
    light_diff := DiffAndWrap(c_light, o_light, 3)

    return Op(hue_diff * 3 + light_diff)
}

func DiffAndWrap(f int, s int, max int) int {
    if s >= f {
        return s - f
    }
    return s + max - f
}

func ColorToCol(c color.Color) Col {
    if col, ok := colorToCol[c]; ok {
        return col
    }
    return Unrecoganized
}

const zero uint8 = 0x00
const mid uint8 = 0xC0
const high uint8 = 0xFF
var colorToCol = map[color.Color]Col {
    color.RGBA{A:high}: Black,
    color.RGBA{A:high, R: high, G: high, B: high}: White,

    color.RGBA{A:high, R:high, G:mid, B:mid}: LightRed,
    color.RGBA{A:high, R:high}:                 MediumRed,
    color.RGBA{A:high, R:mid}:                  DarkRed,

    color.RGBA{A:high, G:high, B:mid, R:mid}: LightGreen,
    color.RGBA{A:high, G:high}:                 MediumGreen,
    color.RGBA{A:high, G:mid}:                  DarkGreen,

    color.RGBA{A:high, B:high, R:mid, G:mid}: LightBlue,
    color.RGBA{A:high, B:high}:                 MediumBlue,
    color.RGBA{A:high, B:mid}:                  DarkBlue,

    color.RGBA{A:high, R:high, G:mid, B:high}:LightMagenta,
    color.RGBA{A:high, R:high, B:high}:        MediumMagenta,
    color.RGBA{A:high, R:mid, B:mid}:          DarkMagenta,

    color.RGBA{A:high, R:high, G:high, B:mid}:LightYellow,
    color.RGBA{A:high, R:high, G:high}:        MediumYellow,
    color.RGBA{A:high, R:mid, G:mid}:          DarkYellow,

    color.RGBA{A:high, R:mid, G:high, B:high}:LightCyan,
    color.RGBA{A:high, G:high, B:high}:        MediumCyan,
    color.RGBA{A:high, G:mid, B:mid}:          DarkCyan,
}

type PietTokens struct {
    Bounds image.Rectangle
    data [][]int
    shapes []*Shape
}
func NewPietTokens(width int, height int) *PietTokens {
    pietTokens := PietTokens{
        Bounds: image.Rectangle{
            Max: image.Point{X:width, Y:height},
        },
        data: make([][] int, width),
        shapes: []*Shape{},
    }
    for x := 0; x < width; x++ {
        pietTokens.data[x] = make([]int, height)
        for y := 0; y < height; y++ {
            pietTokens.data[x][y] = -1
        }
    }
    return &pietTokens
}
func (p *PietTokens) Size() int {
    return len(p.shapes)
}
func (p *PietTokens) At(x int, y int) *Shape {
    if p.data == nil || x < 0 || y < 0 || x >= len(p.data) || p.data[x] == nil || y >= len(p.data[x]) {
        return nil
    }
    if p.data[x][y] == -1 {
        return nil
    }
    return p.shapes[p.data[x][y]]
}
func (p *PietTokens) Add(s *Shape) {
    p.shapes = append(p.shapes, s)
}
func Tokenize(img image.Image) *PietTokens {
    pietTokens := NewPietTokens(img.Bounds().Max.X, img.Bounds().Max.Y)

    pos := image.Point{X:0,Y:0}
    for pos.X = 0; pos.X < img.Bounds().Max.X; pos.X++ {
        for pos.Y = 0; pos.Y < img.Bounds().Max.Y; pos.Y++ {
            if pietTokens.data[pos.X][pos.Y] == -1 {
                shape := Shape{
                    Color: ColorToCol(img.At(pos.X, pos.Y)),
                }
                shape.AddPoint(pos)
                idx := len(pietTokens.shapes)
                pietTokens.data[pos.X][pos.Y] = idx
                pietTokens.Add(&shape)
                fillLeft(img, pos, &shape, pietTokens)
                fillRight(img, pos, &shape, pietTokens)
                fillUp(img, pos, &shape, pietTokens)
                fillDown(img, pos, &shape, pietTokens)
            }         
        }
    }

    return pietTokens
}

// TODO JH optimize this
func fillLeft(img image.Image, pos image.Point, shape *Shape, pietTokens *PietTokens) {
    pos.X -= 1
    if !pos.In(img.Bounds()) {
        return
    }
    if pietTokens.data[pos.X][pos.Y] != -1 {
        return
    }
    colAtPos := ColorToCol(img.At(pos.X, pos.Y))
    if shape.Color != colAtPos {
        return
    }                
    pietTokens.data[pos.X][pos.Y] = len(pietTokens.shapes) - 1
    shape.AddPoint(pos)
    fillLeft(img, pos, shape, pietTokens)
    fillUp(img, pos, shape, pietTokens)
    fillDown(img, pos, shape, pietTokens)
}
func fillUp(img image.Image, pos image.Point, shape *Shape, pietTokens *PietTokens) {
    pos.Y -= 1
    if !pos.In(img.Bounds()) {
        return
    }
    if pietTokens.data[pos.X][pos.Y] != -1 {
        return
    }
    colAtPos := ColorToCol(img.At(pos.X, pos.Y))
    if shape.Color != colAtPos {
        return
    }
    pietTokens.data[pos.X][pos.Y] = len(pietTokens.shapes) - 1
    shape.AddPoint(pos)
    fillUp(img, pos, shape, pietTokens)
    fillLeft(img, pos, shape, pietTokens)
    fillRight(img, pos, shape, pietTokens)
}
func fillRight(img image.Image, pos image.Point, shape *Shape, pietTokens *PietTokens) {
    pos.X += 1
    if !pos.In(img.Bounds()) {
        return
    }
    if pietTokens.data[pos.X][pos.Y] != -1 {
        return
    }
    colAtPos := ColorToCol(img.At(pos.X, pos.Y))
    if shape.Color != colAtPos {
        return
    }
    pietTokens.data[pos.X][pos.Y] = len(pietTokens.shapes) - 1
    shape.AddPoint(pos)
    fillRight(img, pos, shape, pietTokens)
    fillUp(img, pos, shape, pietTokens)
    fillDown(img, pos, shape, pietTokens)
}
func fillDown(img image.Image, pos image.Point, shape *Shape, pietTokens *PietTokens) {
    pos.Y += 1
    if !pos.In(img.Bounds()) {
        return
    }
    if pietTokens.data[pos.X][pos.Y] != -1 {
        return
    }
    colAtPos := ColorToCol(img.At(pos.X, pos.Y))
    if shape.Color != colAtPos {
        return
    }
    pietTokens.data[pos.X][pos.Y] = len(pietTokens.shapes) - 1
    shape.AddPoint(pos)
    fillDown(img, pos, shape, pietTokens)
    fillLeft(img, pos, shape, pietTokens)
    fillRight(img, pos, shape, pietTokens)
}

type Shape struct {
    Color Col
    AdjList []Edge
    Size int32
    xEdges *TreeNode
    yEdges *TreeNode
}
func (s *Shape) AddPoint(p image.Point) {
    s.Size += 1
    if s.Size == 1 {
        s.xEdges = NewTreeNode(p.X, p.Y)
        s.yEdges = NewTreeNode(p.Y, p.X)
    } else {
        s.xEdges.Add(p.X, p.Y)
        s.yEdges.Add(p.Y, p.X)
    }
}
func (s *Shape) AddEdge(e Edge) {
    s.AdjList = append(s.AdjList, e)
}
type PietGraph struct {
    adjList [][]Edge
}
func (pg *PietGraph) GetEdge(idx int, dp Dp, cc Cc) (*Edge, bool) {
    if pg.adjList[idx] == nil {
        return nil, false
    }
    for i := 0; i < len(pg.adjList[idx]); i++ {
        edge := &pg.adjList[idx][i]
        if edge.Dp == dp && edge.Cc == cc {
            return edge, true
        }
    }
    return nil, false
}
func (pg *PietGraph) IsTerminal(idx int) bool {
    return len(pg.adjList[idx]) == 0
}

func Parse(tokens *PietTokens) *PietGraph {
    pg := &PietGraph {
        adjList: make([][]Edge, tokens.Size()),
    }
    for i := 0; i < len(tokens.shapes); i++ {
        if tokens.shapes[i].Color != Black {
            FindEdges(i, tokens, pg)
        }
    }
    return pg
}

type Node struct {
    x, y int
}



func FindEdges(idx int, tokens *PietTokens, pg *PietGraph) {
    if pg.adjList[idx] != nil {
        return
    }
    edges := make([]Edge, 0)                
    pg.adjList[idx] = edges
    for i := 0; i < 4; i++ {
        for j := 0; j < 2; j++ {
            FindEdge(idx, Dp(i), Cc(j), tokens, pg)
        }
    }
}

func SlideToEdge(pt image.Point, shape *Shape, dp Dp, cc Cc, tokens *PietTokens, pg *PietGraph) {
    var xAdj, yAdj int
    switch dp {
    case DpRight:
        xAdj = 1
    case DpDown:
        yAdj = 1
    case DpLeft:
        xAdj = -1
    case DpUp:
        yAdj = -1
    }
    var target *Shape
    x, y := pt.X, pt.Y
    for true {
        target = tokens.At(x + xAdj, y + yAdj)
        if target == nil || target.Color != shape.Color {
            break        
        }
        x += xAdj
        y += yAdj
    }
    

}

func SlideThroughWhite(pt image.Point, shape *Shape, dp Dp, cc Cc, tokens *PietTokens, pg *PietGraph) (image.Point, int, bool) {
    rotations := 0
    for rotations < 4 {
        switch dp {
        case DpRight:
            xNode, _ := shape.yEdges.Get(pt.Y)
            pt.X = xNode.Max 
            target := tokens.At(pt.X + 1, pt.Y)
            if target != nil && target.Color != Black {
                return image.Point{pt.X + 1, pt.Y}, rotations, true
            } else {
                dp = dp.Rotate(1)
                cc = cc.Toggle()
                rotations += 1
            }
        case DpDown:
            yNode, _ := shape.xEdges.Get(pt.X)
            pt.Y = yNode.Max
            target := tokens.At(pt.X, pt.Y + 1)
            if target != nil && target.Color != Black {
                // create an edge based on the target
                return image.Point{pt.X, pt.Y + 1}, rotations, true
            } else {
                dp = dp.Rotate(1)
                cc = cc.Toggle()
                rotations += 1
            }
        case DpLeft:
            xNode, _ := shape.yEdges.Get(pt.Y)
            pt.X = xNode.Min
            target := tokens.At(pt.X - 1, pt.Y)
            if target != nil && target.Color != Black {
                // create an edge based on the target
                return image.Point{pt.X - 1, pt.Y}, rotations, true
            } else {
                dp = dp.Rotate(1)
                cc = cc.Toggle()
                rotations += 1
            }
        case DpUp:
            yNode, _ := shape.xEdges.Get(pt.X)
            pt.Y = yNode.Min
            target := tokens.At(pt.X, pt.Y - 1)
            if target != nil && target.Color != Black {
                return image.Point{pt.X, pt.Y - 1}, rotations, true
            } else {
                dp = dp.Rotate(1)
                cc = cc.Toggle()
                rotations += 1
            }
        }
    }
    // create an exit edge

    // TODO JH should return an exit edge
    panic("Should return an exit edge")
}

func FindEdge(idx int, dp Dp, cc Cc, tokens *PietTokens, pg *PietGraph) {
    shape := tokens.shapes[idx]
    if shape.Color == White {
        panic("White not supported")
    }

    var treeEdge image.Point
    switch dp {
    case DpRight:
        node := shape.xEdges.MaxNode()
        if cc == CcLeft {
            treeEdge = image.Point{X: node.Key, Y: node.Min}
        } else {
            treeEdge = image.Point{X: node.Key, Y: node.Max}
        }
    case DpDown:
        node := shape.yEdges.MaxNode()
        if cc == CcLeft {
            treeEdge = image.Point{X: node.Max, Y: node.Key}
        } else {
            treeEdge = image.Point{X: node.Min, Y: node.Key}
        }
    case DpLeft:
        node := shape.xEdges.MinNode()
        if cc == CcLeft {
            treeEdge = image.Point{X: node.Key, Y: node.Max}
        } else {
            treeEdge = image.Point{X: node.Key, Y: node.Min}
        }
    case DpUp:
        node := shape.yEdges.MinNode()
        if cc == CcLeft {
            treeEdge = image.Point{X: node.Min, Y: node.Key}
        } else {
            treeEdge = image.Point{X: node.Max, Y: node.Key}
        }
    }

    var adjEdge image.Point
    switch dp {
    case DpRight:
        adjEdge = image.Point{X: treeEdge.X + 1, Y: treeEdge.Y}
    case DpDown:
        adjEdge = image.Point{X: treeEdge.X, Y: treeEdge.Y + 1}
    case DpLeft:
        adjEdge = image.Point{X: treeEdge.X - 1, Y: treeEdge.Y}
    case DpUp:
        adjEdge = image.Point{X: treeEdge.X, Y: treeEdge.Y - 1}
    default:
        panic("Unknown direction")
    }

    if !adjEdge.In(tokens.Bounds) {
        return
    }
    targetShape := tokens.At(adjEdge.X, adjEdge.Y)
    if targetShape.Color == Black {
        return
    } else if targetShape.Color == White {
        if whiteEdge, rotations, ok := SlideThroughWhite(adjEdge, targetShape, dp, cc, tokens, pg); ok {
            fmt.Printf("%s - %d\n", adjEdge, rotations)
            dp = dp.Rotate(3)
            if rotations % 2 > 0 {
                cc = cc.Toggle()
            }
            adjEdge = whiteEdge
            fmt.Printf("Need to handle rotations in the edge")
        } else {
            panic("need to return exit edge")
        }
        targetShape = tokens.At(adjEdge.X, adjEdge.Y)
    }
 

    index := tokens.data[adjEdge.X][adjEdge.Y]
    // check if target is white
    
    // create Edge

    newEdge := Edge{
        Op: shape.Color.ToOp(targetShape.Color),
        Data: shape.Size,
        Dp: dp,
        Cc: cc,
        Target: index,
        Entry: adjEdge,
    }
    pg.adjList[idx] = append(pg.adjList[idx], newEdge)
}

type Edge struct {
    Op Op
    Data int32
    Dp Dp
    Cc Cc
    Target int
    Entry image.Point
}

type Instr struct {
    Op Op
    Data uint32
}
type Stack struct {
    data []int32
    head int
    capacity int
}
func NewStack(capacity int) *Stack {
    return &Stack{
        data: make([]int32, capacity),
        head: -1,
        capacity: capacity,
    }
}
func (s Stack) String() string {
    result := fmt.Sprint("[")
    for i := 0; i <= s.head; i++ {
        if i > 0 {
            result += fmt.Sprintf(", %d", s.data[i])
        } else {
            result += fmt.Sprint(s.data[i])
        }
    }
    result += "]"
    return result
}
func (s *Stack) Len() int {
    return s.head + 1
}
func (s *Stack) Roll(depth int32, rolls int32) {
    if s.Len() <= 1 || int(depth) > s.Len() {
        return
    }
    ip := s.Len() - (int(depth) % s.Len())
    i := ip
    j := ip + (int(rolls) % int(depth))
    tmpi := s.data[i]
    var tmpj int32

    for n :=0; n < int(depth);  {
        tmpj = s.data[j]
        s.data[j] = tmpi
        i++
        n++

        if n == int(depth) {
            return
        }
        tmpi = s.data[i]
        s.data[i] = tmpj
        j++
        if j > int(s.head) {
            j = ip
        }
        n++
    }
}
func (s *Stack) Push(val int32) {
    s.head += 1
    if s.head >= s.capacity {
        panic("Stack overflow")
    }
    s.data[s.head] = val
}
func (s *Stack) Pop() (int32, bool) {
    if s.head < 0 {
        return 0, false
    }
    val := s.data[s.head]
    s.head -= 1
    return val, true
}
func (s *Stack) Pop2() (int32, int32, bool) {
    if s.head < 1 {
        return 0, 0, false
    }
    val := s.data[s.head]
    val2 := s.data[s.head - 1]
    s.head -= 2
    return val, val2, true
}
func (s *Stack) Peek() (int32, bool) {
    if s.head < 0 {
        return 0, false
    }
    return s.data[s.head], true
}
func (s *Stack) Dup() bool {
    if s.head + 1 >= s.capacity {
        panic("Stack overflow")
    }
    if s.head < 0 {
        return false
    }
    s.head += 1
    s.data[s.head] = s.data[s.head -1]
    return true
}

func readImage(filename string) (image.Image, error) {
    file, err := os.Open(filename)
    if err != nil {
        // FIXME proper error handling
        return nil, err
    }
    defer file.Close()

    image, _, err := image.Decode(file)
    return image, err
}

func Interpret(pg *PietGraph, sCapacity int) {
    stack := NewStack(sCapacity)
    cur := 0
    cc := CcLeft
    dp := DpRight
    var ccToggle bool
    for !pg.IsTerminal(cur) {
        edge, found := pg.GetEdge(cur, dp, cc)
        if found {
            dp, cc = Exec(stack, edge, dp, cc)
            cur = edge.Target
            ccToggle = true
        } else {
            if ccToggle {
                cc = cc.Toggle()
            } else {
                dp = dp.Rotate(1)
            }
            ccToggle = !ccToggle
        }
    }
}

func ExecForPrune(stack *Stack, edge *Edge, dp Dp, cc Cc) (Dp, Cc) {
    switch edge.Op {
    case Push:
        stack.Push(edge.Data)
    case Pop:
        stack.Pop()
    case Mod:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s % f)
        }
    case Add:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s + f)
        }
    case Sub:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s - f)
        }
    case Div:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s / f)
        }
    case Mult:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s * f)
        }
    case Dup:
        if val, ok := stack.Peek(); ok {
           stack.Push(val) 
        }
    case NumOut:
        stack.Pop()
    case CharOut:
        stack.Pop()
    case Switch:
        if val, ok := stack.Pop(); ok {
            if val % 2 > 0 {
                cc = cc.Toggle() 
            }
        }
    case Pointer:
        if val, ok := stack.Pop(); ok {
            dp = dp.Rotate(val)
        }
    case Roll:
        if f, s, ok := stack.Pop2(); ok {
            stack.Roll(s, f)
        }
    case Greater:
        if f, s, ok := stack.Pop2(); ok {
            if s > f {
                stack.Push(1)
            } else {
                stack.Push(0)
            }
        }
    default:
        panic(fmt.Sprintf("\n%s not implemented", edge.Op))
    }
    // TODO JH print the stack
//    fmt.Printf(" - %s", curStack)
//    fmt.Println()
    return dp, cc
}

func Exec(stack *Stack, edge *Edge, dp Dp, cc Cc) (Dp, Cc) {
//    curStack := fmt.Sprintf("%s - dp:%s cc:%s", stack, dp, cc)
//    fmt.Print(edge.Op)
    switch edge.Op {
    case Push:
        stack.Push(edge.Data)
//        fmt.Printf(" %d", edge.Data)
    case Pop:
        stack.Pop()
    case Mod:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s % f)
        }
    case Add:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s + f)
        }
    case Sub:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s - f)
        }
    case Div:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s / f)
        }
    case Mult:
        if f, s, ok := stack.Pop2(); ok {
            stack.Push(s * f)
        }
    case Dup:
        if val, ok := stack.Peek(); ok {
           stack.Push(val) 
        }
    case NumOut:
        if val, ok := stack.Pop(); ok {
            fmt.Print(val)
        }
    case CharOut:
        if val, ok := stack.Pop(); ok {
            fmt.Print(string(val))
        }
    case Switch:
        if val, ok := stack.Pop(); ok {
            if val % 2 > 0 {
                cc = cc.Toggle() 
            }
        }
    case Pointer:
        if val, ok := stack.Pop(); ok {
            dp = dp.Rotate(val)
        }
    case Roll:
        if f, s, ok := stack.Pop2(); ok {
            stack.Roll(s, f)
        }
    case Greater:
        if f, s, ok := stack.Pop2(); ok {
            if s > f {
                stack.Push(1)
            } else {
                stack.Push(0)
            }
        }
    default:
        panic(fmt.Sprintf("\n%s not implemented", edge.Op))
    }
    // TODO JH print the stack
//    fmt.Printf(" - %s", curStack)
//    fmt.Println()
    return dp, cc
}

func Prune(pg *PietGraph, sCapacity int) {
    cur := 0
    dp := DpRight
    cc := CcLeft
    stack := NewStack(sCapacity)
    visited := make(map[*Edge]bool)
    var ccToggle bool
    for !pg.IsTerminal(cur) {
        edge, found := pg.GetEdge(cur, dp, cc)
        if found {
            visited[edge] = true
            dp, cc = ExecForPrune(stack, edge, dp, cc)
            cur = edge.Target
            ccToggle = true
        } else {
            if ccToggle {
                cc = cc.Toggle()
            } else {
                dp = dp.Rotate(1)
            }
            ccToggle = !ccToggle
        }
    }

    for i := 0; i < len(pg.adjList); i++ {
        prunedEdges := make([]Edge, 0)
        for j := 0; j < len(pg.adjList[i]); j++ {
            if visited[&pg.adjList[i][j]] {
               prunedEdges = append(prunedEdges, pg.adjList[i][j]) 
            }
        }
        pg.adjList[i] = prunedEdges
    }
}

func MoveWhite(x int, y int, tokens [][] *Shape, dp Dp, cc Cc, stack *Stack) (int, int) {
    xAdj, yAdj := 0, 0
    switch dp {
    case DpRight:
        xAdj = 1
    case DpDown:
        yAdj = 1
    case DpLeft:
        xAdj = -1
    case DpUp:
        xAdj = 1
    }
    curShape := tokens[x][y]
    width, height := len(tokens), len(tokens[0])
    for InBounds(x + xAdj, y + yAdj, width, height) {
        x += xAdj
        y += yAdj
        if tokens[x][y] != curShape {
            return x, y
        }
    }
    return x, y
}

func Move(x int, y int, tokens [][]*Shape, dp Dp, cc Cc, stack *Stack) (int, int) {
    xPos, yPos := x, y
    curShape := tokens[x][y]
    if curShape.Color == White {
        return MoveWhite(x, y, tokens, dp, cc, stack)
    }
    switch dp {
    case DpRight:
        rightNode := curShape.xEdges.MaxNode()
        xPos = rightNode.Key
        if cc == CcLeft {
            yPos = rightNode.Min
        } else {
            yPos = rightNode.Max
        }
        xPos += 1
    case DpDown:
        bottomNode := curShape.yEdges.MaxNode()
        yPos = bottomNode.Key
        if cc == CcLeft {
            xPos = bottomNode.Max
        } else {
            xPos = bottomNode.Min
        }
        yPos += 1
    case DpLeft:
        leftNode := curShape.xEdges.MinNode()
        xPos = leftNode.Key
        if cc == CcLeft {
            yPos = leftNode.Max
        } else {
            yPos = leftNode.Min
        }
        xPos -= 1
    case DpUp:
        topNode := curShape.yEdges.MinNode()
        yPos = topNode.Key
        if cc == CcLeft {
            xPos = topNode.Min
        } else {
            xPos = topNode.Max
        }
        yPos -= 1
    }
    return xPos, yPos
}
func InBounds(x int, y int, width int, height int) bool {
    return x >= 0 && y >= 0 && x < width && y < height
}


type TreeNode struct {
    Key int
    Min int
    Max int
    Left *TreeNode
    Right *TreeNode
}
func NewTreeNode(key int, val int) *TreeNode {
    return &TreeNode{Key: key, Min: val, Max: val}
}
func (t *TreeNode) Get(key int) (*TreeNode, bool) {
    if t.Key == key {
        return t, true
    } else if t.Key < key {
        if t.Right == nil {
            return nil, false
        }
        return t.Right.Get(key)
    } else {
        if t.Left == nil {
            return nil, false
        }
        return t.Left.Get(key)
    }
}
func (t *TreeNode) MinNode() (*TreeNode) {
    if t.Left != nil {
        return t.Left.MinNode()
    }
    return t
}
func (t *TreeNode) MaxNode() (*TreeNode) {
    if t.Right != nil {
        return t.Right.MaxNode()
    }
    return t
}
func (t *TreeNode) Add(key int, val int) {
    if t.Key == key {
        if t.Min > val {
            t.Min = val
        }
        if t.Max < val {
            t.Max = val
        }
    } else if t.Key < key {
        if t.Right == nil {
            t.Right = NewTreeNode(key, val)
        } else {
            t.Right.Add(key, val)
        }
    } else {
        if t.Left == nil {
            t.Left = NewTreeNode(key, val)
        } else {
            t.Left.Add(key, val)
        }
    }
}

func ParseStmt(tokens [][]*Shape, capacity int) Stmt {
    dp := DpRight
    cc := CcLeft
    x, y := 0, 0
    curShape := tokens[x][y]

    stack := NewStack(capacity)
    x, y = Move(x, y, tokens, dp, cc, stack)
    nextShape := tokens[x][y]

}

func main() {
    filename := flag.String("f", "", "name of the piet file to interpret")
    codelsize := flag.Int("codel-size", 1, "Size of codels to support enlarged images for better viewing")
    capacity := flag.Int("capacity", 512, "Capacity of the stack")
    help := flag.Bool("h", false, "Print Help/Usage")
    flag.Parse()

    if *help == true{
        flag.Usage()
        os.Exit(0)
    }

    img, err := readImage(*filename)
    if err != nil {
        io.WriteString(os.Stderr, fmt.Sprint(err))
        os.Exit(1)
    } else {
        if *codelsize > 1 {
            img = NewCodelImage(img, *codelsize)
        }
    }
    tokens := Tokenize(img)
    stmt := ParseStmt(tokens, *capacity)
    
    interpreter := NewInterpreter(*capacity)
    interpreter.Interpret(stmt)


    // compile ... 
}

/*
[Stmt]        | (Assign | Call | If)
[Assign]      | Name Int
[If]          | EqExpr Block Stmt?
[EqExpr]      | Name Int
[Block]       | Stmt+
[Call]        | Name Int?
[Int]         | (int32)
*/

type AstNode interface {
}
type Stmt interface {
    AstNode
}
type Assign struct {
    Name string
    val int32    
}
type StmtBlock struct {
    Children []Stmt
}
func (s *StmtBlock) Append(stmt Stmt) {
    s.Children = append(s.Children, stmt)
}
type StmtIf struct {
    Condition EqExpr
    Block StmtBlock
    Else Stmt
}
type EqExpr struct {
    Name string
    val int32
}
type Call struct {
    Op Op
    Args []int32
}

type Interpreter struct {
    Dp Dp
    Cc Cc
    Stack *Stack
}
func NewInterpreter(capacity int) *Interpreter {
    return &Interpreter{
        Stack: NewStack(capacity),
    }
}
func (interpreter *Interpreter) Interpret(stmt Stmt) {
    if stmt == nil {
        return
    }
    if assign, ok := stmt.(Assign); ok {
//        fmt.Printf("Assign %s\n", assign.Name)
        if assign.Name == "dp" {
            interpreter.Dp = Dp(assign.val)
        } else if assign.Name == "cc" {
            interpreter.Cc = Cc(assign.val)
        }
    } else if block, ok := stmt.(StmtBlock); ok {
//        fmt.Println("Block")
        if block.Children != nil {
            for _, s := range block.Children {
                interpreter.Interpret(s)
            }
        }
    } else if call, ok := stmt.(Call); ok {
//        fmt.Printf(" Call %s\n", call.Op)
        switch call.Op {
            case Push:
                interpreter.Stack.Push(call.Args[0])
            case Pop:
                interpreter.Stack.Pop()
            case Add:
                if f, s, ok := interpreter.Stack.Pop2(); ok {
                    interpreter.Stack.Push(s + f)
                }
            case Sub:
                if f, s, ok := interpreter.Stack.Pop2(); ok {
                    interpreter.Stack.Push(s - f)
                }
            case Mult:
                if f, s, ok := interpreter.Stack.Pop2(); ok {
                    interpreter.Stack.Push(s * f)
                }
            case Div:
                if f, s, ok := interpreter.Stack.Pop2(); ok {
                    interpreter.Stack.Push(s / f)
                }
            case Dup:
                if val, ok := interpreter.Stack.Peek(); ok {
                    interpreter.Stack.Push(val)
                }
            case Switch:
                if val, ok := interpreter.Stack.Pop(); ok {
                    if val % 2 > 0 {
                        interpreter.Cc = interpreter.Cc.Toggle()
                    }
                }
            case Pointer:
                if val, ok := interpreter.Stack.Pop(); ok {
                    interpreter.Dp = interpreter.Dp.Rotate(val)
                }
            case CharOut:
                if val, ok := interpreter.Stack.Pop(); ok {
                    fmt.Print(string(val))
                }
            case Roll:
                if f, s, ok := interpreter.Stack.Pop2(); ok {
                    interpreter.Stack.Roll(s, f)
                }
            case Exit:
                fmt.Println()
                os.Exit(0)
        }

    }
}
/*
func Evaluate(stmt Stmt, inspector func(s Stmt)) {
    if stmt == nil {
        return
    }
    if s, ok := stmt.(*StmtBlock); ok {
        if s.Children != nil {
            for _, c := range s.Children {
                Evaluate(c, inspector)
            }
        }
    } else if s, ok := stmt.(*StmtIf); ok {
        // Need to evaluate the expression first
        // if true
        Evaluate(s.Block, inspector)
        // else
        Evaluate(s.Else, inspector)
    } else if s, ok := stmt.(*Call); ok {
        // Call Function
    }

}
*/

func HelloWorld() {
    root := StmtBlock{}
    root.Append(Assign{Name:"dp", val: 0})
    root.Append(Assign{Name:"cc", val: 0})
    root.Append(Call{Op: Switch})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Mult})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Mult})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Mult})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Mult})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pop})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Roll})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pop})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pop})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Mult})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Mult})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Dup})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Add})
    root.Append(Call{Op: Push, Args: []int32 {1}})
    root.Append(Call{Op: Pointer})
    root.Append(Call{Op: Div})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: CharOut})
    root.Append(Call{Op: Exit})
    
    interpreter := NewInterpreter(512)
    interpreter.Interpret(root)
}


