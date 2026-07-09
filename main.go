/*
Copyright © 2026 NAME HERE alessandro.dinato@gmail.com
*/
package main

import (
	"gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/cmd"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
