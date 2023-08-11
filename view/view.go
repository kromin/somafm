package view

import (
	"github.com/nicarl/somafm/radioChannels"
	"github.com/nicarl/somafm/state"
	"github.com/nuttech/bell"
	"github.com/rivo/tview"

	tcell "github.com/gdamore/tcell/v2"
)

func getChannelList(appState *state.AppState, channelDetails *tview.Flex, player *tview.List) *tview.List {
	channelList := tview.NewList()
	channelList.SetBorder(true).SetTitle("Channels")
	channelList.ShowSecondaryText(false)
	channelList.SetBorderColor(tcell.ColorGreenYellow)
	channelList.SetBorderPadding(1, 1, 1, 1)

	for _, radioCh := range appState.Channels {
		channelList.AddItem(radioCh.Title, "", 0, func() {
			appState.PlayMusic()
			player.SetItemText(1, "Pause", "")
		})
	}
	channelList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		appState.SelectCh(i)

		if err := bell.Ring("descriptions_update", appState.GetSelectedCh()); err != nil {
			panic(err)
		}
	})
	return channelList
}

func getChannelDetails(appState *state.AppState) *tview.Flex {
	desc_view := tview.NewTextView()
	img_view := tview.NewImage()
	now_view := tview.NewTextView()
	channelDetails := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(img_view, 0, 3, false).
		AddItem(desc_view, 0, 1, false).
		AddItem(now_view, 0, 1, false)

	channelDetails.SetBorder(true).SetTitle("Details")
	channelDetails.SetBorderPadding(1, 0, 0, 0)
	bell.Listen("descriptions_update", func(message bell.Message) {
		radioChan := message.Value.(radioChannels.RadioChan)
		desc_view.SetText(radioChan.GetDetails())
		go img_view.SetImage(radioChan.LargeImage)
	})
	bell.Listen("now_play_update", func(message bell.Message) {
		radioChan := message.Value.(radioChannels.RadioChan)
		now_view.SetText(radioChan.LastPlaying)
	})
	return channelDetails
}

func getPlayer(appState *state.AppState) *tview.List {
	player := tview.NewList()
	player.SetBorder(true)
	player.ShowSecondaryText(false)
	player.SetBorderPadding(1, 1, 1, 1)

	player.AddItem("Volume +", "", 0, func() {
		appState.IncreaseVolume()
	})
	player.AddItem("Play", "", 0, func() {
		if appState.IsPlaying {
			appState.PauseMusic()
			player.SetItemText(1, "Play", "")
		} else {
			appState.PlayMusic()
			player.SetItemText(1, "Pause", "")
		}
	})
	player.AddItem("Volume -", "", 0, func() {
		appState.DecreaseVolume()
	})
	player.SetCurrentItem(1)
	return player
}

func InitApp(appState *state.AppState) {
	app := tview.NewApplication()

	channelDetails := getChannelDetails(appState)
	player := getPlayer(appState)
	channelList := getChannelList(appState, channelDetails, player)

	channelList.SetSelectedFunc(func(_ int, _ string, _ string, _ rune) {
		app.SetFocus(player)
	})
	player.SetDoneFunc(func() {
		app.SetFocus(channelList)
	})

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(channelList, 0, 3, false).
			AddItem(player, 0, 1, false), 0, 1, false).
		AddItem(channelDetails, 0, 1, false)
	flexWithHeader := tview.NewFrame(flex).SetBorders(2, 2, 2, 2, 4, 4).AddText("SomaFM", true, tview.AlignCenter, tcell.ColorRed)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == rune('1') {
			app.SetFocus(channelList)
		} else if event.Rune() == rune('2') {
			app.SetFocus(player)
		} else if event.Rune() == rune('q') || event.Rune() == rune('Q') {
			app.Stop()
		}
		return event
	})

	if err := bell.Ring("descriptions_update", appState.GetSelectedCh()); err != nil {
		panic(err)
	}

	if err := app.SetRoot(flexWithHeader, true).SetFocus(channelList).Run(); err != nil {
		panic(err)
	}
}
