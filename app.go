package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// App struct
type App struct {
	ctx context.Context
}

type D1Response struct {
	Meta struct {
		Duration int `json:"duration"`
	}
	Results []any
	Success bool
}

type D1Error struct {
	Error struct {
		Text string `json:"text"`
	} `json:"error"`
}

type Db struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Version   string `json:"version"`
	NumTables int    `json:"num_tables"`
	FileSize  int    `json:"file_size"`
}

type Table struct {
	Name string `json:"name"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func runCmd(command string) ([]byte, error) {

	cmd := exec.Command("sh", "-c", command)
	//os.Setenv("PATH", "/opt/homebrew/bin:"+os.Getenv("PATH"))
	cmd.Env = append(os.Environ(), "PATH=/opt/homebrew/bin:"+os.Getenv("PATH"))
	cmd.Dir = "/Users/alexeyzoz/Desktop/Projects/Shaldag/shaldag-cf-worker-nextjs-ts"

	return cmd.CombinedOutput()
}

func fmtExeCmd(db *string, command *string) string {
	return fmt.Sprintf("npx wrangler d1 execute %s --command=\"%s\" --json", *db, *command)
}

func (a *App) Execute(db string, command string) (any, error) {

	cmdFmt := fmtExeCmd(&db, &command)

	cmdRun, err := runCmd(cmdFmt)

	var d1_err D1Error

	if err != nil {
		_ = json.Unmarshal(cmdRun, &d1_err)
		fmt.Println(err)
		return nil, errors.New(d1_err.Error.Text)
	}

	var res []D1Response

	_ = json.Unmarshal(cmdRun, &res)

	fmt.Println(res)

	return res[0].Results, nil
}

// Specifig Funcs

func (a *App) GetOS() (string, error) {
	currentOs := runtime.GOOS

	fmt.Println(currentOs)
	return currentOs, nil
}

func (a *App) GetDbs() ([]Db, error) {
	cmd := "npx wrangler d1 list --json"

	out, err := runCmd(cmd)

	var d1_err D1Error

	if err != nil {
		_ = json.Unmarshal(out, &d1_err)
		return make([]Db, 0), errors.New(d1_err.Error.Text)
	}

	var dbs []Db

	err1 := json.Unmarshal(out, &dbs)

	if err1 != nil || len(dbs) < 1 {
		_ = json.Unmarshal(out, &d1_err)
		return make([]Db, 0), errors.New(d1_err.Error.Text)
	}

	return dbs, nil
}

// script := fmt.Sprintf(`tell application "Terminal"
// 		do script "%s"
// 		end tell`, command)
//
// script := `tell application "Terminal"
// 		do script "ls"
// 		end tell`
//cmd := exec.Command("osascript", "-e", script)
