package main

import "context"

// App exposes Go bindings to the Wails frontend.
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ProductVersion is used by the frontend until login (T3).
func (a *App) ProductVersion() string {
	return "0.1.0"
}

