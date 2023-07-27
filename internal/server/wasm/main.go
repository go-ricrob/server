//go:build wasm

package main

import (
	"fmt"

	"syscall/js"
)

func main() {
	fmt.Println("Go Web Assembly Hello")

	newGame()

	done := make(chan struct{})
	<-done
}

type V js.Value

func (v V) Global() Global     { return Global{Value: js.Value(v)} }
func (v V) Window() Window     { return Window{Value: js.Value(v)} }
func (v V) Document() Document { return Document{Value: js.Value(v)} }
func (v V) Style() Style       { return Style{Value: js.Value(v)} }
func (v V) Element() Element   { return Element{Value: js.Value(v)} }
func (v V) Event() Event       { return Event{Value: js.Value(v)} }

func (v V) Body() Body           { return Body{Element: v.Element()} }
func (v V) Div() Div             { return Div{Element: v.Element()} }
func (v V) Canvas() Canvas       { return Canvas{Element: v.Element()} }
func (v V) Context2d() Context2d { return Context2d{Element: v.Element()} }

type Global struct {
	js.Value
}

func (v Global) Window() Window     { return V(v.Get("window")).Window() }
func (v Global) Document() Document { return V(v.Get("document")).Document() }

type Event struct {
	js.Value
}

type Window struct {
	js.Value
}

func (v Window) InnerWidth() int  { return v.Get("innerWidth").Int() }
func (v Window) InnerHeight() int { return v.Get("innerHeight").Int() }
func (v Window) AddEventListener(ev string, fn func(Event)) {
	v.Call("addEventListener", ev, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn(V(this).Event())
		return nil
	}))
}

type Document struct {
	js.Value
}

func (v Document) Body() Body { return V(v.Get("body")).Body() }

func (v Document) CreateDiv() Div       { return V(v.Call("createElement", "div")).Div() }
func (v Document) CreateCanvas() Canvas { return V(v.Call("createElement", "canvas")).Canvas() }

type Style struct {
	js.Value
}

func (v Style) SetPosition(value string) { v.Set("position", value) }
func (v Style) SetZIndex(value int)      { v.Set("z-index", value) }

type Element struct {
	js.Value
}

func (v Element) Style() Style                 { return V(v.Get("style")).Style() }
func (v Element) AppendChild(child js.Wrapper) { v.Call("appendChild", child.JSValue()) }

type Body struct {
	Element
}

type Div struct {
	Element
}

type Canvas struct {
	Element
}

func (v Canvas) Width() int           { return v.Get("width").Int() }
func (v Canvas) Height() int          { return v.Get("height").Int() }
func (v Canvas) SetWidth(width int)   { v.Set("width", width) }
func (v Canvas) SetHeight(height int) { v.Set("height", height) }

func (v Canvas) GetContext() Context2d { return V(v.Call("getContext", "2d")).Context2d() }

type Context2d struct {
	Element
}

func (v Context2d) FillStyle(style interface{})   { v.Set("fillStyle", style) }
func (v Context2d) StrokeStyle(style interface{}) { v.Set("strokeStyle", style) }

func (v Context2d) FillRect(x, y, width, height int) { v.Call("fillRect", x, y, width, height) }

func (v Context2d) BeginPath()      { v.Call("beginPath") }
func (v Context2d) ClosePath()      { v.Call("closePath") }
func (v Context2d) Stroke()         { v.Call("stroke") }
func (v Context2d) MoveTo(x, y int) { v.Call("moveTo", x, y) }
func (v Context2d) LineTo(x, y int) { v.Call("lineTo", x, y) }

type game struct {
	board *board
}

func newGame() *game {
	document := V(js.Global()).Global().Document()
	div := document.CreateDiv()
	document.Body().AppendChild(div)
	return &game{board: newBoard(div)}
}

const (
	numLayer = 3
)

type board struct {
	layers []Canvas
}

func newBoard(div Div) *board {
	document := V(js.Global()).Global().Document()
	b := &board{layers: make([]Canvas, numLayer)}
	for i := 0; i < numLayer; i++ {
		canvas := document.CreateCanvas()
		canvas.Style().SetPosition("absolute")
		canvas.Style().SetZIndex(i)
		div.AppendChild(canvas)
		b.layers[i] = canvas
	}

	window := V(js.Global()).Global().Window()
	window.AddEventListener("resize", b.onResize)
	b.redraw()
	return b
}

