import globals
from globals import PLATFORM, FROZEN, BASEDIR
from PyQt5.QtWidgets import QSystemTrayIcon, QWidget, QMenu, QSplashScreen
from PyQt5 import QtGui
from PyQt5 import QtPrintSupport
from PyQt5 import QtCore
import sys, webbrowser
from PyQt5.QtCore import *
from PyQt5.QtWebKitWidgets import *
from PyQt5.QtWidgets import QApplication, QMainWindow, QFileDialog

from datetime import datetime

class Aether(QApplication):
    def __init__(self):
        QApplication.__init__(self, sys.argv)
        QApplication.setAttribute(QtCore.Qt.AA_UseHighDpiPixmaps)
        QApplication.setQuitOnLastWindowClosed(False)

    def onClickOnDock(self): # This only gets called on Mac.
        ##print('dock clicked')
        self.view.show()
        self.view.raise_() # because the app needs to fire up in the bg to be brought to front.
        globals.raiseAndFocusApp()


class ModifierEventFilter(QObject):
    def eventFilter(self, receiver, event):
        if (event.type() == QEvent.KeyPress):
            if (event.modifiers() == QtCore.Qt.ControlModifier and event.key() == QtCore.Qt.Key_W):
                self.view.hide()
                return True

            elif (event.key() == QtCore.Qt.Key_Control): # Control is Cmd in Mac.
                global lastModifierKeypressDatetime
                lastModifierKeypressDatetime = datetime.now()
                return True

            elif (event.key() == QtCore.Qt.Key_Escape): # Control is Cmd in Mac.
                return True # Esc key does not move. Return true stops the propagation there.
            else:
                #Call Base Class Method to Continue Normal Event Processing
                return super(ModifierEventFilter,self).eventFilter(receiver, event)
        #elif (event.type() == QEvent.ApplicationActivate):
        #    ##print('app activate fired')
        #    view.show()
        #    toggleVisibilityMenuItem.setText('Hide Aether')
        #
        #    return True
        else:
            #Call Base Class Method to Continue Normal Event Processing
            return super(ModifierEventFilter,self).eventFilter(receiver, event)

class AetherMainWindow(QMainWindow):
    def __init__(self, charon, reactor, baseurl, app):
        super(AetherMainWindow, self).__init__()
        self.resize(1148, 680)
        self.app = app
        self.reactor = reactor
        webView = AetherWebView(reactor, baseurl)
        webView.page().mainFrame().addToJavaScriptWindowObject("Charon", charon)
        webView.page().setLinkDelegationPolicy(QWebPage.DelegateAllLinks)
        self.setContextMenuPolicy(Qt.NoContextMenu)
        self.setCentralWidget(webView)
        self.webView = webView
        self.JSContext = webView.page().mainFrame().evaluateJavaScript
        self.Hermes = charon.Hermes
        #self.setWindowFlags(QtCore.Qt.WindowMinimizeButtonHint)
        if PLATFORM == 'WIN':
            self.setWindowIcon(QtGui.QIcon(BASEDIR + 'Assets/aether-white-tray.ico'))

        # from PyQt5.QtWebKit import QWebSettings
        # QWebSettings.globalSettings().setAttribute(QWebSettings.DeveloperExtrasEnabled, True)
        # self.inspector = QWebInspector()
        # self.inspector.resize(1450, 300)
        # self.inspector.move(0,0)
        # self.inspector.setPage(self.webView.page())
        # self.setContextMenuPolicy(Qt.DefaultContextMenu)


    def hideEvent(self, QHideEvent):
        if PLATFORM == 'LNX':
            # Linux GUI behaves differently: since there is no tray icon or a background process, when
            # the app is closed by pressing the X button, it actually has to be closed. It can't keep running.
            self.pp.stop()
            if self.reactor.threadpool is not None:
                self.reactor.threadpool.stop()
            self.close()
            self.reactor.stop()
            self.app.quit()
            sys.exit()

class AetherWebView(QWebView):
    def __init__(self, reactor, baseurl):
        super(AetherWebView, self).__init__()
        self.reactor = reactor

        if FROZEN:
            self.load(QUrl('file:///' + BASEDIR + 'GUI/WebKitApp/index.html'))
        else:
            self.load(QUrl('file:///' + baseurl + 'GUI/WebKitApp/index.html'))

        self.page().action(QWebPage.Reload).setVisible(False)
        self.page().setLinkDelegationPolicy(QWebPage.DelegateAllLinks)

        def linkClick(url):
            webbrowser.open(str(url.toString()))

        self.linkClicked.connect(linkClick)

