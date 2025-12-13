package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func setVPadding(v float32) *canvas.Rectangle {
	padding := canvas.NewRectangle(color.Transparent)
	padding.SetMinSize(fyne.NewSize(0, v))

	return padding
}

func NewWelcomePanel(windowNote func(), windowTimeBreak func(), windowBase64 func(), windowStringGenerator func(), windowQRGenerator func()) fyne.CanvasObject {
	img := canvas.NewImageFromResource(theme.HomeIcon())
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(100, 100))

	title := canvas.NewText("Sesterdamp Apps", theme.Color(theme.ColorNameForeground))
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 36

	subtitle := canvas.NewText("Secure • Simple • Fast", theme.Color(theme.ColorNameForeground))
	subtitle.TextSize = 16

	header := container.NewVBox(
		container.NewCenter(img),
		container.NewCenter(title),
		container.NewCenter(subtitle),
	)

	cardNotes := widget.NewCard(
		"My Notes",
		"Mulai membuat catatan baru",
		widget.NewButton("Open", func() {
			if windowNote != nil {
				windowNote()
			}
		}),
	)

	cardTimeBreak := widget.NewCard(
		"Time for Break",
		"Ingatkan untuk istirahat secara berkala",
		widget.NewButton("Open", func() {
			if windowTimeBreak != nil {
				windowTimeBreak()
			}
		}),
	)

	cardBase64 := widget.NewCard(
		"Base64 Manager",
		"Encode dan Decode teks dengan Base64",
		widget.NewButton("Open", func() {
			if windowBase64 != nil {
				windowBase64()
			}
		}),
	)

	cardString := widget.NewCard(
		"String Generator",
		"Buat string acak dengan berbagai opsi",
		widget.NewButton("Open", func() {
			if windowStringGenerator != nil {
				windowStringGenerator()
			}
		}),
	)

	cardQRGenerator := widget.NewCard(
		"QR Generator",
		"Buat kode QR dari teks atau URL",
		widget.NewButton("Open", func() {
			if windowQRGenerator != nil {
				windowQRGenerator()
			}
		}),
	)

	cardComingSoon := widget.NewCard(
		"Coming Soon",
		"Fitur-fitur menarik lainnya akan segera hadir!",
		widget.NewButton("Open", func() {}),
	)

	grid := container.NewGridWithColumns(3,
		cardNotes,
		cardTimeBreak,
		cardBase64,
		cardString,
		cardQRGenerator,
		cardComingSoon,
	)

	panel := container.NewVBox(
		container.NewCenter(header),
		setVPadding(50),
		widget.NewSeparator(),
		setVPadding(20),
		container.NewCenter(grid),
	)

	return container.NewCenter(panel)
}

func welcomes(w fyne.Window) fyne.CanvasObject {
	welcome := NewWelcomePanel(
		windowNote(w),
		windowTimeBreak(w),
		windowBase64(w),
		windowStringGenerator(w),
		windowQRGenerator(w),
	)

	return welcome
}

func windowTimeBreak(w fyne.Window) func() {
	return func() {
		w.SetContent(
			container.NewBorder(toolbars(w), nil, nil, nil, widget.NewLabel("Time Break Reminder")),
		)
	}
}

func windowBase64(w fyne.Window) func() {
	return func() {
		w.SetContent(
			container.NewBorder(toolbars(w), nil, nil, nil, widget.NewLabel("Base64 Manager")),
		)
	}
}

func windowQRGenerator2(w fyne.Window) func() {
	return func() {
		w.SetContent(
			container.NewBorder(toolbars(w), nil, nil, nil, widget.NewLabel("QR Generator")),
		)
	}
}

var emptyFunc = func() {
	fmt.Println("click inisiate func!")
}

var notePreviewToolbar *widget.ToolbarAction

func toolbars(w fyne.Window) *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {
			fmt.Println("Home clicked")
			w.SetContent(welcomes(w))
		}),
		notePreviewToolbar,
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			fmt.Println("Settings clicked")
		}),
	)

	return toolbar
}

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())

	w := a.NewWindow("Welcome Panel Sesterdamp")

	notePreviewToolbar = widget.NewToolbarAction(theme.VisibilityOffIcon(), emptyFunc)

	w.SetContent(welcomes(w))
	w.Resize(fyne.NewSize(1200, 800))
	w.ShowAndRun()
}
