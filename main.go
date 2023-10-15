package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
//    "embed"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
    "bytes"
    "text/template"
)

// templates/macho64/main.tmpl

var (
    //mainTmpl embed.FS
    asmTemplate *template.Template
)

func init() {
    tmpl, err := template.New("main.tmpl").Funcs(template.FuncMap{
          "IsBlock": func(stmt Stmt) bool {
              _, ok := stmt.(StmtBlock) 
              return ok
          },
          "IsCall": func(stmt Stmt) bool {
              _, ok := stmt.(Call)
              return ok
          },
          "IsOp": func(stmt Stmt, op string) bool {
              if _, ok := stmt.(Call); ok {
                  return (stmt.(Call)).Op.String() == op
              }
              return false
          },
          "HasArgs": func(stmt Stmt) bool {
              if _, ok := stmt.(Call); ok {
                  return (stmt.(Call)).Op == Push
              }
              return false
          },
      }).ParseFiles("templates/macho64/main.tmpl")//, "templates/macho64/stmt.tmpl")
      if err != nil {
          panic(err)
      }
      asmTemplate = tmpl
  //    baseLayout := template.Must(template.New("layout").ParseFS(mainTmpl, templateLayout))
}

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
func (p *PietTokens) Width() int {
    return len(p.data)
}
func (p *PietTokens) Height() int {
    return len(p.data[0])
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

type Node struct {
    x, y int
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
    min := s.Len() - int(depth)
    s.Reverse(min, s.Len())

    rolls = rolls % depth
    mid := min + int(rolls)
    s.Reverse(min, mid)
    s.Reverse(mid, s.Len())
}

func (s *Stack) Reverse(from int, to int) {
    to -= 1
    count := (to - from )/ 2
    for i := 0; i < count; i++ {
        s.Swap(from + i, to - i)
    }
}

func (s *Stack) Swap(source int, target int) {
    tmp := s.data[target]
    s.data[target] = s.data[source]
    s.data[source] = tmp
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

func ParseStmt(tokens *PietTokens, capacity int) Stmt {
    dp := DpRight
    cc := CcLeft

    carrot := Carrot{X:0, Y:0, tokens: tokens}

    root := StmtBlock{}
    stack := NewStack(capacity)

    curShape := carrot.CurrentShape()

    attempts := 8
    for true {
        ok := carrot.Move(dp, cc)
        if !ok {
            attempts -= 1
            if attempts == 0 {
                root.Append(Call{Op:Exit})
                return root
            }
            root.Append(Call{Op: Push, Args: []int32 {1}})
            if attempts % 2 > 0 {
                root.Append(Call{Op: Switch})
                cc = cc.Toggle()
            } else {
                root.Append(Call{Op: Pointer})
                dp = dp.Rotate(1)
            }
            continue
        }
        nextShape := carrot.CurrentShape()
        if nextShape.Color == White {
            root.Append(Call{Op: Push, Args: []int32 {1}})
            root.Append(Call{Op: Dup})
            root.Append(Call{Op: Switch})
            root.Append(Call{Op: Pointer})
            attempts -= 2
            curShape = nextShape
            continue
        }

//        fmt.Printf("Move to (%d, %d) - %s %s : [%d]\n", carrot.X, carrot.Y, dp, cc, curShape.Size)
        op := curShape.Color.ToOp(nextShape.Color)
//        fmt.Println(op)
        attempts = 8
        switch op {
        case Switch: 
            root.Append(Call{Op:Pointer})
            if val, ok := stack.Pop(); ok {
               if val % 2 > 0 {
                   cc = cc.Toggle()
               }
            }
        case Pointer:
            root.Append(Call{Op:Switch})
            if val, ok := stack.Pop(); ok {
                dp = dp.Rotate(val)
            }
        case Push: 
            stack.Push(curShape.Size)
            root.Append(Call{Op: op, Args: []int32 {curShape.Size}})
        case Add: 
            if f, s, ok := stack.Pop2(); ok {
               stack.Push(s + f)
            }
            root.Append(Call{Op: op})
        case Sub: 
            if f, s, ok := stack.Pop2(); ok {
                stack.Push(s - f)
            }
            root.Append(Call{Op: op})
        case Mult:
            if f, s, ok := stack.Pop2(); ok {
                stack.Push(s * f)
            }
            root.Append(Call{Op: op})
        case Div:
            if f, s, ok := stack.Pop2(); ok {
                stack.Push(s / f)
            }
            root.Append(Call{Op: op})
        case NumOut:
            stack.Pop()
            root.Append(Call{Op: op})
        case CharOut:
            stack.Pop()
            root.Append(Call{Op: op})
        case NumIn:
            panic("NumIn not supported")
        case CharIn:
            panic("NumIn not supported")
        case Roll:
            if f, s, ok := stack.Pop2(); ok {
                stack.Roll(s, f)
            }
            root.Append(Call{Op: op})
        case Pop:
            stack.Pop()
            root.Append(Call{Op: op})
        case Dup:
            if val, ok := stack.Peek(); ok {
                stack.Push(val)
            }
            root.Append(Call{Op: op})
        case Greater:
            if f, s, ok := stack.Pop2(); ok {
                if s > f {
                    stack.Push(1)
                } else {
                    stack.Push(0)
                }
            }
            root.Append(Call{Op: op})
        case Mod:
            if f, s, ok := stack.Pop2(); ok {
                stack.Push(s % f)
            }
            root.Append(Call{Op: op})
        case Noop:
        default:
            panic(fmt.Sprintf("Unhandled operator %s", op))
        }
        curShape = nextShape
    }

    return root
}

type Carrot struct {
    X, Y int
    tokens *PietTokens
}
func (c *Carrot) CurrentShape() *Shape {
    return c.tokens.At(c.X, c.Y)
}
func (c *Carrot) slideWhite(dp Dp, cc Cc) bool {
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
    curShape := c.tokens.At(c.X, c.Y)
    width, height := c.tokens.Width(), c.tokens.Height()
    for InBounds(c.X + xAdj, c.Y + yAdj, width, height) {
        c.X += xAdj
        c.Y += yAdj
        if c.tokens.At(c.X, c.Y) != curShape {
            if curShape.Color == Black {
                c.X -= xAdj
                c.Y -= yAdj
                return false
            }
            return true
        }
    }
    return false
}
func (c *Carrot) Move(dp Dp, cc Cc) bool {
    xPos, yPos := c.X, c.Y
    curShape := c.tokens.At(c.X, c.Y)
    if curShape.Color == White {
        return c.slideWhite(dp, cc)
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
    if InBounds(xPos, yPos, c.tokens.Width(), c.tokens.Height()) {
        if c.tokens.At(xPos, yPos).Color == Black {
            return false
        }
        c.X = xPos
        c.Y = yPos
        return true
    }
    return false
}

func main() {
    filename := flag.String("f", "", "name of the piet file to interpret")
    codelsize := flag.Int("codel-size", 1, "Size of codels to support enlarged images for better viewing")
    capacity := flag.Int("capacity", 512, "Capacity of the stack")
    mode := flag.String("m", "run", "(run | compile)")
    help := flag.Bool("h", false, "Print Help/Usage")
    flag.Parse()

    if *help == true {
        flag.Usage()
        os.Exit(0)
    }
    if *mode != "run" && *mode != "compile" {
        fmt.Printf("Unrecogznied mode %s, expected one of (run, compile)\n", *mode)
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

    if *mode == "compile" {
        segments := strings.Split(*filename, "/")
        name := strings.Split(segments[len(segments) - 1], ".")[0]

        tmplName := fmt.Sprintf("%s.asm", name)
        tmplF, err := os.Create(tmplName)
        if err != nil {
            panic(err)
        }
        defer tmplF.Close()
        CompileTmpl(stmt, tmplF)

        // nasm -fmacho64 tetris.asm
        cmd := exec.Command("nasm", "-fmacho64", tmplName)
        out, err := cmd.Output()
        if err != nil {
            // if there was any error, print it here
            panic(err)
        } else {
            cmd = exec.Command("xcode-select", "-p")
            xcodePath, err := cmd.Output()
            if err != nil {
                panic(err)
            }
            xcodePathStr := strings.TrimSpace(string(xcodePath))
            objFile := fmt.Sprintf("%s.o", name)
    // ld -e _main  -macosx_version_min 10.10 -arch x86_64 -lSystem -L$(xcode-select -p)/SDKs/MacOSX.sdk/usr/lib -o tetris tetris.o
            cmd = exec.Command("ld", "-v", 
                                      "-e", 
                                      "_main", 
                                      "-macosx_version_min",
                                      "10.10", 
                                      "-arch",
                                      "x86_64",
                                      "-lSystem",
                                      fmt.Sprintf("-L%s/SDKs/MacOSX.sdk/usr/lib", xcodePathStr),
                                      "-o",
                                      name,
                                      objFile)

            var out bytes.Buffer
            var stderr bytes.Buffer
            cmd.Stdout = &out
            cmd.Stderr = &stderr
            err = cmd.Run()
            if err != nil {
                fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
                return
            }
            fmt.Println("Result: " + out.String())
                                      
// ld -e _main  -macosx_version_min 10.10 -arch x86_64 -lSystem -L$(xcode-select -p)/SDKs/MacOSX.sdk/usr/lib -o tetris tetris.o*
        }

        fmt.Println("Output: ", string(out))
    } else {
        interpreter := NewInterpreter(*capacity)
        interpreter.Interpret(stmt)
    }
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
        if assign.Name == "dp" {
            interpreter.Dp = Dp(assign.val)
        } else if assign.Name == "cc" {
            interpreter.Cc = Cc(assign.val)
        }
    } else if block, ok := stmt.(StmtBlock); ok {
        if block.Children != nil {
            for _, s := range block.Children {
                interpreter.Interpret(s)
            }
        }
    } else if call, ok := stmt.(Call); ok {
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
            case Greater:
                if f, s, ok := interpreter.Stack.Pop2(); ok {
                    if s > f {
                        interpreter.Stack.Push(1)
                    } else {
                        interpreter.Stack.Push(0)
                    }
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
            case NumOut:
                if val, ok := interpreter.Stack.Pop(); ok {
                    fmt.Print(val)
                }
            case CharOut:
                if val, ok := interpreter.Stack.Pop(); ok {
                    fmt.Print(string(val))
                }
            case NumIn:
                b := make([]byte, 1)
                os.Stdin.Read(b)
                val, err := strconv.Atoi(string(b[0]))
                if err != nil {
                    panic(err)
                }
                interpreter.Stack.Push(int32(val))
            case CharIn:
                b := make([]byte, 1)
                os.Stdin.Read(b)
                interpreter.Stack.Push(int32(b[0]))
            case Roll:
                if f, s, ok := interpreter.Stack.Pop2(); ok {
                    interpreter.Stack.Roll(s, f)
                }
            case Exit:
                fmt.Println()
                os.Exit(0)
            default:
                panic(fmt.Sprintf("%s not supported", call.Op))
        }

    }
}

func CompileTmpl(stmt Stmt, f io.Writer) {
    err := asmTemplate.Execute(f, stmt)
    if err != nil {
        panic(err)
    }
}