class SystemTrayIcon(QSystemTrayIcon):
    def __init__(self, basedir, app, parent=None):
        if PLATFORM == 'OSX':
            if app.devicePixelRatio() == 2:
                self.icon = QtGui.QIcon(BASEDIR+'Assets/aether-black-tray.svg')
                self.iconActive = QtGui.QIcon(BASEDIR+'Assets/aether-white-tray.svg')
                self.iconHighlight =  QtGui.QIcon(BASEDIR+'Assets/aether-blue-tray.svg')
            else:
                self.icon = QtGui.QIcon(BASEDIR+'Assets/aether-black-tray.png')
                self.iconActive = QtGui.QIcon(BASEDIR+'Assets/aether-white-tray.png')
                self.iconHighlight =  QtGui.QIcon(BASEDIR+'Assets/aether-blue-tray.png')
        elif PLATFORM == 'WIN':
            self.icon = QtGui.QIcon(BASEDIR+'Assets/aether-white-tray-win.svg')
            self.iconActive = self.icon
            self.iconHighlight = QtGui.QIcon(BASEDIR+'Assets/aether-green-tray-win.svg')
        else:
            pass

        QSystemTrayIcon.__init__(self, self.icon, parent)

        self.menu = QMenu(parent)
        if globals.appIsPaused:
            self.menu.addAction('Paused').setDisabled(True)
        else:
            self.menu.addAction('Online').setDisabled(True)
        self.globalStatusMenuItem = self.menu.actions()[0]

        self.menu.addSeparator() # 1
        self.menu.addAction('You have no replies.').setDisabled(True)
        self.messagesMenuItem = self.menu.actions()[2]
        def goToMessages():
            self.messagesMenuItem.setText('You have no replies.')
            self.messagesMenuItem.setDisabled(True)
            parent.show()
            parent.raise_()
            jsString = \
            ("firstFrameScope = angular.element(document.getElementById('first-frame-contents')).scope();"
             "firstFrameScope.repliesButtonClick();"
             "firstFrameScope.$apply();"
            )
            self.webView.JSContext(jsString)
            # reach out to jscontext and
            # Here, I need to call qtwebkit and tell it to open messages.
        self.messagesMenuItem.triggered.connect(goToMessages)
        self.menu.addSeparator() # 3
        if globals.appIsPaused:
            self.menu.addAction('Resume')
        else:
            self.menu.addAction('Pause')
        self.togglePauseMenuItem = self.menu.actions()[4]
        def togglePause():
            if globals.appIsPaused:
                globals.appIsPaused = False
                self.togglePauseMenuItem.setText('Pause')
                self.globalStatusMenuItem.setText('Online')
            else:
                globals.appIsPaused = True
                self.togglePauseMenuItem.setText('Resume')
                self.globalStatusMenuItem.setText('Paused')
        self.togglePauseMenuItem.triggered.connect(togglePause)
        self.menu.addAction('Show Aether')
        self.toggleVisibilityMenuItem = self.menu.actions()[5]
        def makeVisible():
            parent.show()
            parent.raise_()
            if PLATFORM == 'OSX':
                globals.raiseAndFocusApp()

        self.toggleVisibilityMenuItem.triggered.connect(makeVisible)

        self.menu.addAction('Email the developer')
        self.emailDevMenuItem = self.menu.actions()[6]
        def emailDev():
            mailInitialiser = \
                QUrl('mailto:burak@nehbit.net'
                     '?subject=Feedback for Aether'
                     '&body=Hello there! Thanks for taking time to give feedback, I really appreciate it. '
                     'If you are having problems, please follow the directions at www.getaether.net/sending_logs, '
                     'and send me the produced logs. Thanks! You can delete this text before sending. '
                     'You can find my PGP key here: pgp.mit.edu:11371/pks/lookup?search=Burak+Nehbit')
            QtGui.QDesktopServices.openUrl(mailInitialiser)
        self.emailDevMenuItem.triggered.connect(emailDev)

        self.menu.addSeparator() # 5

        self.menu.addAction('Settings')
        self.settingsMenuItem = self.menu.actions()[8]
        def goToSettings():
            self.settingsMenuItem.setText('Settings')
            self.settingsMenuItem.setDisabled(False)
            if parent.isHidden():
                parent.show()
                parent.raise_()
            jsString = \
            ("firstFrameScope = angular.element(document.getElementById('first-frame-contents')).scope();"
             "firstFrameScope.settingsButtonClick();"
             "firstFrameScope.$apply();"
            )
            self.webView.JSContext(jsString)
        self.settingsMenuItem.triggered.connect(goToSettings)

        self.menu.addSeparator() # 6

        self.menu.addAction('Quit Aether')
        self.quitAppMenuItem = self.menu.actions()[10]
        # This is below reactor.run to allow access from other places outside main.
        def quitApp():
            # This is buggy...
            if parent.reactor.threadpool is not None:
                parent.reactor.threadpool.stop()
            parent.close()
            parent.reactor.stop()
            app.quit()
            sys.exit()

            # self.protInstance.killApp()
            # def finishExit():
            #     parent.reactor.stop()
            #     app.quit()
            #     sys.exit()
            # d.addCallback(finishExit)
        self.quitAppMenuItem.triggered.connect(quitApp)



        self.setContextMenu(self.menu)
        self.setIcon(self.icon)

        def changeIconToActiveState():
            self.setIcon(self.iconActive)
        def changeIconToPassiveState():
            self.setIcon(self.icon)
        self.menu.aboutToShow.connect(changeIconToActiveState)
        self.menu.aboutToHide.connect(changeIconToPassiveState)
        if PLATFORM == 'WIN':
            def showOnLeftClick(reason):
                if reason == self.Trigger:
                    makeVisible() # I hate that Python doesn't have anonymous functions.
            self.activated.connect(showOnLeftClick)

    def lightUpIcon(self):
        self.setIcon(self.iconHighlight)
        self.messagesMenuItem.setText('New replies available.')
        self.messagesMenuItem.setDisabled(False)

    def makeIconGoDark(self):
        self.setIcon(self.icon)
        self.messagesMenuItem.setText('You have no replies.')
        self.messagesMenuItem.setDisabled(True)