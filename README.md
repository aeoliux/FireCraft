# FireCraft
An open source Minecraft Launcher written in Go which uses Qt for GUI. 

### Supported operating systems
FireCraft currently supports Linux and Windows.

### Supported Minecraft Versions
All versions without 1.8.x should work. 1.8.x needs some work to get it working.

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
Work in progress.