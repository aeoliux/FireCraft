package window

import (
	"fmt"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type EditProfileWindow struct {
	Window     *widgets.QWidget
	layout     *widgets.QFormLayout
	nameEn     *widgets.QLineEdit
	verList    *widgets.QComboBox
	javaArgsEn *widgets.QLineEdit
	javaBinEn  *widgets.QLineEdit
	gameDirEn  *widgets.QLineEdit

	bottom    *widgets.QWidget
	bottomLay *widgets.QHBoxLayout
	okBt      *widgets.QPushButton
	cancelBt  *widgets.QPushButton
	deleteBt  *widgets.QPushButton
	errLabel  *widgets.QLabel

	parent *FWindow
}

func NewEditProfileWindow(parent *FWindow) *EditProfileWindow {
	this := EditProfileWindow{}
	this.parent = parent

	this.Window = widgets.NewQWidget(nil, 0)
	this.Window.SetWindowTitle("Profile editor")

	this.layout = widgets.NewQFormLayout(this.Window)

	this.nameEn = widgets.NewQLineEdit(this.Window)
	this.layout.AddRow3("Name: ", this.nameEn)

	this.gameDirEn = widgets.NewQLineEdit(this.Window)
	this.layout.AddRow3("Game directory: ", this.gameDirEn)

	this.verList = widgets.NewQComboBox(this.Window)
	for _, v := range vm.Versions {
		this.verList.AddItem(fmt.Sprintf("%s %s", v.Type, v.Id), core.NewQVariant())
	}
	this.layout.AddRow3("Version: ", this.verList)

	this.javaArgsEn = widgets.NewQLineEdit(this.Window)
	this.layout.AddRow3("Java args: ", this.javaArgsEn)

	this.javaBinEn = widgets.NewQLineEdit(this.Window)
	this.layout.AddRow3("Java binary path: ", this.javaBinEn)

	this.bottom = widgets.NewQWidget(this.Window, 0)
	this.bottomLay = widgets.NewQHBoxLayout()
	this.bottom.SetLayout(this.bottomLay)
	this.okBt = widgets.NewQPushButton2("OK", this.bottom)
	this.okBt.ConnectClicked(func(checked bool) {
		if this.saveProfile() {
			lpf.Save()
			this.Window.Destroy(true, true)
			this.Window.DeleteLater()
		} else {
			this.errLabel.SetVisible(true)
		}
	})
	this.deleteBt = widgets.NewQPushButton2("Delete profile", this.bottom)
	this.deleteBt.ConnectClicked(func(checked bool) {
		selected := this.nameEn.Text()
		if selected != "New profile" && selected != "" {
			lpf.deleteProfile(selected)
			lpf.Save()
			this.parent.reloadProfileSelector("")
		}

		this.Window.Destroy(true, true)
		this.Window.DeleteLater()
	})
	this.cancelBt = widgets.NewQPushButton2("Cancel", this.bottom)
	this.cancelBt.ConnectClicked(func(checked bool) { this.Window.Destroy(true, true); this.Window.DeleteLater() })
	this.bottomLay.AddWidget(this.okBt, 0, core.Qt__AlignRight)
	this.bottomLay.AddWidget(this.deleteBt, 0, core.Qt__AlignRight)
	this.bottomLay.AddWidget(this.cancelBt, 0, core.Qt__AlignRight)

	this.errLabel = widgets.NewQLabel2("Invalid profile", this.Window, 0)
	this.errLabel.SetVisible(false)
	this.layout.AddRow(this.errLabel, this.bottom)

	this.loadProfile()

	this.Window.SetLayout(this.layout)
	return &this
}

func (epw *EditProfileWindow) loadProfile() {
	selected := epw.parent.profilesSelector.CurrentText()
	if selected == "New profile" {
		epw.nameEn.SetText("New profile")
		epw.verList.SetCurrentIndex(FindVerById(vm.Latest.Release))
		epw.javaArgsEn.SetText("-Xmx2048M")
	} else {
		profile, ok := lpf.Profiles[selected]
		if !ok {
			epw.nameEn.SetText(selected)
			epw.verList.SetCurrentIndex(FindVerById(vm.Latest.Release))
			epw.javaArgsEn.SetText("-Xmx2048M")
		} else {
			epw.nameEn.SetText(selected)
			epw.javaBinEn.SetText(profile.JavaBin)
			epw.javaArgsEn.SetText(profile.JavaArgs)
			epw.verList.SetCurrentIndex(FindVerById(profile.LastVersionId))
			epw.gameDirEn.SetText(profile.GameDir)
		}
	}
}

func (epw *EditProfileWindow) saveProfile() bool {
	name := epw.nameEn.Text()

	if name == "New profile" {
		return false
	}

	for k := range lpf.Profiles {
		if k == name {
			epw.verList.RemoveItem(findPosition(name))
		}
	}

	lpf.Profiles[name] = LProfile{
		name,
		"custom",
		strings.Split(epw.verList.CurrentText(), " ")[1],
		epw.gameDirEn.Text(),
		epw.javaBinEn.Text(),
		epw.javaArgsEn.Text(),
	}

	epw.parent.reloadProfileSelector(name)

	return true
}

func findPosition(key string) int {
	ind := 0
	for k := range lpf.Profiles {
		if k == key {
			return ind
		}

		ind++
	}

	return 0
}

func FindVerById(ver string) int {
	for i, j := range vm.Versions {
		if j.Id == ver {
			return i
		}
	}

	return 0
}
