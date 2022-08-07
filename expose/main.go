package main

import (
    "log"
    "os"
    "os/signal"
    "os/exec"
    "math"
    "math/rand"
    "time"
    "context"
    "github.com/veandco/go-sdl2/sdl"
)

func HasGlxinfo() bool {
    glxinfo_path, err := exec.LookPath("glxinfo")
    if err != nil {
        return true
    }
    glxinfo := exec.Command(glxinfo_path)
    err = glxinfo.Run()
    return err == nil
}

type Window struct {
    OriginalX int
    OriginalY int
    OriginalWidth int
    OriginalHeight int

    X float64
    Y float64
    Width int
    Height int
    Color sdl.Color
}

/* return true if the point (x,y) is inside the rectangle given be the upper left coordinate (x1, x2)
 * and lower right coordinate (x2, y2)
 */
func insideRect(x int, y int, x1 int, y1 int, x2 int, y2 int) bool {
    if x < x1 {
        return false
    }
    if x > x2 {
        return false
    }
    if y < y1 {
        return false
    }
    if y > y2 {
        return false
    }

    return true
}

/* return true if any of the 4 points in this window is inside the rect given by the other window */
func (window Window) Overlaps(other Window, margin int) bool {
    x1 := int(window.X) - margin
    y1 := int(window.Y) - margin

    x2 := int(window.X) + window.Width + margin
    y2 := int(window.Y) - margin

    x3 := int(window.X) - margin
    y3 := int(window.Y) + window.Height + margin

    x4 := int(window.X) + window.Width + margin
    y4 := int(window.Y) + window.Height + margin

    ox1 := int(other.X) - margin
    oy1 := int(other.Y) - margin
    ox4 := int(other.X) + other.Width + margin
    oy4 := int(other.Y) + other.Height + margin

    return insideRect(x1, y1, ox1, oy1, ox4, oy4) ||
           insideRect(x2, y2, ox1, oy1, ox4, oy4) ||
           insideRect(x3, y3, ox1, oy1, ox4, oy4) ||
           insideRect(x4, y4, ox1, oy1, ox4, oy4)
}

func createWindow(x int, y int, width int, height int, color sdl.Color) Window {
    return Window{
        OriginalX: x,
        OriginalY: y,
        OriginalWidth: width,
        OriginalHeight: height,
        X: float64(x),
        Y: float64(y),
        Width: width,
        Height: height,
        Color: color,
    }
}

func isOverlapping(windows []Window) bool {
    for i := 0; i < len(windows); i++ {
        for j := i+1; j < len(windows); j++ {
            margin := 10
            if windows[i].Overlaps(windows[j], margin) || windows[j].Overlaps(windows[i], margin) {
                return true
            }
        }
    }

    return false
}

func doMinimize(windows []Window){
    margin := 10
    for i := 0; i < len(windows); i++ {

        overlapping := false
        fx := 0.0
        fy := 0.0

        for j := 0; j < len(windows); j++ {
            if i == j {
                continue
            }

            if windows[i].Overlaps(windows[j], 10) || windows[j].Overlaps(windows[i], margin) {
                overlapping = true
                cx1 := windows[i].X + float64(windows[i].Width) / 2
                cy1 := windows[i].Y + float64(windows[i].Height) / 2

                cx2 := windows[j].X + float64(windows[j].Width) / 2
                cy2 := windows[j].Y + float64(windows[j].Height) / 2

                radians := math.Atan2(cy1 - cy2, cx1 - cx2)
                // log.Printf("Window %v pushed %v\n", i, radians * 180 / math.Pi)

                fx += math.Cos(radians)
                fy += math.Sin(radians)
            }
        }

        if overlapping {
            if windows[i].Width > 20 {
                windows[i].Width -= 1
            }
            if windows[i].Height > 20 {
                windows[i].Height -= 1
            }
        }

        windows[i].X += fx
        windows[i].Y += fy

        if windows[i].X < 0 {
            windows[i].X = 0
        }
        if windows[i].Y < 0 {
            windows[i].Y = 0
        }

        if windows[i].X + float64(windows[i].Width) > 1000 {
            windows[i].X = float64(1000 - windows[i].Width)
        }
        if windows[i].Y + float64(windows[i].Height) > 1000 {
            windows[i].Y = float64(1000 - windows[i].Height)
        }
    }
}

