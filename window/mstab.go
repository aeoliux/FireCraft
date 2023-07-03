package window

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type MSTab struct {
	widget       *widgets.QWidget
	layout       *widgets.QVBoxLayout
	parent       widgets.QWidget_ITF
	parentWindow *FWindow

	container    *widgets.QWidget
	contLayout   *widgets.QVBoxLayout
	messages     [4]*widgets.QLabel
	RedirectLink *widgets.QLineEdit
}

func NewMSTab(parent widgets.QWidget_ITF, parentWindow *FWindow) *MSTab {
	ms := MSTab{}
	ms.parent = parent
	ms.parentWindow = parentWindow
	ms.widget = widgets.NewQWidget(parent, 0)
	ms.layout = widgets.NewQVBoxLayout()
	ms.widget.SetLayout(ms.layout)

	ms.container = widgets.NewQWidget(ms.widget, 0)
	ms.contLayout = widgets.NewQVBoxLayout()
	ms.container.SetLayout(ms.contLayout)

	ms.messages[0] = widgets.NewQLabel2("1. Go to <a href=\"https://login.live.com/oauth20_authorize.srf?client_id=000000004C12AE6F&redirect_uri=https://login.live.com/oauth20_desktop.srf&scope=service::user.auth.xboxlive.com::MBI_SSL&display=touch&response_type=token&locale=en\">this link</a>", ms.widget, 0)
	ms.messages[0].SetOpenExternalLinks(true)
	ms.messages[1] = widgets.NewQLabel2("2. Log into MS", ms.widget, 0)
	ms.messages[2] = widgets.NewQLabel2("3. When you see a blank page, copy link from the web browser and paste it below", ms.widget, 0)
	ms.messages[3] = widgets.NewQLabel2("4. Click Play", ms.widget, 0)
	for _, j := range ms.messages[:3] {
		ms.contLayout.AddWidget(j, 0, core.Qt__AlignTop)
	}

	ms.RedirectLink = widgets.NewQLineEdit(ms.widget)
	ms.RedirectLink.ConnectTextChanged(ms.updateUsername)
	ms.contLayout.AddWidget(ms.RedirectLink, 0, core.Qt__AlignTop)

	ms.contLayout.AddWidget(ms.messages[3], 0, core.Qt__AlignTop)

	ms.layout.AddWidget(ms.container, 0, core.Qt__AlignTop)

	return &ms
}

func (ms *MSTab) updateUsername(text string) {
	if text != "" {
		ms.parentWindow.usernameTv.SetText("")
		ms.parentWindow.usernameTv.SetEnabled(false)
		lpf.Save()
	} else {
		ms.parentWindow.usernameTv.SetEnabled(true)
	}
}
