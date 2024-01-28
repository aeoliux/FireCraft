package window

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/zapomnij/firecraft/pkg/auth"
	"strings"
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
	authButton   *widgets.QPushButton
	unauthButton *widgets.QPushButton
	logger       *widgets.QTextEdit

	AccessToken       string
	Uuid              string
	HaveBoughtTheGame bool
	Authed            bool
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
	ms.messages[3] = widgets.NewQLabel2("4. Click Authenticate", ms.widget, 0)
	for _, j := range ms.messages[:3] {
		ms.contLayout.AddWidget(j, 0, core.Qt__AlignTop)
	}

	ms.RedirectLink = widgets.NewQLineEdit(ms.widget)
	ms.RedirectLink.ConnectTextChanged(func(text string) {
		ms.parentWindow.usernameTv.SetReadOnly(text != "")
	})
	ms.contLayout.AddWidget(ms.RedirectLink, 0, core.Qt__AlignTop)

	ms.contLayout.AddWidget(ms.messages[3], 0, core.Qt__AlignTop)

	ms.logger = widgets.NewQTextEdit(ms.container)
	ms.logger.SetReadOnly(true)
	ms.contLayout.AddWidget(ms.logger, 0, core.Qt__AlignBottom)

	ms.authButton = widgets.NewQPushButton2("Authenticate", ms.container)
	ms.authButton.ConnectClicked(func(checked bool) {
		ms.logger.Clear()
		if ms.RedirectLink.Text() == "" {
			ms.logger.SetText("Follow the instructions\n")
			return
		}

		ms.parentWindow.playBt.SetEnabled(false)
		ms.authButton.SetEnabled(false)
		ch := make(chan string)
		go func() {
			for {
				msg := <-ch
				ms.logger.MoveCursor(gui.QTextCursor__End, gui.QTextCursor__MoveAnchor)
				ms.logger.InsertPlainText(msg)
				if strings.HasPrefix(msg, "error") {
					ms.parentWindow.playBt.SetEnabled(true)
					ms.authButton.SetEnabled(true)
					return
				} else if msg == "authenticator: success\n" {
					ms.parentWindow.playBt.SetEnabled(true)
					return
				}
			}
		}()

		go ms.Auth(ch)
	})
	ms.contLayout.AddWidget(ms.authButton, 0, core.Qt__AlignBottom)

	ms.unauthButton = widgets.NewQPushButton2("Unauthenticate", ms.container)
	ms.unauthButton.ConnectClicked(ms.Unauth)
	ms.unauthButton.SetEnabled(false)
	ms.contLayout.AddWidget(ms.unauthButton, 0, core.Qt__AlignBottom)

	ms.layout.AddWidget(ms.container, 0, core.Qt__AlignTop)

	return &ms
}

func (ms *MSTab) Auth(ch chan string) {
	au, err := auth.NewAuthentication(ms.RedirectLink.Text())
	if err != nil {
		ch <- "error: " + err.Error() + "\n"
		return
	}

	ch <- "authenticator: authenticating Minecraft\n"
	mc, err := auth.NewMinecraftAuthentication(au.MsAccessToken, au.HtClient)
	if err != nil {
		ch <- "error: " + err.Error() + "\n"
		return
	}

	ch <- "authenticator: fetching profile\n"
	usrProf, err := mc.GetProfile()
	if err != nil {
		ch <- "error: " + err.Error() + "\n"
		return
	}

	ms.HaveBoughtTheGame = mc.OwnsGame()
	ms.AccessToken = mc.MinecraftToken
	ms.Uuid = usrProf.Id
	ms.Authed = true
	ms.parentWindow.usernameTv.SetText(usrProf.Name)

	ms.unauthButton.SetEnabled(true)
	ch <- "authenticator: success\n"
}

func (ms *MSTab) Unauth(_ bool) {
	ms.RedirectLink.Clear()
	ms.parentWindow.usernameTv.Clear()
	ms.Authed = false
	ms.unauthButton.SetEnabled(false)
	ms.authButton.SetEnabled(true)
	ms.logger.Clear()
	ms.logger.SetText("authenticator: unauthenticated")
}
