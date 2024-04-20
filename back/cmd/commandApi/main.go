package main

import "github.com/enchik0reo/commandApi/internal/app"

// @title Script Executor API
// @version 1.0
// @description API Server for Script Executor
// @host localhost:8008
// @BasePath /
func main() {
	app.New().MustRun()
}
