package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type application struct {
	db            *Database
	items         []MediaItem
	filteredItems []MediaItem
	selectedItem  *MediaItem

	// fyne
	app fyne.App

	// explorer
	explorerContainer fyne.CanvasObject

	// detail
	detailContainer fyne.CanvasObject
	titleValue      *widget.Label
	pathValue       *widget.Label
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

	var err error
	app.items, err = app.db.LoadItems()
	if err != nil {
		return err
	}
	app.filteredItems = app.items

	app.initExplorer()
	app.initDetail()

	mainSplit := container.NewHSplit(app.explorerContainer, app.detailContainer)
	mainSplit.SetOffset(0.6)
	window.SetContent(mainSplit)
	window.ShowAndRun()
	return nil
}

func (app *application) initExplorer() {
	filterInput := widget.NewEntry()
	filterInput.SetPlaceHolder("Enter query...")
	filterInput.OnSubmitted = func(s string) {
		// TODO
		app.filteredItems = nil
		for _, item := range app.items {
			if strings.Contains(strings.ToLower(item.Title), s) ||
				strings.Contains(strings.ToLower(filepath.Base(item.Path)), s) {
				app.filteredItems = append(app.filteredItems, item)
			}
		}
		app.explorerContainer.Refresh()
	}

	filterForm := container.NewVBox(filterInput)

	itemList := widget.NewList(
		func() int {
			return len(app.filteredItems)
		},
		func() fyne.CanvasObject {
			titleText := widget.NewLabel("title")
			episodeText := widget.NewLabel("episode")
			return container.New(layout.NewGridLayout(2), titleText, episodeText)
		},
		func(item widget.ListItemID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			titleText := container.Objects[0].(*widget.Label)
			episodeText := container.Objects[1].(*widget.Label)

			tm := app.filteredItems[item]
			if tm.Title != "" {
				titleText.SetText(tm.Title)
			} else {
				titleText.SetText(filepath.Base(tm.Path))
			}
			episodeText.SetText(strconv.Itoa(tm.Episode))
		},
	)
	itemList.OnSelected = func(id widget.ListItemID) {
		itemList.Refresh()
		app.selectedItem = &app.filteredItems[id]
		app.updateDetail()
	}

	split := container.NewVSplit(filterForm, itemList)
	split.SetOffset(0.2)
	app.explorerContainer = split
}

func (app *application) initDetail() {
	titleLabel := widget.NewLabel("Title")
	app.titleValue = widget.NewLabel("")
	app.titleValue.Wrapping = fyne.TextTruncate

	pathLabel := widget.NewLabel("Path")
	app.pathValue = widget.NewLabel("")
	app.pathValue.Wrapping = fyne.TextTruncate

	openButton := widget.NewButton("Open", func() {
		openMedia(*app.selectedItem)
	})

	app.detailContainer = container.NewVBox(
		container.New(layout.NewFormLayout(),
			titleLabel, app.titleValue,
			pathLabel, app.pathValue,
		),
		openButton,
	)
}

func (app *application) updateDetail() {
	if app.selectedItem == nil {
		return
	}

	app.titleValue.SetText(app.selectedItem.Title)
	app.pathValue.SetText(app.selectedItem.Path)

	app.detailContainer.Refresh()
}

func openMedia(tm MediaItem) {
	cmd := exec.Command("cmd", "/C", "start", "", tm.Path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Printf("failed to open %s: %s", tm.Path, err)
		return
	}
}
