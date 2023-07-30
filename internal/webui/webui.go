//go:build wasm

// Package webui implements the ricrob web UI functionality as web assembly.
package webui

import (
	"syscall/js"
)

var global = newGlobal(js.Global())

type _global struct {
	js.Value
	window   *window
	document *document
}

func newGlobal(value js.Value) *_global {
	return &_global{
		Value:    value,
		window:   newWindow(value.Get("window")),
		document: newDocument(value.Get("document")),
	}
}

type event struct {
	js.Value
}

func newEvent(value js.Value) *event { return &event{Value: value} }

type window struct {
	js.Value
}

func newWindow(value js.Value) *window { return &window{Value: value} }

func (v *window) innerWidth() int  { return v.Get("innerWidth").Int() }
func (v *window) innerHeight() int { return v.Get("innerHeight").Int() }
func (v *window) addEventListener(ev string, fn func(*event)) {
	v.Call("addEventListener", ev, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn(newEvent(this))
		return nil
	}))
}

type document struct {
	js.Value
}

func newDocument(value js.Value) *document { return &document{Value: value} }

func (v *document) body() *body { return newBody(newElement(v.Value.Get("body"))) }

func (v *document) createDiv() *div { return newDiv(newElement(v.Call("createElement", "div"))) }
func (v *document) createCanvas() *canvas {
	return newCanvas(newElement(v.Call("createElement", "canvas")))
}

type style struct {
	js.Value
}

func newStyle(value js.Value) *style { return &style{Value: value} }

func (v *style) setPosition(value string) { v.Set("position", value) }
func (v *style) setZIndex(value int)      { v.Set("z-index", value) }

type element struct {
	js.Value
}

func newElement(value js.Value) *element { return &element{Value: value} }

func (v *element) style() *style              { return newStyle(v.Get("style")) }
func (v *element) appendChild(child js.Value) { v.Call("appendChild", child) }

type body struct {
	*element
}

func newBody(element *element) *body { return &body{element: element} }

type div struct {
	*element
}

func newDiv(element *element) *div { return &div{element: element} }

type canvas struct {
	*element
}

func newCanvas(element *element) *canvas { return &canvas{element: element} }

// func (v canvas) width() int           { return v.Get("width").Int() }
// func (v canvas) height() int          { return v.Get("height").Int() }
func (v *canvas) setWidth(width int)   { v.Set("width", width) }
func (v *canvas) setHeight(height int) { v.Set("height", height) }

func (v *canvas) getContext() *context2d { return newContext2d(newElement(v.Call("getContext", "2d"))) }

type context2d struct {
	*element
}

func newContext2d(element *element) *context2d { return &context2d{element: element} }

func (v *context2d) fillStyle(style interface{})   { v.Set("fillStyle", style) }
func (v *context2d) strokeStyle(style interface{}) { v.Set("strokeStyle", style) }

func (v *context2d) fillRect(x, y, width, height int) { v.Call("fillRect", x, y, width, height) }

func (v *context2d) beginPath()      { v.Call("beginPath") }
func (v *context2d) closePath()      { v.Call("closePath") }
func (v *context2d) stroke()         { v.Call("stroke") }
func (v *context2d) moveTo(x, y int) { v.Call("moveTo", x, y) }
func (v *context2d) lineTo(x, y int) { v.Call("lineTo", x, y) }

// WebUI represents a ricrob web ui.
type WebUI struct {
	board *board
}

// New returns a new WebUI instance.
func New() *WebUI {
	div := global.document.createDiv()
	global.document.body().appendChild(div.Value)
	return &WebUI{board: newBoard(div)}
}

const (
	numLayer = 3
)

type board struct {
	window *window
	layers []*canvas
}

func newBoard(div *div) *board {
	b := &board{window: global.window, layers: make([]*canvas, numLayer)}
	for i := 0; i < numLayer; i++ {
		canvas := global.document.createCanvas()
		canvas.style().setPosition("absolute")
		canvas.style().setZIndex(i)
		div.appendChild(canvas.Value)
		b.layers[i] = canvas
	}

	b.window.addEventListener("resize", b.onResize)
	b.redraw()
	return b
}

func (b *board) onResize(e *event) { println("resize"); b.redraw() }

func (b *board) redraw() {
	width := b.window.innerWidth()
	height := b.window.innerHeight()
	size := 0
	if width > height {
		size = height
	} else {
		size = width
	}

	for _, canvas := range b.layers {
		canvas.setWidth(size)
		canvas.setHeight(size)
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

func (b *board) drawBackground(canvas *canvas, size int) {
	ctx := canvas.getContext()
	ctx.fillStyle("black")
	ctx.fillRect(0, 0, size, size)
}

func (b *board) drawGrid(canvas *canvas, size int) {
	ctx := canvas.getContext()
	ctx.beginPath()
	defer ctx.closePath()
	for x := 0; x <= size; x += (size / 16) {
		ctx.moveTo(x, 0)
		ctx.lineTo(x, size)
	}
	for y := 0; y <= size; y += (size / 16) {
		ctx.moveTo(0, y)
		ctx.lineTo(size, y)
	}
	ctx.strokeStyle("white")
	ctx.stroke()
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
