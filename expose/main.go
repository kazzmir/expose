package main

import (
    "log"
    "os"
    "os/signal"
    "os/exec"
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
    X int
    Y int
    Width int
    Height int
    Color sdl.Color
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

    windows := []Window{
        Window{
            X: 100,
            Y: 100,
            Width: 100,
            Height: 100,
            Color: sdl.Color{R: 255, G: 0, B: 0, A: 255},
        },
        Window{
            X: 500,
            Y: 200,
            Width: 100,
            Height: 100,
            Color: sdl.Color{R: 0, G: 255, B: 0, A: 255},
        },
    }

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
    }

    for quit.Err() == nil {
        select {
            case <-timer.C:
                render(windows)
            case <-quit.Done():
                break
        }
    }
}

func main() {
    log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)

    sdl.Main(run)
    log.Printf("Bye")
}
