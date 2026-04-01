/*
 * ╔════════════════════════════════════════════════════════════════╗
 * ║ TERMINAL                                                       ║
 * ║ Plik / File: main.go                                           ║
 * ╠════════════════════════════════════════════════════════════════╣
 * ║ Autor / Author:                                                ║
 * ║   GufiM                                                        ║
 * ╠════════════════════════════════════════════════════════════════╣
 * ║ GitHub  : moj-terminal                                         ║
 * ║ Forum   : https://forum.lothar-team.pl                         ║
 * ║                                                                ║
 * ║ Licencja / License: MIT                                        ║
 * ║ Rok / Year: 2026                                               ║
 * ╚════════════════════════════════════════════════════════════════╝
 */
package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
)

// --- STAŁE APLIKACJI ---
const AppVersion = "0.0.20"

// --- KOLORY GLOBALNE ---
var (
	ColorBackground = color.RGBA{10, 10, 10, 255}
	ColorTerminalBg = color.RGBA{20, 20, 20, 255}
	ColorInputBg    = color.RGBA{30, 30, 10, 255}
	ColorTextSystem = color.RGBA{0, 255, 255, 255}
	ColorTextError  = color.RGBA{255, 0, 0, 255}
	ColorIndIdle    = color.RGBA{100, 100, 100, 255}
	ColorBlue       = color.RGBA{0, 150, 255, 255}
	ColorRed        = color.RGBA{255, 0, 0, 255}
	ColorGreen      = color.RGBA{0, 200, 0, 255}
	ColorOrange     = color.RGBA{255, 165, 0, 255}
	ColorPurple     = color.RGBA{160, 32, 240, 255}

	UserColors = []color.Color{
		color.White,
		color.RGBA{255, 230, 100, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 200, 255, 255},
		color.RGBA{255, 100, 100, 255},
		color.RGBA{200, 100, 255, 255},
	}
	ColorNames = []string{"White", "Sand", "Green", "Cyan", "Red", "Pink"}
)

type Config struct {
	Port          string  `json:"port"`
	BaudRate      int     `json:"baud"`
	DataBits      int     `json:"data"`
	Parity        string  `json:"parity"`
	StopBits      string  `json:"stop"`
	Echo          bool    `json:"echo"`
	ShowLR        bool    `json:"show_lr"`
	ShowHEX       bool    `json:"show_hex"`
	HexFormat     string  `json:"hex_format"`
	FontSize      float32 `json:"font_size"`
	ShowLines     bool    `json:"show_lines"`
	ShowTimestamp bool    `json:"show_timestamp"`
	TXBold        bool    `json:"tx_bold"`
	RXBold        bool    `json:"rx_bold"`
	TXColorIdx    int     `json:"tx_color_idx"`
	RXColorIdx    int     `json:"rx_color_idx"`
}

var (
	activePort   serial.Port
	keepReading  bool
	lineCounter  int = 1
	conf         Config
	terminalBox  *fyne.Container
	inputEntry   *CustomEntry
	portSelect   *tappableSelect
	baudSelect   *widget.Select
	dataSelect   *widget.Select
	paritySelect *widget.Select
	stopSelect   *widget.Select
	chkEchoMain  *widget.Check
	chkLRMain    *widget.Check
	chkHEXMain   *widget.Check
)

type tappableSelect struct {
	widget.Select
	OnTapped func()
}

func newTappableSelect(options []string, onTapped func()) *tappableSelect {
	s := &tappableSelect{OnTapped: onTapped}
	s.Options = options
	s.ExtendBaseWidget(s)
	return s
}
func (s *tappableSelect) Tapped(e *fyne.PointEvent) {
	if s.OnTapped != nil {
		s.OnTapped()
	}
	s.Select.Tapped(e)
}

type customEntryRenderer struct{ fyne.WidgetRenderer }

func (r *customEntryRenderer) Objects() []fyne.CanvasObject {
	objs := r.WidgetRenderer.Objects()
	for _, obj := range objs {
		if txt, ok := obj.(*canvas.Text); ok {
			txt.Color = UserColors[conf.TXColorIdx]
			txt.TextStyle = fyne.TextStyle{Bold: conf.TXBold}
		}
	}
	return objs
}

