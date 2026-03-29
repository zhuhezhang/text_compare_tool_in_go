package main

import (
	"context"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// App struct
type App struct {
	ctx context.Context
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

// CompareText compares two texts and returns the diff in HTML format along with the number of lines in the result
func (a *App) CompareText(text1, text2 string) string {
	lines1 := strings.Split(text1, "\n")
	lines2 := strings.Split(text2, "\n")

	dmp := diffmatchpatch.New()
	var diff []diffmatchpatch.Diff
	var result = ""
	i, j := 0, 0

	for i < len(lines1) || j < len(lines2) {
		if i < len(lines1) && j < len(lines2){
			if lines1[i] == lines2[j] && lines1[i] == "" {
				result = result + "<br>"
			}else{
				diff = dmp.DiffMain(lines1[i], lines2[j], false)
				result = result + "<p>" + dmp.DiffPrettyHtml(diff) + "</p>"
			}
            i++
            j++
		}else if i < len(lines1) {
			if lines1[i] != "" {
			    diff = dmp.DiffMain(lines1[i], "", false)
			    result = result + "<p>" + dmp.DiffPrettyHtml(diff) + "</p>"
		    }else{
			    result = result + "<br>"
		    }
			i++
		}else if j < len(lines2) {
			if lines2[j] != "" {
			    diff = dmp.DiffMain("", lines2[j], false)
			    result = result + "<p>" + dmp.DiffPrettyHtml(diff) + "</p>"
		    }else{
			    result = result + "<br>"
		    }
			j++
		}
	}

    return result
}