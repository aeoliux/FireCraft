# Disclaimer
Because of Qt bindings not being updated for 4 years and not working on Silicon macs, project is discontinued.

# FireCraft
An open source Minecraft Launcher written in Go which uses Qt5 for GUI. 

### Supported operating systems
FireCraft currently supports Linux, macOS and Windows.

### Versions that don't work
- anything newer than 1.19, including 1.19 (due to changed management over native libraries, working on that, please be patient)
- 1.8.x - unknown issue unfortunately

### Launching Minecraft
To set up Qt for Go, follow instructions at https://github.com/therecipe/qt/wiki/Installation

#### Linux and macOS:
If you're using homebrew Qt on macOS, first do `export QT_HOMEBREW=true`.

To deploy app:
```sh
make clean deploy # after that, you should see something like ./deploy/darwin/firecraft.app or ./deploy/linux/firecraft
make linux-install # can be used on Linux to install launcher with its desktop files
```

#### Windows (under MSYS2):
Unfortunately deployment doesn't work.
```sh
make clean firecraft
./firecraft.exe # launcher's .exe
```

### Minecraft Premium support
Works for Microsoft accounts, see MS Authentication tab in the launcher.

### Libs which were used
- Qt5 bindings for Go: https://github.com/therecipe/qt
