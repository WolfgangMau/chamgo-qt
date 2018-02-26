## chamgo-qt
QT based GUI for Chameleon-RevE-Rebooted &amp; Chameleon RevG - written in Golang

what works in general (on RevE-Rebooted & RevG):
- USB-Device detection
- Serial connection
- display /edit  Slot-Config
- Serial Terminal
- upload / download dump
- get & decode detection-nonces for the use with mfkey32v2 (RevE-Rebooted)
- logdownload  (RevG - untested)
- implementation of crc16 / BCC 
- display RSSI
- integration of mfkey32v2 (RevE-Rebooted - workaround with binary)
- diff/edit Dump-/Slot-Data

what is missing:
- logmode=live (RevG)

to install/create the qt-bindings follow the instuctions in the wiki: https://github.com/therecipe/qt/wiki

there are also some pre-compiled binaries in the ['release-section'](https://github.com/WolfgangMau/chamgo-qt/releases)

### Screenshots
#### Serial-Terminal
![Serial-Terminal](https://github.com/WolfgangMau/chamgo-qt/blob/master/screenshots/Serial-Terminal.png)

#### Tag-Editor
![Tag-Editor](https://github.com/WolfgangMau/chamgo-qt/blob/master/screenshots/Tag-Editor.png)

#### Data-Diff
![Tag-Editor](https://github.com/WolfgangMau/chamgo-qt/blob/master/screenshots/Data-Diff.png)
