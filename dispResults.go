// dispResults.go

/*
	Source file auto-generated on Wed, 23 Oct 2019 18:40:49 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019 H.F.M - Functions & Library Manager
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
	// gltsbh "github.com/hfmrow/genLib/tools/bench"
)

/*
	Main TreeStore (results)
*/

// displayTreeStore: Fill TreeViewFound with found results
func displayTreeStore(in []toDispTreeStore) (err error) {
	var iter *gtk.TreeIter

	if len(in) > 0 { // Detach & clear listStore
		tvsTreeSearch.StoreDetach()
		tvsTreeSearch.TreeStore.Clear()

		// Add parents
		for _, row := range in {
			if iter, err = tvsTreeSearch.AddRow(nil, row.Name, row.Type, row.Exported, row.Path, row.Score, row.Idx); err != nil {
				tvsTreeSearch.StoreAttach()
				DlgErr("displayTreeStore:AddParents", err)
				return
			}
			// Add childs if exists (structures' methods)
			for _, subRow := range row.Methods {
				if _, err = tvsTreeSearch.AddRow(iter, subRow.Name, subRow.Type, subRow.Exported, subRow.Path, subRow.Score, subRow.Idx); err != nil {
					tvsTreeSearch.StoreAttach()
					DlgErr("displayTreeStore:AddChilds", err)
					return
				}
			}
		}
		// Attach listStore
		tvsTreeSearch.StoreAttach()
	}
	// Update statusbar
	if tvsTreeSearch.CountRows() > 0 {
		updateStatusBar()
	} else {
		updateStatusBar(fmt.Sprintf(sts["noResult"]+"\"%s\"", GetEntryText(mainObjects.EntrySearchFor)))
	}
	return
}

/*
	Popup window
*/

// popupTreeview: Display content as TextView
func popupSourceView(index int) {

	var err error

	if svs == nil {

		svs, err = SourceViewStructNew(mainObjects.WindowSource, mainObjects.Source, mainObjects.SourceMap)
		DlgErr("popupSourceView:SourceViewStructNew", err)
		svs.View.SetEditable(false)

		// Handling "populate-popup" signal to add some personal entries
		svs.View.Connect("populate-popup", popupTextViewPopulateMenu)

		// TODO Think to use search in preview window
		// Make a tag to indicate found element (when HighlightFound not checked)
		tag := make(map[string]interface{})
		tag["background"] = "#ABF6FF"
		markFound = svs.Buffer.CreateTag("markFound", tag)

		// Language & style, add a personal version for Golang (directory content)
		svs.UserStylePath = mainOptions.HighlightUserDefined
		svs.UserLanguagePath = mainOptions.HighlightUserDefined

		// Setting Language and style scheme
		if currentLanguage := svs.SetLanguage(mainOptions.DefaulLanguage); currentLanguage != nil {
			if currentStyleScheme := svs.SetStyleScheme(mainOptions.DefaultStyle); currentStyleScheme != nil {

				// Fill comboboxes with languages and styles
				for _, id := range svs.LanguageIds {
					mainObjects.ComboboxSourceLanguage.AppendText(id)
				}
				for _, id := range svs.StyleShemeIds {
					mainObjects.ComboboxSourceStyle.AppendText(id)
				}

				// Just indicate id must be set as first model column.
				mainObjects.ComboboxSourceLanguage.SetIDColumn(0)
				mainObjects.ComboboxSourceStyle.SetIDColumn(0)

				// Set ComboBox current values display.
				mainObjects.ComboboxSourceLanguage.SetActiveID(mainOptions.DefaulLanguage)
				mainObjects.ComboboxSourceStyle.SetActiveID(mainOptions.DefaultStyle)
			}
		}

		mainObjects.WindowSource.Resize(mainOptions.SourceWinWidth, mainOptions.SourceWinHeight)
		mainObjects.WindowSource.Move(mainOptions.SourceWinPosX, mainOptions.SourceWinPosY)
		mainObjects.PanedSource.SetPosition(mainOptions.MainWinWidth - mainOptions.PanedWidth)
	}

	indexCurrText = index // needed to always get the last selected choice
	dispTextView(index)
}

// dispTreeStore:
func dispTextView(index int) {
	var err error

	descr, ok := declIdexes.GetDescr(index)
	if ok {
		if err = svs.LoadSource(filepath.Join(desc.RootLibs, descr.File)); err == nil {

			mainObjects.WindowSource.SetTitle(descr.File)

			svs.TxtBgCol = mainOptions.TxtBgCol
			svs.TxtFgCol = mainOptions.TxtFgCol

			svs.ColorBgRangeSet = mainOptions.DefRangeCol
			svs.SelBgCol = mainOptions.SelBgCol

			svs.ColorBgRange(descr.LineStart+1, descr.LineEnd+1)
			svs.SelectRange(descr.LineStart+1, descr.LineEnd+1)

			svs.RunAfterEvents(func() {
				mainObjects.WindowSource.SetKeepAbove(true)
				svs.ScrollToLine(descr.LineStart)
				mainObjects.WindowSource.ShowAll()
				mainObjects.WindowSource.SetKeepAbove(false)
			})

		}
		DlgErr("ReadFile", err)
		// tvn.TextView.GrabFocus()
		return
	}
	DlgErr("Description not found", fmt.Errorf("#%d", index))
}
