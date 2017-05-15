# Aether

Aether 1 last version: 1.2.3 (deprecated)   
Aether 2 last version: 2.0.0 (currently non-functional, in development)

*NOTICE: Aether 1 is currently deprecated and not supported. Aether 2 is currently in development.*

Aether is a free app that you use to read, write in, and create community moderated, distributed, and anonymous forums, an “anonymous reddit without servers.” — [The Verge](http://www.theverge.com/2013/11/27/5150758/aether-aims-to-be-a-reddit-for-the-privacy-conscious)

More information at the [website.](http://www.getaether.net) 

Aether is available for OS X, Windows and Linux. If you're just looking to try it out, [download the app directly](http://getaether.net/download). All provided binaries are signed, signatures are [here](https://github.com/nehbit/aether-public/releases). 

## Aether 2

*(Not yet working, 80% complete)*

Aether 2 is a Go application built from ground up to provide a scalable architecture that can accept third party clients. The front-end is Electron. The backend and frontend talk to each other over gRPC.

The backend is reasonably complete, however, the frontend and the backend - frontend glue code is still in progress. 

## Aether 1
*(Deprecated, only here for documentation purposes)*

###  How to build

If you choose to not trust the provided binaries, you can build it yourself, instructions are below. 

You need a working PyQt that should be linked to your system Qt installation. That's impossible to automate, so please check the respective documentation on how to do that. In short, you should install [Qt](http://qt-project.org/downloads) first, and then follow [these instructions](http://pyqt.sourceforge.net/Docs/PyQt5/introduction.html) for PyQt.

Afterwards, install the requirements.txt via pip, and do `python main.py` on the project folder.

If you want to package the app, the PyInstaller spec files for Windows, Mac and OS X are provided. Fill in the required folder names on spec files, and PyInstaller should handle the rest. 

**Note:** Unless you're building for OS X, you do not need PyObjC dependencies in requirements.txt. Feel free to remove them.
