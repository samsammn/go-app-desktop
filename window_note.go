package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var selectedItem widget.ListItemID = -1
var titleList = []string{}
var notesList = []string{}

func resegmentMarkdown(segments []widget.RichTextSegment) (newSegments []widget.RichTextSegment) {
	for i, seg := range segments {
		var isList bool
		if len(segments) > i+1 {
			_, isList = segments[i+1].(*widget.ListSegment)
		}

		var nextIsImage bool
		if len(segments) > i+1 {
			_, nextIsImage = segments[i+1].(*widget.ImageSegment)
		}

		var previousIsImage bool
		if i > 0 && len(segments) > i-1 {
			_, previousIsImage = segments[i-1].(*widget.ImageSegment)
		}

		if (seg.Textual() == "" && !isList) || nextIsImage {
			newSegments = append(newSegments, &widget.TextSegment{Text: ""})
		}

		if previousIsImage && seg.Textual() == "" {
			continue
		} else {
			newSegments = append(newSegments, seg)
		}

		var textSeg, isText = seg.(*widget.TextSegment)
		if isText {
			if textSeg.Style.SizeName == theme.SizeNameHeadingText || textSeg.Style.SizeName == theme.SizeNameSubHeadingText {
				newSegments = append(newSegments, &widget.TextSegment{Text: ""})
			}
		}
	}

	return newSegments
}

func windowNote(w fyne.Window) func() {
	return func() {
		notePreviewToolbar.ToolbarObject().Show()

		titleInput := widget.NewEntry()
		titleInput.SetPlaceHolder("Title here...")

		noteInput := widget.NewMultiLineEntry()
		noteInput.SetPlaceHolder("Notes here...")
		noteInput.Wrapping = fyne.TextWrapWord

		notePreview := widget.NewRichText()
		notePreview.Wrapping = fyne.TextWrapWord

		list := widget.NewList(
			func() int {
				return len(titleList)
			},
			func() fyne.CanvasObject {
				lbl := widget.NewLabel("")
				lbl.Wrapping = fyne.TextWrapWord

				return container.NewStack(lbl)
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				lbl := o.(*fyne.Container).Objects[0].(*widget.Label)
				lbl.SetText(titleList[i])
			},
		)

		addBtn := widget.NewButtonWithIcon("Add New Note", theme.ContentAddIcon(), func() {
			// add new notes on the first list position
			titleList = append([]string{"New Title here..."}, titleList...)
			notesList = append([]string{"New Notes here..."}, notesList...)
			selectedItem = 0

			list.Refresh()
			list.Select(selectedItem)

			titleInput.SetText("New Title here...")
			noteInput.SetPlaceHolder("New Notes here...")
		})

		saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
			if selectedItem >= 0 && selectedItem < len(titleList) {
				titleList[selectedItem] = titleInput.Text
				notesList[selectedItem] = noteInput.Text
				list.Refresh()

				dialog.ShowInformation("Saved", "Your notes have been saved.", w)
				w.Content().Refresh()
			}
		})

		scrollEditor := container.NewVScroll(noteInput)
		scrollPreview := container.NewVScroll(notePreview)
		editorAndPreview := container.NewStack(scrollEditor)

		updateContentPreview := func() {
			notePreview.ParseMarkdown(noteInput.Text)
			notePreview.Segments = resegmentMarkdown(notePreview.Segments)
			notePreview.Refresh()
		}

		notePreviewToolbar.OnActivated = func() {
			var isToggleOn = notePreviewToolbar.Icon.Name() == theme.VisibilityIcon().Name()
			if isToggleOn {
				notePreviewToolbar.SetIcon(theme.VisibilityOffIcon())

				editorAndPreview.Objects = []fyne.CanvasObject{scrollEditor}
				editorAndPreview.Refresh()
			} else {
				notePreviewToolbar.SetIcon(theme.VisibilityIcon())

				editorAndPreview.Objects = []fyne.CanvasObject{scrollPreview}
				editorAndPreview.Refresh()

				updateContentPreview()
			}
		}

		list.OnSelected = func(id widget.ListItemID) {
			selectedItem = id
			titleInput.SetText(titleList[id])
			noteInput.SetText(notesList[id])

			updateContentPreview()
		}

		if c := w.Canvas(); c != nil {
			c.AddShortcut(&desktop.CustomShortcut{
				KeyName:  fyne.KeyS,
				Modifier: fyne.KeyModifierControl,
			}, func(shortcut fyne.Shortcut) {
				saveBtn.OnTapped()
			})
		}

		listHeader := container.NewBorder(nil, nil, nil, addBtn)
		leftPannel := container.NewBorder(listHeader, nil, nil, nil, list)

		notesHeader := container.NewBorder(nil, nil, nil, saveBtn, titleInput)
		rightPanel := container.NewBorder(notesHeader, nil, nil, nil, editorAndPreview)

		split := container.NewHSplit(padding(5, leftPannel), padding(5, rightPanel))
		split.SetOffset(0.3)

		w.SetContent(
			container.NewBorder(toolbars(w), nil, nil, nil, split),
		)
	}
}