func (b *board) onResize(e Event) { println("resize"); b.redraw() }

func (b *board) redraw() {
	window := V(js.Global()).Global().Window()

	width := window.InnerWidth()
	height := window.InnerHeight()
	size := 0
	if width > height {
		size = height
	} else {
		size = width
	}

	for _, canvas := range b.layers {
		canvas.SetWidth(size)
		canvas.SetHeight(size)
	}
	size -= 50 // TODO
	size = (size / 16) * 16

	b.drawBackground(b.layers[0], size)
	b.drawGrid(b.layers[1], size)

	// drawBackground ...

	// func (m *trackMap) resize(evt *Event) {
	// 	width := Window.InnerWidth()
	// 	height := Window.InnerHeight()

	// 	for _, canvas := range m.layers {
	// 		canvas.SetWidth(width)
	// 		canvas.SetHeight(height)
	// 	}
	// 	m.redraw(width, height)
	// }

}

func (b *board) drawBackground(canvas Canvas, size int) {
	ctx := canvas.GetContext()
	ctx.FillStyle("black")
	ctx.FillRect(0, 0, size, size)
}

func (b *board) drawGrid(canvas Canvas, size int) {
	ctx := canvas.GetContext()
	ctx.BeginPath()
	for x := 0; x <= size; x += (size / 16) {
		ctx.MoveTo(x, 0)
		ctx.LineTo(x, size)
	}
	for y := 0; y <= size; y += (size / 16) {
		ctx.MoveTo(0, y)
		ctx.LineTo(size, y)
	}
	ctx.StrokeStyle("white")
	ctx.Stroke()
}

// root, ok := Dom.GetElementById("root").(*Div)
// if !ok {
// 	fmt.Printf("invalid root element type %v\n", root)
// 	return
// }

// eventSourceConstructor := js.Global().Get("EventSource")
// eventSource := eventSourceConstructor.New("/sse/")

// eventSource.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 	fmt.Println("message fired")
// 	return nil
// }))

// eventSource.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 	fmt.Println("event error")
// 	return nil
// }))

// eventSource.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 	fmt.Println("event open")
// 	return nil
// }))

// eventSource.Call("addEventListener", "ping", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 	fmt.Printf("go event data %s\n", args[0].Get("data"))
// 	//fmt.Printf("%v\n", this)
// 	return nil
// }))

// //h := newSSEHandler("/sse/")

// trackMap := newTrackMap(root)

// type sseHandler struct {
// }

// func newSSEHandler(url string) *sseHandler {
// 	// create eventsource

// 	fmt.Printf("setup sse handler %s\n", url)

// 	eventSourceConstructor := js.Global().Get("EventSource")
// 	eventSource := eventSourceConstructor.New(url)

// 	eventSource.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		panic("event")
// 		fmt.Println("event fired")
// 		// fmt.Println(this)
// 		return nil
// 	}))

// 	eventSource.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		fmt.Println("event error")
// 		return nil
// 	}))

// 	eventSource.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		fmt.Println("event open")
// 		return nil
// 	}))

// 	return &sseHandler{}
// }

// func (h *sseHandler) cleanup() {
// 	return
// }

// // const evtSource = new EventSource("ssedemo.php");

// type grid struct{}

// func (g *grid) redraw(ctx *CanvasRenderingContext2D, tileWidth, tileHeight, width, height int) {
// 	ctx.BeginPath()
// 	for x := 0; x <= width; x += tileWidth {
// 		ctx.MoveTo(x, 0)
// 		ctx.LineTo(x, height)
// 	}
// 	for y := 0; y <= height; y += tileHeight {
// 		ctx.MoveTo(0, y)
// 		ctx.LineTo(width, y)
// 	}
// 	ctx.StrokeStyle("white")
// 	ctx.Stroke()
// }

// type pos struct {
// 	i, j int
// }

// type tiles struct {
// 	//	en 38 x 63 mm. Ab der Bauform Sp Dr S60 wurden die Tischfelder auf 34 x 54 mm v

// 	iPos, jPos            int // '0' coordinate tiles
// 	tileWidth, tileHeight int

// 	m map[pos]tile
// }

