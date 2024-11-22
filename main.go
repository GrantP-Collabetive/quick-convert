package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Sql struct {
	Name      string `json:"name"`
	Statement string `json:"statement"`
}

var (
	sqlSlice = []Sql{}
	sqlFile  = "sql.json"
)

const pageCount = 3

const logo = `
  /$$$$$$            /$$           /$$              /$$$$$$                                                      /$$    
 /$$__  $$          |__/          | $$             /$$__  $$                                                    | $$    
| $$  \ $$ /$$   /$$ /$$  /$$$$$$$| $$   /$$      | $$  \__/  /$$$$$$  /$$$$$$$  /$$    /$$ /$$$$$$   /$$$$$$  /$$$$$$  
| $$  | $$| $$  | $$| $$ /$$_____/| $$  /$$/      | $$       /$$__  $$| $$__  $$|  $$  /$$//$$__  $$ /$$__  $$|_  $$_/  
| $$  | $$| $$  | $$| $$| $$      | $$$$$$/       | $$      | $$  \ $$| $$  \ $$ \  $$/$$/| $$$$$$$$| $$  \__/  | $$    
| $$/$$ $$| $$  | $$| $$| $$      | $$_  $$       | $$    $$| $$  | $$| $$  | $$  \  $$$/ | $$_____/| $$        | $$ /$$
|  $$$$$$/|  $$$$$$/| $$|  $$$$$$$| $$ \  $$      |  $$$$$$/|  $$$$$$/| $$  | $$   \  $/  |  $$$$$$$| $$        |  $$$$/
 \____ $$$ \______/ |__/ \_______/|__/  \__/       \______/  \______/ |__/  |__/    \_/    \_______/|__/         \___/  
      \__/                                                                                                              
                                                                                                                        
                                                                                                                        
`

const (
	version    = `V0.01`
	subtitle   = `Quick Convert - Quickly relates DB tables to form sql statements in a JSON object for easy use`
	navigation = `[#FF8C00]Ctrl-N[#00FFFF]: Next Page   [#FF8C00]Ctrl-P[#00FFFF]: Previous Page    [#FF8C00]Ctrl-C[#00FFFF]: Exit`
)

var pageNames = []string{"Home", "Query Builder", "List View"}

func loadSql() {
	if _, err := os.Stat(sqlFile); err == nil {
		data, err := os.ReadFile(sqlFile)
		if err != nil {
			log.Fatal("Error reading sql file: ", err)
		}
		json.Unmarshal(data, &sqlSlice)
	}
}

func saveSql() {
	data, err := json.MarshalIndent(sqlSlice, "", " ")
	if err != nil {
		log.Fatal("Error Saving SQL: ", err)
	}
	os.WriteFile(sqlFile, data, 0644)
}

func refreshSql(sqlList *tview.TextView) {
	sqlList.Clear()
	if len(sqlSlice) == 0 {
		fmt.Fprintln(sqlList, "No Items in list")
	} else {
		for i, item := range sqlSlice {
			fmt.Fprintf(sqlList, "[%d] %s (Statement: %s)\n", i+1, item.Name, item.Statement)
		}
	}
}

//TODO: Reimplement deleting from the sql file
// func deleteItem(index int) {
// 	if index < 0 || index >= len(sqlSlice) {
// 		fmt.Println("Invalid item index!")
// 		return
// 	}
// 	sqlSlice = append(sqlSlice[:index], sqlSlice[index+1:]...)
// 	saveSql()
// }

func createWelcomePage() *tview.Flex {
	lines := strings.Split(logo, "\n")
	logoWidth := 0
	logoHeight := len(lines)
	for _, line := range lines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}

	logoBox := tview.NewTextView().
		SetTextColor(tcell.NewRGBColor(57, 255, 20))
	fmt.Fprint(logoBox, logo)

	frame := tview.NewFrame(tview.NewBox()).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(version, true, tview.AlignCenter, tcell.ColorOrange).
		AddText("", true, tview.AlignCenter, tcell.ColorWhite).
		AddText(subtitle, true, tview.AlignCenter, tcell.ColorWhite).
		AddText("", true, tview.AlignCenter, tcell.ColorWhite).
		AddText(navigation, true, tview.AlignCenter, tcell.ColorDarkMagenta)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 7, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(logoBox, logoWidth, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, true).
		AddItem(frame, 0, 10, false)

	return flex
}

