#from twisted.protocols import amp
from __future__ import print_function
from interprocessAPI import *
from twisted.internet import reactor
from ORM.Demeter import committer, connectToLastConnected
import globals
from globals import PLATFORM, BASEDIR
from PyQt5 import QtWidgets

if not globals.userProfile.get('debugDetails', 'debugLogging'):
    def print(*a, **kwargs):
        pass

class MainProcessProtocol(amp.AMP):
    def killApp(self):
        # This is probably not needed in the new arch. TODO
        print('sending kill signal')
        return self.callRemote(killApp)

    def commit(self, PostFingerprint):
        print('commit is called to parent')
        return self.callRemote(commit, PostFingerprint=PostFingerprint)

    def connectWithIP(self, IP, Port):
        return self.callRemote(connectWithIP, IP=IP, Port=int(Port))

    # These are methods arriving from Main thread to GUI thread (this thread.)

    @thereAreReplies.responder
    def respondToThereAreReplies(self):
        # Do stuff here.

        # Here I need to check the reply count and post that reply count to the main page by changing the scope element.
        print('I received a NEW REPLIES signal')
        print('globals notification shown is', globals.notificationShown)
        if not globals.notificationShown:
            # Prevents further notifications from happening, unless user clears them up.
            globals.notificationShown = True
            if PLATFORM == 'OSX':
                self.trayIcon.showMessage('New Messages', 'You have new replies.', QtWidgets.QSystemTrayIcon.Information)
                self.trayIcon.lightUpIcon()
            elif PLATFORM == 'WIN':
                self.trayIcon.showMessage('Aether', 'You have new replies.', QtWidgets.QSystemTrayIcon.Information)
                self.trayIcon.lightUpIcon()
            elif PLATFORM == 'LNX':
                try:
                    import subprocess
                    pid = subprocess.Popen(['notify-send',
                                            'Aether',
                                            'You have new replies.',
                                            '--icon=' + BASEDIR + 'Assets/splash.png']).pid
                except: pass
        replyCount = self.Hermes.countReplies()
        # This doesn't seem to be working. It might be that it's running attached to PyCharm, or it might legitimately be broken. TEST
        jsString = \
            ("rootScope = angular.element(document.getElementById('root-body')).scope();"
             "rootScope.totalReplyCount = %s;"
             "rootScope.$apply();" % replyCount
            )
        self.JSContext(jsString)
        return {}

    # def connectionLost(self, reason):
    #     #do nothing?
    #     pass
    #     # amp.AMP.connectionLost(self, reason)
    #     # from twisted.internet import reactor
    #     # try:
    #     #     reactor.stop()
    #     # except error.ReactorNotRunning:
    #     #     # woa, this means that something bad happened,
    #     #     # most probably we received a SIGINT. Now this is only
    #     #     # a problem when you use Ctrl+C to stop the main process
    #     #     # because it would send the SIGINT to child processes too.
    #     #     # In all other cases receiving a SIGINT here would be an
    #     #     # error condition and correctly restarted. maybe we should
    #     #     # use sigprocmask?
    #     #     pass
    #     # if not self.shutdown:
    #     #     # if the shutdown wasn't explicit we presume that it's an
    #     #     # error condition and thus we return a -1 error returncode.
    #     #     import os
    #     #     os._exit(-1)