type CustomEntry struct{ widget.Entry }

func NewCustomEntry() *CustomEntry {
	e := &CustomEntry{}
	e.ExtendBaseWidget(e)
	return e
}
func (e *CustomEntry) CreateRenderer() fyne.WidgetRenderer {
	return &customEntryRenderer{e.Entry.CreateRenderer()}
}

func saveSettingsToFile() {
	data, _ := json.MarshalIndent(conf, "", "  ")
	os.WriteFile("settings.json", data, 0644)
}

func createBlinker(w, h float32) *canvas.Rectangle {
	r := canvas.NewRectangle(color.Transparent)
	r.StrokeColor = ColorIndIdle
	r.StrokeWidth = 2
	r.CornerRadius = 4
	r.Resize(fyne.NewSize(w, h))
	return r
}

func wrapBtn(blk *canvas.Rectangle, btn *widget.Button) *fyne.Container {
	return container.New(layout.NewGridWrapLayout(fyne.NewSize(85, 35)),
		container.NewStack(blk, container.NewCenter(container.New(layout.NewGridWrapLayout(fyne.NewSize(81, 31)), btn))))
}

func flash(r *canvas.Rectangle, c color.Color) {
	go func() {
		r.StrokeColor = c
		r.Refresh()
		time.Sleep(time.Millisecond * 200)
		r.StrokeColor = ColorIndIdle
		r.Refresh()
	}()
}

type myTheme struct{}

func (m myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if n == theme.ColorNameBackground {
		return ColorBackground
	}
	if n == theme.ColorNameInputBackground {
		return ColorInputBg
	}
	return theme.DefaultTheme().Color(n, v)
}
func (m myTheme) Font(s fyne.TextStyle) fyne.Resource     { return theme.DefaultTheme().Font(s) }
func (m myTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(n) }
func (m myTheme) Size(n fyne.ThemeSizeName) float32       { return theme.DefaultTheme().Size(n) }