func createQueryBuilder(sqlList *tview.TextView) *tview.Flex {
	fromTableInput := tview.NewInputField().SetLabel("From Table: ")
	fromSqlInput := tview.NewInputField()
	toTableInput := tview.NewInputField().SetLabel("To Table: ")
	toSqlInput := tview.NewInputField()

	joinInput := tview.NewInputField().SetLabel("Join Side 1")
	joinInput1 := tview.NewInputField().SetLabel("Join Side 2")

	fromTable := tview.NewForm().
		AddFormItem(fromTableInput)

	toTable := tview.NewForm().
		AddFormItem(toTableInput)

	fromForm := tview.NewForm().
		AddFormItem(fromSqlInput)

	toForm := tview.NewForm().
		AddFormItem(toSqlInput)

	joinForm := tview.NewForm().
		AddFormItem(joinInput).
		AddFormItem(joinInput1)

	buttonForm := tview.NewForm().
		AddButton("Add Item", func() {
			var builder strings.Builder
			fromTable := fromTableInput.GetText()
			fromInput := fromSqlInput.GetText()
			toTable := toTableInput.GetText()
			toInput := toSqlInput.GetText()
			joinInputLeft := joinInput.GetText()
			joinInputRight := joinInput1.GetText()
			if fromInput != "" && toInput != "" {

				if joinInputLeft != "" && joinInputRight != "" {
					builder.WriteString(fmt.Sprintf("INSERT INTO %s", toTable))
					builder.WriteString(fmt.Sprintf(" (%s)", toInput))
					builder.WriteString(fmt.Sprintf(" SELECT %s", fromInput))
					builder.WriteString(fmt.Sprintf(" FROM %s", fromTable))
					builder.WriteString(fmt.Sprintf(" INNER JOIN %s ON %s = %s", joinInputLeft, joinInputLeft, joinInputRight))
				} else {
					builder.WriteString(fmt.Sprintf("INSERT INTO %s", toTable))
					builder.WriteString(fmt.Sprintf(" (%s)", toInput))
					builder.WriteString(fmt.Sprintf(" SELECT %s", fromInput))
					builder.WriteString(fmt.Sprintf(" FROM %s", fromTable))
				}

				sqlSlice = append(sqlSlice, Sql{Name: "Testing", Statement: builder.String()})
				saveSql()
				refreshSql(sqlList)
				fromSqlInput.SetText("")
				toSqlInput.SetText("")
				fromTableInput.SetText("")
				toTableInput.SetText("")
				joinInput.SetText("")
				joinInput1.SetText("")
				builder.Reset()
			}
		})

	fromTable.SetBorder(true).SetTitle("From Table Name").SetTitleAlign(tview.AlignCenter)
	toTable.SetBorder(true).SetTitle("To Table Name").SetTitleAlign(tview.AlignCenter)
	fromForm.SetBorder(true).SetTitle("SQL From").SetTitleAlign(tview.AlignCenter)
	toForm.SetBorder(true).SetTitle("SQL To").SetTitleAlign(tview.AlignCenter)
	joinForm.SetBorder(true).SetTitle("Joins").SetTitleAlign(tview.AlignCenter)
	buttonForm.SetBorder(true)

	tableNameRow := tview.NewFlex().
		AddItem(fromTable, 0, 1, true).
		AddItem(toTable, 0, 1, true)

	topRow := tview.NewFlex().
		AddItem(fromForm, 0, 1, true).
		AddItem(toForm, 0, 1, true)

	bottomRow := tview.NewFlex().
		AddItem(joinForm, 0, 3, true).
		AddItem(buttonForm, 0, 1, true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tableNameRow, 0, 1, true).
		AddItem(topRow, 0, 2, true).
		AddItem(bottomRow, 0, 1, true)

	return layout
}

func createNavBar(pages *tview.Pages) *tview.TextView {
	navBar := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			if len(added) == 0 {
				return
			}
			pages.SwitchToPage(added[0])
		})
	return navBar
}

func populateNavBar(pageNames []string, navBar *tview.TextView) {
	for i := 0; i < len(pageNames); i++ {
		fmt.Fprintf(navBar, `["%d"][#00FFFF]%s[white][""]  `, i, pageNames[i])
	}
	navBar.Highlight("0")
}

func main() {
	app := tview.NewApplication()
	app.EnableMouse(true)
	pages := tview.NewPages()
	currentPage := 0

	loadSql()
	sqlList := tview.NewTextView().SetWordWrap(true)

	pages.AddPage("0", createWelcomePage(), true, true)

	pages.AddPage("1", createQueryBuilder(sqlList), true, false)

	pages.AddPage("2", sqlList, true, false)

	refreshSql(sqlList)

	pages.SetBorder(true).SetTitle("Quick Convert").SetTitleAlign(tview.AlignCenter)

	navBar := createNavBar(pages)
	populateNavBar(pageNames, navBar)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(navBar, 1, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			currentPage = (currentPage + 1) % pageCount
			pages.SwitchToPage(fmt.Sprintf("%d", currentPage))
			navBar.Highlight(fmt.Sprintf("%d", currentPage)).ScrollToHighlight()
			return nil
		} else if event.Key() == tcell.KeyCtrlP {
			currentPage = (currentPage - 1 + pageCount) % pageCount
			pages.SwitchToPage(fmt.Sprintf("%d", currentPage))
			navBar.Highlight(fmt.Sprintf("%d", currentPage)).ScrollToHighlight()
			return nil
		}
		return event
	})

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
