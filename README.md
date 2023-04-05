# SpeakerKeeper

SpeakerKeeper is a simple utility that prevents your computer's speaker from going to sleep after a period of inactivity.
It works by playing a nearly silent audio file at specified intervals, which keeps the speaker active.

## How it works

After launching the app, you can select the audio output device and the time interval for playing the audio file.
The app runs silently in the background and automatically places an icon in the system tray while running. 

From the menu, you can reset the configuration or stop playing the audio file.
The file "sound.mp3" can be replaced by any MP3 file you want the application to play.
Simply make sure that the file name remains as "sound.mp3"



## Notes

To make the application run on system startup, you can place a shortcut to the .exe in:

C:\Users\\"USER"\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup


This is potentially compatible with Linux/MacOS but I need to set up the compilers for them since this uses Cgo in the dependencies.
Maybe someone who knows more than me can help set that up. 