func main() {
	myApp := app.NewWithID("com.serial.terminal")
	myApp.Settings().SetTheme(&myTheme{})
	win := myApp.NewWindow("SERIAL TERMINAL")
	win.Resize(fyne.NewSize(1050, 650))

	if res, err := fyne.LoadResourceFromPath("terminal.png"); err == nil {
		win.SetIcon(res)
		myApp.SetIcon(res)
	}

	conf = Config{BaudRate: 9600, DataBits: 8, Parity: "None", StopBits: "1", Echo: true, RXColorIdx: 2, HexFormat: "XX", FontSize: 14}
	if data, err := os.ReadFile("settings.json"); err == nil {
		json.Unmarshal(data, &conf)
	}

	terminalBox = container.NewVBox()
	terminalScroll := container.NewScroll(terminalBox)
	terminalArea := container.NewStack(canvas.NewRectangle(ColorTerminalBg), container.NewPadded(terminalScroll))

	appendTerminal := func(txt string, col color.Color, isBold bool, isSystem bool) {
		if txt == "" {
			return
		}
		displayTxt := txt

		if !isSystem && conf.ShowHEX {
			var hexParts []string
			for i := 0; i < len(txt); i++ {
				f := "%02X"
				switch conf.HexFormat {
				case "0xXX":
					f = "0x%02X"
				case "\\xXX":
					f = "\\x%02X"
				}
				hexParts = append(hexParts, fmt.Sprintf(f, txt[i]))
			}
			displayTxt = strings.Join(hexParts, " ")
		} else {
			if conf.ShowLR {
				displayTxt = strings.ReplaceAll(displayTxt, "\r", "[CR]")
				displayTxt = strings.ReplaceAll(displayTxt, "\n", "[LF]")
			} else {
				displayTxt = strings.NewReplacer("\r", "", "\n", "").Replace(displayTxt)
			}
		}

		if displayTxt == "" {
			return
		}

		prefix := ""
		if conf.ShowLines {
			prefix = fmt.Sprintf("[%03d] ", lineCounter)
			lineCounter++
		}
		if conf.ShowTimestamp {
			prefix += fmt.Sprintf("[%s] ", time.Now().Format("15:04:05.000"))
		}

		t := canvas.NewText(prefix+displayTxt, col)
		t.TextSize = conf.FontSize
		t.TextStyle = fyne.TextStyle{Bold: isBold}
		terminalBox.Add(t)
		terminalBox.Refresh()
		terminalScroll.ScrollToBottom()
	}

	txRect := createBlinker(45, 35)
	rxRect := createBlinker(45, 35)
	txInd := container.New(layout.NewGridWrapLayout(fyne.NewSize(45, 35)), container.NewStack(txRect, container.NewCenter(canvas.NewText("TX", color.White))))
	rxInd := container.New(layout.NewGridWrapLayout(fyne.NewSize(45, 35)), container.NewStack(rxRect, container.NewCenter(canvas.NewText("RX", color.White))))

	startBlk := createBlinker(85, 35)
	stopBlk := createBlinker(85, 35)
	sendBlk := createBlinker(85, 35)
	setBlk := createBlinker(85, 35)
	saveBlk := createBlinker(85, 35)
	clrBlk := createBlinker(85, 35)

	inputEntry = NewCustomEntry()
	inputEntry.SetPlaceHolder("Message...")
	inputArea := container.NewStack(canvas.NewRectangle(ColorInputBg), inputEntry)

	refreshPorts := func() {
		pList, _ := serial.GetPortsList()
		if len(pList) > 0 {
			portSelect.Options = pList
			if portSelect.Selected == "" {
				portSelect.SetSelected(pList[0])
			}
		} else {
			portSelect.Options = []string{}
			portSelect.ClearSelected()
		}
		portSelect.Refresh()
	}
	portSelect = newTappableSelect([]string{}, refreshPorts)
	portSelect.PlaceHolder = "FIND PORT"

	baudList := []string{"300", "600", "1200", "2400", "4800", "9600", "14400", "19200", "38400", "57600", "115200", "230400", "460800", "921800"}
	baudSelect = widget.NewSelect(baudList, nil)
	dataSelect = widget.NewSelect([]string{"5", "6", "7", "8"}, nil)
	paritySelect = widget.NewSelect([]string{"None", "Odd", "Even", "Mark", "Space"}, nil)
	stopSelect = widget.NewSelect([]string{"1", "1.5", "2"}, nil)

	openSettings := func() {
		flash(setBlk, ColorOrange)
		sWin := myApp.NewWindow("SETTINGS")
		sWin.Resize(fyne.NewSize(450, 680)) // Lekko powiększone okno

		lineChk := widget.NewCheck("Show Line Numbers", func(b bool) { conf.ShowLines = b })
		lineChk.SetChecked(conf.ShowLines)
		timeChk := widget.NewCheck("Show Timestamp", func(b bool) { conf.ShowTimestamp = b })
		timeChk.SetChecked(conf.ShowTimestamp)

		txBld := widget.NewCheck("TX Bold", func(b bool) { conf.TXBold = b; inputEntry.Refresh() })
		txBld.SetChecked(conf.TXBold)
		rxBld := widget.NewCheck("RX Bold", func(b bool) { conf.RXBold = b })
		rxBld.SetChecked(conf.RXBold)

		txC := widget.NewSelect(ColorNames, func(s string) {
			for i, n := range ColorNames {
				if n == s {
					conf.TXColorIdx = i
					break
				}
			}
			inputEntry.Refresh()
		})
		txC.SetSelectedIndex(conf.TXColorIdx)
		rxC := widget.NewSelect(ColorNames, func(s string) {
			for i, n := range ColorNames {
				if n == s {
					conf.RXColorIdx = i
					break
				}
			}
		})
		rxC.SetSelectedIndex(conf.RXColorIdx)

		sizeLabel := widget.NewLabel(fmt.Sprintf("Font Size: %.0f", conf.FontSize))
		sizeSlider := widget.NewSlider(8, 24)
		sizeSlider.SetValue(float64(conf.FontSize))
		sizeSlider.OnChanged = func(v float64) {
			conf.FontSize = float32(v)
			sizeLabel.SetText(fmt.Sprintf("Font Size: %.0f", conf.FontSize))
		}

		sBaud := widget.NewSelect(baudList, func(s string) { baudSelect.SetSelected(s); conf.BaudRate, _ = strconv.Atoi(s) })
		sBaud.SetSelected(baudSelect.Selected)
		sData := widget.NewSelect(dataSelect.Options, func(s string) { dataSelect.SetSelected(s); conf.DataBits, _ = strconv.Atoi(s) })
		sData.SetSelected(dataSelect.Selected)
		sPar := widget.NewSelect(paritySelect.Options, func(s string) { paritySelect.SetSelected(s); conf.Parity = s })
		sPar.SetSelected(paritySelect.Selected)
		sStop := widget.NewSelect(stopSelect.Options, func(s string) { stopSelect.SetSelected(s); conf.StopBits = s })
		sStop.SetSelected(stopSelect.Selected)

		sEcho := widget.NewCheck("Echo Local", func(b bool) { chkEchoMain.SetChecked(b); conf.Echo = b })
		sEcho.SetChecked(conf.Echo)
		sLR := widget.NewCheck("Show [CR][LF]", func(b bool) { chkLRMain.SetChecked(b); conf.ShowLR = b })
		sLR.SetChecked(conf.ShowLR)
		sHEX := widget.NewCheck("Show HEX", func(b bool) { chkHEXMain.SetChecked(b); conf.ShowHEX = b })
		sHEX.SetChecked(conf.ShowHEX)

		hexFmt := widget.NewSelect([]string{"XX", "0xXX", "\\xXX"}, func(s string) { conf.HexFormat = s })
		hexFmt.SetSelected(conf.HexFormat)

		btnSaveExit := container.New(layout.NewGridWrapLayout(fyne.NewSize(160, 35)),
			container.NewStack(createBlinker(160, 35), container.NewCenter(container.New(layout.NewGridWrapLayout(fyne.NewSize(156, 31)),
				widget.NewButton("SAVE & EXIT", func() { saveSettingsToFile(); sWin.Close() })))))

		// ETYKIETA WERSJI (na dole z lewej)
		versionLabel := widget.NewLabel("version:" + AppVersion)

		// DOLNY PASEK Z WERSJĄ I PRZYCISKIEM
		bottomFooter := container.NewBorder(nil, nil, versionLabel, nil, container.NewCenter(btnSaveExit))

		sWin.SetContent(container.NewPadded(container.NewVBox(
			widget.NewLabelWithStyle("VISUALS", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewGridWithColumns(2, lineChk, timeChk),
			container.NewGridWithColumns(2, txBld, rxBld),
			container.NewGridWithColumns(2, widget.NewLabel("TX Color:"), txC),
			container.NewGridWithColumns(2, widget.NewLabel("RX Color:"), rxC),
			container.NewGridWithColumns(2, sizeLabel, sizeSlider),
			widget.NewSeparator(),
			widget.NewLabelWithStyle("SERIAL CONFIG", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewGridWithColumns(2, widget.NewLabel("Baud:"), sBaud),
			container.NewGridWithColumns(2, widget.NewLabel("Data Bits:"), sData),
			container.NewGridWithColumns(2, widget.NewLabel("Parity:"), sPar),
			container.NewGridWithColumns(2, widget.NewLabel("Stop Bits:"), sStop),
			container.NewGridWithColumns(2, widget.NewLabel("HEX Format:"), hexFmt),
			container.NewGridWithColumns(3, sEcho, sLR, sHEX),
			layout.NewSpacer(),
			bottomFooter, // Dodanie stopki z wersją
		)))
		sWin.Show()
	}

	btnStart := wrapBtn(startBlk, widget.NewButton("START", func() {
		flash(startBlk, ColorBlue)
		if portSelect.Selected == "" {
			return
		}
		b, _ := strconv.Atoi(baudSelect.Selected)
		d, _ := strconv.Atoi(dataSelect.Selected)
		var pMode serial.Parity
		switch paritySelect.Selected {
		case "Odd":
			pMode = serial.OddParity
		case "Even":
			pMode = serial.EvenParity
		default:
			pMode = serial.NoParity
		}
		p, err := serial.Open(portSelect.Selected, &serial.Mode{BaudRate: b, DataBits: d, Parity: pMode})
		if err != nil {
			appendTerminal("Error: "+err.Error(), ColorTextError, true, true)
			return
		}
		activePort = p
		keepReading = true
		appendTerminal("Connected to "+portSelect.Selected, ColorTextSystem, true, true)
		go func() {
			buff := make([]byte, 1024)
			for keepReading {
				n, _ := activePort.Read(buff)
				if n > 0 {
					flash(rxRect, UserColors[conf.RXColorIdx])
					appendTerminal(string(buff[:n]), UserColors[conf.RXColorIdx], conf.RXBold, false)
				}
			}
		}()
	}))

	btnStop := wrapBtn(stopBlk, widget.NewButton("STOP", func() {
		flash(stopBlk, ColorRed)
		keepReading = false
		if activePort != nil {
			activePort.Close()
		}
		appendTerminal("Disconnected", ColorTextSystem, true, true)
	}))

	sendFunc := func() {
		if inputEntry.Text == "" || activePort == nil {
			return
		}
		activePort.Write([]byte(inputEntry.Text + "\r\n"))
		flash(txRect, UserColors[conf.TXColorIdx])
		flash(sendBlk, UserColors[conf.TXColorIdx])
		if conf.Echo {
			appendTerminal(inputEntry.Text, UserColors[conf.TXColorIdx], conf.TXBold, false)
		}
		inputEntry.SetText("")
	}

	btnSend := wrapBtn(sendBlk, widget.NewButton("SEND", sendFunc))
	inputEntry.OnSubmitted = func(s string) { sendFunc() }

	btnSettings := wrapBtn(setBlk, widget.NewButton("SETTINGS", openSettings))
	btnSave := wrapBtn(saveBlk, widget.NewButton("SAVE", func() {
		flash(saveBlk, ColorGreen)
		conf.Port = portSelect.Selected
		saveSettingsToFile()
		appendTerminal("Settings saved", ColorTextSystem, false, true)
	}))
	btnClear := wrapBtn(clrBlk, widget.NewButton("CLEAR", func() {
		flash(clrBlk, ColorPurple)
		terminalBox.Objects = nil
		lineCounter = 1
		terminalBox.Refresh()
	}))

	chkEchoMain = widget.NewCheck("ECHO", func(b bool) { conf.Echo = b })
	chkEchoMain.SetChecked(conf.Echo)
	chkLRMain = widget.NewCheck("LR", func(b bool) { conf.ShowLR = b })
	chkLRMain.SetChecked(conf.ShowLR)
	chkHEXMain = widget.NewCheck("HEX", func(b bool) { conf.ShowHEX = b })
	chkHEXMain.SetChecked(conf.ShowHEX)

	topBar := container.New(layout.NewGridLayout(5), portSelect, baudSelect, dataSelect, paritySelect, stopSelect)
	topRow := container.NewBorder(nil, nil, container.NewHBox(btnStart, txInd), btnSend, container.NewStack(container.New(layout.NewGridWrapLayout(fyne.NewSize(0, 35))), inputArea))
	bottomRow := container.NewBorder(nil, nil, container.NewHBox(btnStop, rxInd), container.NewHBox(layout.NewSpacer(), btnSettings, btnSave, btnClear), container.NewHBox(chkEchoMain, chkLRMain, chkHEXMain))

	go func() {
		time.Sleep(time.Millisecond * 200)
		pList, _ := serial.GetPortsList()
		if len(pList) > 0 {
			portSelect.Options = pList
			foundSaved := false
			if conf.Port != "" {
				for _, p := range pList {
					if p == conf.Port {
						portSelect.SetSelected(conf.Port)
						foundSaved = true
						break
					}
				}
			}
			if !foundSaved {
				portSelect.SetSelected(pList[0])
			}
		} else {
			portSelect.Options = []string{}
			portSelect.ClearSelected()
		}
		portSelect.Refresh()
	}()

	baudSelect.SetSelected(strconv.Itoa(conf.BaudRate))
	dataSelect.SetSelected(strconv.Itoa(conf.DataBits))
	paritySelect.SetSelected(conf.Parity)
	stopSelect.SetSelected(conf.StopBits)

	win.SetContent(container.NewBorder(topBar, container.NewPadded(container.NewVBox(topRow, bottomRow)), nil, nil, terminalArea))
	win.ShowAndRun()
}