func doMaximize(windows []Window){
    epsilon := 0.1
    for i := 0; i < len(windows); i++ {
        if math.Abs(windows[i].X - float64(windows[i].OriginalX)) > epsilon {
            if math.Abs(windows[i].X - float64(windows[i].OriginalX)) <= 1 {
                windows[i].X = float64(windows[i].OriginalX)
            } else if windows[i].X > float64(windows[i].OriginalX) {
                windows[i].X -= 1
            } else if windows[i].X < float64(windows[i].OriginalX) {
                windows[i].X += 1
            }
        }

        if math.Abs(windows[i].Y - float64(windows[i].OriginalY)) > epsilon {
            if math.Abs(windows[i].Y - float64(windows[i].OriginalY)) <= 1 {
                windows[i].Y = float64(windows[i].OriginalY)
            } else if windows[i].Y > float64(windows[i].OriginalY) {
                windows[i].Y -= 1
            } else if windows[i].Y < float64(windows[i].OriginalY) {
                windows[i].Y += 1
            }
        }

        if windows[i].Width < windows[i].OriginalWidth {
            windows[i].Width += 1
        }
        if windows[i].Height < windows[i].OriginalHeight {
            windows[i].Height += 1
        }
    }
}

func randomInt(max int) int {
    if max <= 0 {
        return 0
    }
    return rand.Intn(max)
}

func randomWindows(max int) []Window {
    var out []Window

    for i := 0; i < max; i++ {
        x := randomInt(1000)
        y := randomInt(1000)
        w := 200 + randomInt(1000 - x - 200)
        h := 200 + randomInt(1000 - y - 200)
        out = append(out, createWindow(x, y, w, h, sdl.Color{
            R: uint8(randomInt(255)),
            G: uint8(randomInt(255)),
            B: uint8(randomInt(255)),
            A: 255,
        }))
    }

    return out
}

func run(){
    var err error

    if !HasGlxinfo() {
        sdl.Do(func(){
            sdl.SetHint(sdl.HINT_RENDER_DRIVER, "software")
        })
    }

    sdl.Do(func(){
        log.Printf("Init sdl")
        err = sdl.Init(sdl.INIT_EVERYTHING)
    })

    if err != nil {
        log.Printf("Could not initialize sdl: %v", err)
        return
    }

    defer sdl.Do(func(){
        sdl.Quit()
    })

    var window *sdl.Window
    var renderer *sdl.Renderer

    sdl.Do(func(){
        log.Printf("Creating window")
        window, renderer, err = sdl.CreateWindowAndRenderer(1000, 1000, sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE)

        if window != nil {
            window.SetTitle("Expose demo")
        }
    })

    if err != nil {
        log.Printf("Could not create window and renderer: %v", err)
        return
    }

    defer sdl.Do(func(){
        window.Destroy()
    })
    defer sdl.Do(func(){
        renderer.Destroy()
    })

    /*
    windows := []Window{
        createWindow(100, 100, 300, 300, sdl.Color{R: 255, G: 0, B: 0, A: 255}),
        createWindow(300, 200, 200, 300, sdl.Color{R: 0, G: 255, B: 0, A: 255}),
    }
    */
    windows := randomWindows(5)

    quit, cancel := context.WithCancel(context.Background())
    signals := make(chan os.Signal, 2)
    go func(){
        <-signals
        cancel()
    }()

    signal.Notify(signals, os.Interrupt)

    timer := time.NewTicker(time.Second / 30)
    defer timer.Stop()

    render := func(windows []Window){
        renderer.SetDrawColor(0, 0, 0, 0)
        renderer.Clear()

        for _, window := range windows {
            renderer.SetDrawColor(window.Color.R, window.Color.G, window.Color.B, window.Color.A)
            renderer.FillRect(&sdl.Rect{
                X: int32(window.X),
                Y: int32(window.Y),
                W: int32(window.Width),
                H: int32(window.Height),
            })
        }

        renderer.Present()
    }

    state := ""

    handleEvents := func(){
        event := sdl.WaitEventTimeout(1)
        if event != nil {
            switch event.GetType() {
                case sdl.QUIT: cancel()
                case sdl.KEYDOWN:
                    key := event.(*sdl.KeyboardEvent)
                    switch key.Keysym.Sym {
                        case sdl.K_ESCAPE: cancel()
                        case sdl.K_MINUS:
                            log.Printf("Minimize")
                            state = "minimize"
                        case sdl.K_EQUALS:
                            log.Printf("Maximize")
                            state = "maximize"
                    }
            }
        }
    }

    move := time.NewTicker(time.Second / 30)
    defer move.Stop()

    speed := 8
    for quit.Err() == nil {
        select {
            case <-move.C:
                switch state {
                    case "": break
                    case "minimize":
                        for i := 0; i < speed; i++ {
                            doMinimize(windows)
                        }
                    case "maximize":
                        for i := 0; i < speed; i++ {
                            doMaximize(windows)
                        }
                }
            case <-timer.C:
                sdl.Do(func(){
                    render(windows)
                })
            case <-quit.Done():
                break
            default:
                sdl.Do(func(){
                    handleEvents()
                })
        }
    }
}

func main() {
    log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)
    rand.Seed(time.Now().UnixNano())

    log.Printf("Press - to apply expose, and = to make the windows their original size")

    sdl.Main(run)
    log.Printf("Bye")
}
