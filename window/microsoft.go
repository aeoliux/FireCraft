package window

import "github.com/therecipe/qt/widgets"

type MicrosoftAuthTab struct {
	Main   *widgets.QWidget
	Window *FWindow

	Layout *widgets.QFormLayout
	Email  *widgets.QLineEdit
	Passwd *widgets.QLineEdit
}

func NewMicrosoftAuthTab(parent widgets.QWidget_ITF, window *FWindow) *MicrosoftAuthTab {
	m := MicrosoftAuthTab{}
	m.Main = widgets.NewQWidget(parent, 0)
	m.Window = window

	m.Layout = widgets.NewQFormLayout(m.Main)
	m.Main.SetLayout(m.Layout)

	m.Email = widgets.NewQLineEdit(parent)
	m.Email.ConnectTextChanged(m.disableUsername)
	m.Email.SetText(lpf.AuthenticationDatabase.Email)
	m.Passwd = widgets.NewQLineEdit(parent)
	m.Passwd.SetEchoMode(widgets.QLineEdit__Password)

	m.Layout.AddRow3("Email: ", m.Email)
	m.Layout.AddRow3("Password: ", m.Passwd)

	return &m
}

func (m *MicrosoftAuthTab) disableUsername(text string) {
	if text == "" {
		m.Window.usernameTv.SetReadOnly(false)
	} else {
		m.Window.usernameTv.SetReadOnly(true)
	}

	lpf.AuthenticationDatabase.Email = text
	lpf.Save()
}
