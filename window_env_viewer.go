package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type EnvViewer struct {
	ProjectName string
	EnvPath     string
}

var currentEnvData string = ""
var listSourceEnvApps = []EnvViewer{
	{
		ProjectName: "Linkreator",
		EnvPath:     `D:\Workspace\MEA\mea-lp-builder-api\.env`,
	},
}

func catEnvFile(path string) string {
	fileByte, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("err", err)
	}

	return string(fileByte)
}

func envToRichText(input string) *widget.RichText {
	rt := widget.NewRichText()
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		switch {
		case strings.HasPrefix(strings.TrimSpace(line), "#"):
			rt.Segments = append(rt.Segments, &widget.TextSegment{
				Text: line + "\n",
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameDisabled,
				},
			})
		case strings.Contains(line, "="):
			parts := strings.SplitN(line, "=", 2)
			rt.Segments = append(rt.Segments,
				&widget.TextSegment{
					Text: parts[0] + "=",
					Style: widget.RichTextStyle{
						Inline:    true,
						TextStyle: fyne.TextStyle{Bold: true},
					},
				},
				&widget.TextSegment{
					Text: parts[1] + "",
				},
			)
		default:
			rt.Segments = append(rt.Segments, &widget.TextSegment{Text: line + "\n"})
		}
	}

	return rt
}

func envToJSON(input string) (string, error) {
	result := make(map[string]interface{})

	lines := strings.Split(input, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// skip empty & comment
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// remove quotes
		if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
			val = strings.Trim(val, `"`)
			result[key] = val
			continue
		}

		// bool
		if val == "true" || val == "false" {
			result[key] = (val == "true")
			continue
		}

		// int
		if i, err := strconv.Atoi(val); err == nil {
			result[key] = i
			continue
		}

		// float
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			result[key] = f
			continue
		}

		// default string
		result[key] = val
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func jsonToRichText(input string) *widget.RichText {
	rt := widget.NewRichText()
	var segs []widget.RichTextSegment

	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		segs = append(segs, &widget.TextSegment{
			Text: input,
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameError,
			},
		})

		rt.Segments = segs
		return rt
	}

	formatted, _ := json.MarshalIndent(data, "", "    ")
	lines := strings.Split(string(formatted), "\n")

	for _, line := range lines {
		switch {
		case strings.Contains(line, ":"):
			parts := strings.SplitN(line, ":", 2)
			segs = append(segs, &widget.TextSegment{
				Text: parts[0] + ":",
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNamePrimary,
					Inline:    true,
				},
			}, &widget.TextSegment{
				Text: parts[1] + "\n",
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					Inline:    true,
				},
			})
		default:
			segs = append(segs, &widget.TextSegment{Text: line + "\n"})
		}
	}

	rt.Segments = segs
	return rt
}

func windowEnvViewer(w fyne.Window) func() {
	return func() {
		notePreviewToolbar.ToolbarObject().Hide()

		viewOptions := widget.NewSelectEntry([]string{"Markdown", "Textarea", "JSON"})
		viewOptions.SetText("Markdown")

		list := widget.NewList(
			func() int {
				return len(listSourceEnvApps)
			},
			func() fyne.CanvasObject {
				lbl := widget.NewLabel("")
				lbl.Wrapping = fyne.TextWrapWord

				return container.NewStack(lbl)
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				lbl := o.(*fyne.Container).Objects[0].(*widget.Label)
				lbl.SetText(listSourceEnvApps[i].ProjectName)
			},
		)

		richText := widget.NewRichText()
		richText.Wrapping = fyne.TextWrapBreak

		textArea := widget.NewMultiLineEntry()
		textArea.Wrapping = fyne.TextWrapBreak
		containerValue := container.NewStack(richText)

		list.OnSelected = func(id widget.ListItemID) {
			envPath := listSourceEnvApps[id].EnvPath
			envData := catEnvFile(envPath)
			viewOption := viewOptions.Text
			currentEnvData = envData

			if viewOption == "Markdown" {
				richText = envToRichText(envData)
				richText.Wrapping = fyne.TextWrapBreak

				containerValue.Objects = []fyne.CanvasObject{richText}
			} else if viewOption == "JSON" {
				jsonString, _ := envToJSON(envData)
				richText = jsonToRichText(jsonString)
				richText.Wrapping = fyne.TextWrapBreak

				containerValue.Objects = []fyne.CanvasObject{richText}
			} else {
				textArea.SetText(envData)

				containerValue.Objects = []fyne.CanvasObject{textArea}
			}

			containerValue.Refresh()
		}

		viewOptions.OnChanged = func(s string) {
			if s == "Markdown" {
				richText = envToRichText(currentEnvData)
				richText.Wrapping = fyne.TextWrapBreak

				containerValue.Objects = []fyne.CanvasObject{richText}
			} else if s == "JSON" {
				jsonString, _ := envToJSON(currentEnvData)
				richText = jsonToRichText(jsonString)
				richText.Wrapping = fyne.TextWrapBreak

				containerValue.Objects = []fyne.CanvasObject{richText}
			} else {
				textArea.SetText(currentEnvData)

				containerValue.Objects = []fyne.CanvasObject{textArea}
			}

			containerValue.Refresh()
		}

		viewDropdown := container.NewGridWrap(fyne.NewSize(130, 40), viewOptions)
		btnAdd := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			var modal *widget.PopUp

			projectName := widget.NewEntry()
			projectPath := widget.NewEntry()

			btnInsert := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
				listSourceEnvApps = append(listSourceEnvApps, EnvViewer{
					ProjectName: projectName.Text,
					EnvPath:     projectPath.Text,
				})

				list.Refresh()
				modal.Hide()
			})

			btnCancel := widget.NewButtonWithIcon("Cancel", theme.WindowCloseIcon(), func() {
				modal.Hide()
			})

			projectName.SetPlaceHolder("Enter Project Name")
			projectPath.SetPlaceHolder("Enter .env location")

			content := container.NewVBox(projectName, projectPath, vPadding(10), container.NewHBox(btnInsert, btnCancel))

			grid := container.NewGridWrap(fyne.NewSize(300, 175), padding(20, content))

			modal = widget.NewModalPopUp(container.NewCenter(grid), w.Canvas())
			modal.Show()
		})

		btnCopy := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
			if viewOptions.Text == "Textarea" {
				w.Clipboard().SetContent(textArea.Text)
			} else {
				w.Clipboard().SetContent(richText.String())
			}
		})

		labelProject := widget.NewLabel("List Project")
		labelProject.SizeName = theme.SizeNameSubHeadingText
		labelProject.TextStyle = fyne.TextStyle{
			Bold: true,
		}

		actionBtn := container.NewBorder(nil, nil, container.NewHBox(btnAdd, labelProject), container.NewHBox(viewDropdown, btnCopy))

		rightPanel := container.NewVScroll(containerValue)
		leftPanel := container.NewGridWrap(fyne.NewSize(240, w.Content().Size().Height), list)
		topPanel := container.NewVBox(
			widget.NewSeparator(),
			actionBtn,
			widget.NewSeparator(),
		)

		content := container.NewBorder(topPanel, nil, leftPanel, nil, rightPanel)

		w.SetContent(
			container.NewBorder(toolbars(w), nil, nil, nil, content),
		)
	}
}
