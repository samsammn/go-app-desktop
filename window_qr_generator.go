package main

import (
	"bytes"
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

var listQRLogoName = []string{"Sebari", "Linkreator"}
var listQRLogoAssets = map[string]string{
	"Sebari":     "assets/sebari_logo_qr_code.png",
	"Linkreator": "assets/linkreator_logo_qr_code.png",
}

func generateQR(text string, logoName string, widthImage uint8, widthBorder int) ([]byte, error) {
	qr, err := qrcode.New(text)
	if err != nil {
		return nil, err
	}

	var bufQRImage bytes.Buffer
	var logoAssets = listQRLogoAssets[logoName]

	pr, pw := io.Pipe()

	defer pr.Close()
	defer pw.Close()

	go func() {
		writer := standard.NewWithWriter(pw,
			standard.WithQRWidth(widthImage),
			standard.WithBorderWidth(widthBorder),
			standard.WithLogoImageFilePNG(logoAssets),
			standard.WithLogoSizeMultiplier(1),
		)

		err := qr.Save(writer)
		if err != nil {
			return
		}
	}()

	_, err = io.Copy(&bufQRImage, pr)
	if err != nil {
		return nil, err
	}

	return bufQRImage.Bytes(), nil
}

func windowQRGenerator(w fyne.Window) func() {
	return func() {
		notePreviewToolbar.ToolbarObject().Hide()

		textEntry := widget.NewMultiLineEntry()
		textEntry.SetPlaceHolder("Fill any of text do you want")

		widthImageEntry := widget.NewEntry()
		widthImageEntry.SetPlaceHolder("Width Image")
		widthImageEntry.SetText("15")

		logoOptions := widget.NewSelectEntry(listQRLogoName)
		logoOptions.SetText("Sebari")

		widthBorderEntry := widget.NewEntry()
		widthBorderEntry.SetPlaceHolder("Width Border")
		widthBorderEntry.SetText("0")

		resultImage := canvas.NewImageFromResource(theme.BrokenImageIcon())

		generate := func() {
			var widthImage uint8
			var widthBorder int

			if textEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("text must be fill"), w)
				return
			}

			_, err := fmt.Sscan(widthImageEntry.Text, &widthImage)
			if err != nil || widthImage <= 0 {
				dialog.ShowError(fmt.Errorf("invalid length: %v", widthImageEntry.Text), w)
				return
			}

			_, err = fmt.Sscan(widthBorderEntry.Text, &widthBorder)
			if err != nil || widthImage <= 0 {
				dialog.ShowError(fmt.Errorf("invalid length: %v", widthBorderEntry.Text), w)
				return
			}

			qrCode, err := generateQR(textEntry.Text, logoOptions.Text, widthImage, widthBorder)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			qr := canvas.NewImageFromReader(bytes.NewReader(qrCode), "qrImage.png")

			resultImage.Resource = qr.Resource
			resultImage.FillMode = canvas.ImageFillContain
			resultImage.Refresh()

			dialog.ShowInformation("Generate", "QR has been Generated ✔", w)
		}

		generateBtn := widget.NewButtonWithIcon("Generate", theme.DocumentPrintIcon(), nil)
		generateBtn.OnTapped = generate

		clearBtn := widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), func() {
			brokenImage := canvas.NewImageFromResource(theme.BrokenImageIcon())
			resultImage.Resource = brokenImage.Resource
			resultImage.Refresh()

			textEntry.Text = ""
			textEntry.Refresh()
		})

		// shortcuts: Ctrl+G for generate, Ctrl+C for copy
		if c := w.Canvas(); c != nil {
			c.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyG, Modifier: fyne.KeyModifierControl}, func(fyne.Shortcut) {
				generate()
			})
		}

		configs := widget.NewForm()
		configs.Append("Text:", textEntry)
		configs.Append("Width Image:", widthImageEntry)
		configs.Append("Width Border:", widthBorderEntry)
		configs.Append("Logo:", logoOptions)

		left := container.NewVBox(
			widget.NewLabelWithStyle("Configuration", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			configs,
			widget.NewSeparator(),
			container.NewHBox(generateBtn, clearBtn),
		)

		scrollLeft := container.NewVScroll(left)
		scrollLeft.SetMinSize(fyne.NewSize(150, 0))

		imageCard := container.NewGridWrap(fyne.NewSize(300, 300), resultImage)

		rightTop := container.NewBorder(nil, nil, nil, nil, widget.NewLabelWithStyle("Result", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		right := container.NewBorder(rightTop, nil, nil, nil, container.NewCenter(imageCard))

		split := container.NewHSplit(scrollLeft, right)
		split.SetOffset(0.7)

		w.SetContent(
			container.NewBorder(toolbars(w), nil, nil, nil, split),
		)
	}
}
