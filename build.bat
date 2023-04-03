go build -ldflags "-H=windowsgui" -o SpeakerKeeper_windows.exe

setlocal

set FILENAMES=SpeakerKeeper_windows.exe,sound.mp3
set ZIPNAME=SpeakerKeeper_windows.zip

powershell -Command "Compress-Archive -Path %FILENAMES% -DestinationPath %ZIPNAME%"

endlocal