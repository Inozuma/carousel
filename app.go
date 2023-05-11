package main

import (
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type application struct {
	db *Database

	app fyne.App
}

func newApplication(db *Database) *application {
	return &application{
		db: db,
	}
}

func (app *application) run() error {
	app.app = fyneapp.New()
	window := app.app.NewWindow("Carousel")
	window.Resize(fyne.NewSize(1024, 768))

	items, err := app.db.LoadItems()
	if err != nil {
		return err
	}

	titleList := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			titleText := canvas.NewText("title", color.White)
			episodeText := canvas.NewText("episode", color.White)
			return container.New(layout.NewGridLayout(2), titleText, episodeText)
		},
		func(item widget.ListItemID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			titleText := container.Objects[0].(*canvas.Text)
			episodeText := container.Objects[1].(*canvas.Text)

			tm := items[item]
			if tm.Title != "" {
				titleText.Text = tm.Title
			} else {
				titleText.Text = filepath.Base(tm.Path)
			}
			episodeText.Text = strconv.Itoa(tm.Episode)
		},
	)
	titleList.OnSelected = func(id widget.ListItemID) {
		titleList.Refresh()
		tm := items[id]
		openMedia(tm)
	}

	grid := container.New(layout.NewGridLayout(2), titleList)
	window.SetContent(grid)
	window.ShowAndRun()
	return nil
}

func openMedia(tm MediaItem) {
	cmd := exec.Command("cmd", "/C", "start", "", tm.Path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	log.Println(cmd)

	err := cmd.Run()
	if err != nil {
		log.Printf("failed to open %s: %s", tm.Path, err)
		return
	}
}
