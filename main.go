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
    if p.data == nil || p.data[x] == nil {
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
    // TODO JH need a better data structure than this
    leftPoints []image.Point
    rightPoints []image.Point
    topPoints []image.Point
    bottomPoints []image.Point
}
func (s Shape) String() string {
    result := "{ right: ["
    for i := 0; i < len(s.rightPoints); i++ {
        if i > 0 {
            result += ", "
        }
        result += fmt.Sprintf("%s", s.rightPoints[i])
    }
    result += "], bottom: ["
    for i := 0; i < len(s.bottomPoints); i++ {
        if i > 0 {
            result += ", "
        }
        result += fmt.Sprintf("%s", s.bottomPoints[i])
    }
    result += fmt.Sprint("], left: [")
    for i := 0; i < len(s.leftPoints); i++ {
        if i > 0 {
            result += ", "
        }
        result += fmt.Sprintf("%s", s.leftPoints[i])
    }
    result += fmt.Sprint("], top: [")
    for i := 0; i < len(s.topPoints); i++ {
        if i > 0 {
            result += ", "
        }
        result += fmt.Sprintf("%s", s.topPoints[i])
    }
    result += fmt.Sprint("] }")
    return result
}
func (s *Shape) AddPoint(p image.Point) {
    s.Size += 1
    if s.Size == 1 {
        s.leftPoints = append(s.leftPoints, p)
        s.rightPoints = append(s.rightPoints, p)
        s.topPoints = append(s.topPoints, p)
        s.bottomPoints = append(s.bottomPoints, p)
    } else {
        if s.leftPoints[0].X == p.X {
            s.leftPoints = append(s.leftPoints, p)
        } else if s.leftPoints[0].X > p.X {
            s.leftPoints = make([]image.Point, 1)
            s.leftPoints[0] = p
        }
        if s.bottomPoints[0].Y == p.Y {
            s.bottomPoints = append(s.bottomPoints, p)
        } else if s.bottomPoints[0].Y < p.Y {
            s.bottomPoints = make([]image.Point, 1)
            s.bottomPoints[0] = p
        }
        if s.rightPoints[0].X == p.X {
            s.rightPoints = append(s.rightPoints, p)
        } else if s.rightPoints[0].X < p.X {
            s.rightPoints = make([]image.Point, 1)
            s.rightPoints[0] = p
        }
        if s.topPoints[0].Y == p.Y {
            s.topPoints = append(s.topPoints, p)
        } else if s.topPoints[0].Y > p.Y {
            s.topPoints = make([]image.Point, 1)
            s.topPoints[0] = p
        }
    }
}
func (s *Shape) FindPoint(pos image.Point, dp Dp, cc Cc) image.Point {
   switch dp {
   case DpLeft:
       if len(s.leftPoints) == 1 {
           return s.leftPoints[0]
       }
       panic("Need to implement")
   case DpRight:
       if len(s.rightPoints) == 1 {
           return s.rightPoints[0]
       }
       panic("Need to implement")
   case DpDown:
        if len(s.bottomPoints) == 1 {
            return s.bottomPoints[0]
        }
        panic("Need to implement")
   case DpUp:
        if len(s.bottomPoints) == 1 {
            return s.bottomPoints[0]
        }
        panic("Need to implement")
   }
   panic("Unknown direction")
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

func FindEdge(idx int, dp Dp, cc Cc, tokens *PietTokens, pg *PietGraph) {
    shape := tokens.shapes[idx]
    if shape.Color == White {
        panic("White not supported")
    }
    var edge image.Point
    switch dp {
    case DpRight:
        edge = shape.rightPoints[0]
        if len(shape.rightPoints) > 1{
            if cc == CcRight {
                for i := 1; i < len(shape.rightPoints); i++ {
                    if edge.Y < shape.rightPoints[i].Y {
                        edge = shape.rightPoints[i]
                    }
                }
            } else {
                for i := 1; i < len(shape.rightPoints); i++ {
                    if edge.Y > shape.rightPoints[i].Y {
                        edge = shape.rightPoints[i]
                    }
                }
            }
        }
    case DpDown:
        edge = shape.bottomPoints[0]
        if len(shape.bottomPoints) > 1 {
            if cc == CcRight {
                for i := 1; i < len(shape.bottomPoints); i++ {
                    if edge.X > shape.bottomPoints[i].X {
                        edge = shape.bottomPoints[i]
                    }
                }
            } else {
                for i := 1; i > len(shape.bottomPoints); i++ {
                    if edge.X < shape.bottomPoints[i].X {
                        edge = shape.bottomPoints[i]
                    }
                }
            }
        }
    case DpLeft:
        edge = shape.leftPoints[0]
        if len(shape.leftPoints) > 1 {
            if cc == CcRight {
                for i := 1; i < len(shape.leftPoints); i++ {
                    if edge.Y > shape.leftPoints[i].Y {
                        edge = shape.leftPoints[i]
                    }
                }
            } else {
                for i := 1; i < len(shape.leftPoints); i++ {
                    if edge.Y < shape.leftPoints[i].Y {
                        edge = shape.leftPoints[i]
                    }
                }
            }
        }
    case DpUp:
        edge = shape.topPoints[0]
        if len(shape.topPoints) > 1 {
            if cc == CcRight {
                for i:= 1; i < len(shape.topPoints); i++ {
                    if edge.X < shape.topPoints[i].X {
                        edge = shape.topPoints[i]
                    } 
                }
            } else {
                for i:= 1; i < len(shape.topPoints); i++ {
                    if edge.X > shape.topPoints[i].X {
                        edge = shape.topPoints[i]
                    } 
                }
            }
        }
    default:
        panic("Unkown dp direction")
    }

    var adjEdge image.Point
    switch dp {
    case DpRight:
        adjEdge = image.Point{X: edge.X + 1, Y: edge.Y}
    case DpDown:
        adjEdge = image.Point{X: edge.X, Y: edge.Y + 1}
    case DpLeft:
        adjEdge = image.Point{X: edge.X - 1, Y: edge.Y}
    case DpUp:
        adjEdge = image.Point{X: edge.X, Y: edge.Y - 1}
    default:
        panic("Unknown direction")
    }

    if !adjEdge.In(tokens.Bounds) {
        return
    }
    targetShape := tokens.At(adjEdge.X, adjEdge.Y)
    if targetShape.Color == Black {
        return
    } 

    index := tokens.data[adjEdge.X][adjEdge.Y]
    // create Edge

    newEdge := Edge{
        Op: shape.Color.ToOp(targetShape.Color),
        Data: shape.Size,
        Dp: dp,
        Cc: cc,
        Target: index,
    }
    pg.adjList[idx] = append(pg.adjList[idx], newEdge)
}

type Edge struct {
    Op Op
    Data int32
    Dp Dp
    Cc Cc
    Target int
}

type Instr struct {
    Op Op
    Data uint32
}
type Block struct {
    Label string
    Instructions []Instr
}
func (b *Block) Add(instr Instr) {
    b.Instructions = append(b.Instructions, instr)
}

type Executable struct {
    blocks map[string]*Block
}
func NewExecutable() *Executable {
    return &Executable{blocks: make(map[string]*Block)}
}
func (e *Executable) Get(label string) (*Block, bool) {
    val, ok := e.blocks[label]
    return val, ok
}
func (e *Executable) Add(block *Block) {
    e.blocks[block.Label] = block
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

/*
type Parser struct {
    cc Cc
    dp Dp
    pos image.Point
    stack Stack
    tokens *PietTokens
}
func NewParser(tokens *PietTokens, capacity int) *Parser {
    return &Parser{
        tokens: tokens,
        stack: NewStack(capacity),
    }
}
func (p *Parser) Parse() *Executable {
    exec := NewExecutable()
    p.parseBlock(exec)
    return exec
}

func (p *Parser) parseBlock(exec *Executable) {
    label := p.curLabel()
    if _, contains := exec.Get(label); !contains {
        block := Block{Label: label}
        exec.Add(&block)
        
        // parse instructions
        for true {
            curShape := p.tokens[p.pos.X][p.pos.Y]
            
            if !p.moveToNextShape() {
                var dpRot uint32 = 0
                ccTog := 0
                rotateCc := true
                valid := false
                // TODO JH this needs to be 7 once I have better test data
                for i := 0; i < 3; i++ {
                    // find next move in a loop here and create a block
                    if rotateCc {
                        fmt.Println("Rotate CC")
                        ccTog += 1
                        p.cc = p.cc.Toggle(1)
                    } else {
                        fmt.Println("Rotate DP")
                        dpRot += 1
                        p.dp = p.dp.Rotate(1)
                    }
                    rotateCc = !rotateCc
                    if p.moveToNextShape() {
                        valid = true
                        break
                    }
                }
                if !valid {
                    block.Add(Instr{Op:Exit})
                    return
                }
                if ccTog % 2 > 0 {
                    block.Add(Instr{Op:Switch})
                }
                if dpRot > 0 {
                    block.Add(Instr{Op:Pointer, Data: dpRot})
                }
            }
            nextShape := p.tokens[p.pos.X][p.pos.Y]
            op := curShape.Color.ToOp(nextShape.Color)
            instr := Instr{Op: op}
            if op == Push {
                instr.Data = curShape.Size
            }
            block.Add(instr)
        }
    }
}
*/

/*
func (p *Parser) moveToNextShape() bool {
    var adj = image.Point{}
    switch p.dp {
    case DpLeft:
        adj.X = -1
    case DpRight:
        adj.X = 1
    case DpUp:
        adj.Y = -1
    case DpDown:
        adj.Y = 1
    default:
        panic("Unknown dp")
    }
    cur := p.pos
    startShape := p.tokens[cur.X][cur.Y]
    for p.tokens[cur.X][cur.Y] == startShape {
        cur = cur.Add(adj)
        if cur.X < 0 || cur.Y < 0 || cur.X >= len(p.tokens) || cur.Y >= len(p.tokens[cur.X]) {
            return false
        }
        if p.tokens[cur.X][cur.Y].Color == Black {
            fmt.Printf("%d, %d - is Black\n", cur.X, cur.Y)
            return false
        }
    }
    p.pos = cur
    fmt.Printf("Move to %d, %d\n", p.pos.X, p.pos.Y)
    return true
}

func (p *Parser) curLabel() string {
    return Label(p.pos.X, p.pos.Y, p.dp, p.cc)
}
*/
func Label(x int, y int, dp Dp, cc Cc) string {
    return fmt.Sprintf("block_%d_%d__%d_%d", x, y, dp, cc)
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

    fmt.Printf("VISITED %d\n", len(visited))

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
    /*
    for i := 0; i < tokens.Size(); i++ {
        fmt.Printf("%d - %s\n", i, *tokens.shapes[i])
    }
    fmt.Println()
    */
    pg := Parse(tokens)
    Prune(pg, *capacity)
    /*
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
    */
    Interpret(pg, *capacity)

    // compile ... 
}

