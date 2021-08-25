@echo off
tasklist /nh /fi "imagename eq go-mosbrew-scales.exe" | find /i "go-mosbrew-scales.exe" > nul || (Start /I "" "C:\Auto-Control\go-mosbrew-scales.exe")
exit
