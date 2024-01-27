# FireCraft
An open source Minecraft Launcher written in Go which uses Qt for GUI. 

### Supported operating systems
FireCraft currently supports Linux and Windows... and macOS, but without Java finder (also working on that).

### Versions that don't work
- anything newer than 1.19, including 1.19 (due to changed management over native libraries, working on that, please be patient)
- 1.8.x - unknown issue unfortunately

### Launching Minecraft
Launcher installation on Linux:
```sh
make linux-install
```

Compilation on Windows (it should be done under MSYS2 or cygwin):
```sh
make clean firecraft
./firecraft.exe # launcher's .exe
```

### Minecraft Premium support
Works!