package window

import (
	"log"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/zapomnij/firecraft/pkg/downloader"
)

var (
	lpf *LauncherProfiles
	vm  *downloader.VersionManifest
)

type FWindow struct {
	Window *widgets.QWidget

	container *widgets.QWidget
	layout    *widgets.QGridLayout

	logger   *widgets.QTextEdit
	notebook *widgets.QTabWidget

	microsoft *MicrosoftAuthTab

	bottombar       *widgets.QWidget
	bottombarLayout *widgets.QHBoxLayout
	playBt          *widgets.QPushButton

	userBox    *widgets.QWidget
	userLay    *widgets.QGridLayout
	usernameTv *widgets.QLineEdit

	profilesBox      *widgets.QWidget
	profilesLay      *widgets.QGridLayout
	profilesSelector *widgets.QComboBox
	editProfile      *widgets.QPushButton
}

var oldemail string

func NewFWindow() *FWindow {
	var err error
	lpf, err = loadProfiles()
	if err != nil {
		log.Fatalln(err)
	}
	oldemail = lpf.AuthenticationDatabase.Email

	vm, err = downloader.GetVersionManifest()
	if err != nil {
		log.Fatalln(err)
	}

	var this = FWindow{}
	this.Window = widgets.NewQWidget(nil, 0)
	this.Window.SetWindowTitle(core.QCoreApplication_ApplicationName())

	this.container = widgets.NewQWidget(this.Window, 0)
	this.layout = widgets.NewQGridLayout(this.container)
	this.layout.SetSpacing(0)
	this.layout.SetContentsMargins(0, 0, 0, 0)

	this.notebook = widgets.NewQTabWidget(this.container)
	this.logger = widgets.NewQTextEdit(this.notebook)
	this.logger.SetReadOnly(true)
	this.notebook.AddTab(this.logger, "Launcher logs")
	this.layout.AddWidget(this.notebook)

	this.bottombar = widgets.NewQWidget(this.container, 0)
	this.bottombar.SetFixedHeight(70)
	this.bottombarLayout = widgets.NewQHBoxLayout()
	this.bottombarLayout.SetContentsMargins(0, 0, 0, 0)
	this.bottombar.SetLayout(this.bottombarLayout)

	this.profilesBox = widgets.NewQWidget(this.bottombar, 0)
	this.profilesLay = widgets.NewQGridLayout(this.profilesBox)
	this.profilesBox.SetLayout(this.profilesLay)
	this.editProfile = widgets.NewQPushButton2("Edit profile", this.profilesBox)
	this.editProfile.ConnectClicked(this.editProfileHandle)
	this.profilesSelector = widgets.NewQComboBox(this.profilesBox)
	for k := range lpf.Profiles {
		this.profilesSelector.AddItem(k, core.NewQVariant())
	}
	this.profilesSelector.AddItem("New profile", core.NewQVariant())

	this.profilesLay.AddWidget(this.profilesSelector)
	this.profilesLay.AddWidget(this.editProfile)

	this.playBt = widgets.NewQPushButton2("Play", this.bottombar)
	this.playBt.ConnectClicked(func(checked bool) { go this.Launch() })
	this.playBt.SetFixedHeight(60)
	this.playBt.SetFixedWidth(300)

	this.userBox = widgets.NewQWidget(this.bottombar, 0)
	this.userLay = widgets.NewQGridLayout(this.userBox)
	this.userLay.AddWidget(widgets.NewQLabel2("Username", this.userBox, 0))
	this.usernameTv = widgets.NewQLineEdit2(lpf.AuthenticationDatabase.Username, this.userBox)
	this.usernameTv.ConnectTextChanged(this.saveUsername)
	this.userLay.AddWidget(this.usernameTv)
	this.userBox.SetLayout(this.userLay)

	this.bottombarLayout.AddWidget(this.profilesBox, 0, core.Qt__AlignLeft)
	this.bottombarLayout.AddWidget(this.playBt, 0, core.Qt__AlignHCenter)
	this.bottombarLayout.AddWidget(this.userBox, 0, core.Qt__AlignRight)

	this.microsoft = NewMicrosoftAuthTab(this.notebook, &this)
	this.notebook.AddTab(this.microsoft.Main, "Microsoft Authentication")

	this.layout.AddWidget(this.bottombar)

	this.container.SetLayout(this.layout)
	this.Window.SetLayout(this.layout)

	return &this
}

func (fw *FWindow) editProfileHandle(checked bool) {
	epw := NewEditProfileWindow(fw)
	epw.Window.Resize2(300, 200)
	epw.Window.Show()
}

func (fw FWindow) saveUsername(text string) {
	lpf.AuthenticationDatabase.Username = text
	lpf.Save()
}

func (fw *FWindow) reloadProfileSelector() {
	fw.profilesSelector.Clear()

	for k := range lpf.Profiles {
		fw.profilesSelector.AddItem(k, core.NewQVariant())
	}

	fw.profilesSelector.AddItem("New profile", core.NewQVariant())
}

func (fw *FWindow) appendToLog(msg string) {
	fw.logger.MoveCursor(gui.QTextCursor__End, gui.QTextCursor__MoveAnchor)
	fw.logger.InsertPlainText(msg)
}