// func newTiles() *tiles {
// 	ts := &tiles{tileWidth: 126, tileHeight: 76, m: map[pos]tile{}}
// 	ts.init() // create some test data
// 	return ts
// }

// func (ts *tiles) init() {
// 	ts.setTile(newTrack(0, 0))
// 	ts.setTile(newTrack(5, 5))
// 	ts.setTile(newTrack(6, 5))
// 	ts.setTile(newTrack(10, 10))
// }

// func (ts *tiles) setTile(t tile) {
// 	i, j := t.coord()
// 	ts.m[pos{i, j}] = t
// }

// func (ts *tiles) visibleTiles(width, height int) (i0, j0, i1, j1 int) {
// 	iMax := width / ts.tileWidth
// 	if width%ts.tileWidth != 0 {
// 		iMax++
// 	}
// 	jMax := height / ts.tileHeight
// 	if height%ts.tileHeight != 0 {
// 		jMax++
// 	}
// 	return ts.iPos, ts.jPos, ts.iPos + iMax, ts.jPos + jMax
// }

// func (ts *tiles) redraw(ctx *CanvasRenderingContext2D, width, height int) {
// 	i0, j0, i1, j1 := ts.visibleTiles(width, height)

// 	for _, t := range ts.m {
// 		i, j := t.coord()
// 		if i >= i0 && i <= i1 && j >= j0 && j <= j1 { // check if visible

// 			// ctx.SetTransform(1, 0, 0, 1, 0, 0)
// 			// ctx.Translate(i*ts.tileWidth, j*ts.tileHeight) // TODO iPos, jPos

// 			ctx.SetTransform(1, 0, 0, 1, (i-ts.iPos)*ts.tileWidth, (j-ts.jPos)*ts.tileHeight)

// 			t.draw(ctx, ts.tileWidth, ts.tileHeight)
// 		}
// 	}
// }

// type tile interface {
// 	coord() (int, int)
// 	draw(ctx *CanvasRenderingContext2D, width, height int)
// }

// type track struct {
// 	i, j int
// }

// func newTrack(i, j int) *track {
// 	return &track{i: i, j: j}
// }

// func (t *track) coord() (int, int) {
// 	return t.i, t.j
// }

// func (t *track) draw(ctx *CanvasRenderingContext2D, width, height int) {
// 	ctx.LineWidth(6)
// 	ctx.BeginPath()
// 	ctx.MoveTo(0, height/2)
// 	ctx.LineTo(width, height/2)
// 	ctx.StrokeStyle("yellow")
// 	ctx.Stroke()

// 	// ctx.FillStyle("red")
// 	// ctx.FillRect(0, 0, width, height)
// }

// type trackMap struct {
// 	layers []*Canvas

// 	grid  *grid
// 	tiles *tiles
// }

// func newTrackMap(div *Div) *trackMap {
// 	m := &trackMap{tiles: newTiles(), grid: new(grid)}
// 	m.init(div)
// 	return m
// }

// func (m *trackMap) cleanup() {} // TODO

// func (m *trackMap) init(div *Div) {
// 	m.layers = make([]*Canvas, 3)
// 	for i := range m.layers {
// 		canvas := NewCanvas()
// 		canvas.Style().SetPosition(VAbsolute)
// 		canvas.Style().SetZIndex(i)
// 		div.AppendChild(canvas)
// 		m.layers[i] = canvas
// 	}
// 	m.resize(nil)
// 	Window.AddEventListener(EvtResize, m.resize)
// }

// func (m *trackMap) resize(evt *Event) {
// 	width := Window.InnerWidth()
// 	height := Window.InnerHeight()

// 	for _, canvas := range m.layers {
// 		canvas.SetWidth(width)
// 		canvas.SetHeight(height)
// 	}
// 	m.redraw(width, height)
// }

// func (m *trackMap) redraw(width, height int) {
// 	// background
// 	ctx := m.layers[0].GetContext()
// 	ctx.FillStyle("black")
// 	ctx.FillRect(0, 0, width, height)

// 	// tiles
// 	ctx = m.layers[1].GetContext()
// 	m.tiles.redraw(ctx, width, height)

// 	// grid
// 	ctx = m.layers[2].GetContext()
// 	m.grid.redraw(ctx, m.tiles.tileWidth, m.tiles.tileHeight, width, height)

// }

// // Window.requestAnimationFrame()
