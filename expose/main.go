package main

import (
    "log"
    "os/exec"
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
}

func main() {
    log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)

    sdl.Main(run)
    log.Printf("Bye")
}
