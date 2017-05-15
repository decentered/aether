import globals
from globals import FROZEN
from interprocessAPI import *
from ampoule_aether import child
from ORM import Hermes, Demeter

from InputOutput.aetherProtocol import connectWithIP as connWithIP # Clashes with the interprocess call name.
import os

# This is the code responding to the protocol in the child.
class NetworkingProcessProtocol(child.AMPChild):

    def __init__(self, reactor):
        super(NetworkingProcessProtocol, self).__init__()
        self.reactor = reactor
        print('Networking Process Protocol initialized')

    def notifyMainProcessOfNewReplies(self):
        self.callRemote(thereAreReplies)

    @commit.responder
    def guiCommitted(self, PostFingerprint):
        #committer.commit()
        # TODO: Do other stuff required when user adds something.
        # I need to find the last n dudes I connected at, and ignore the cooldown and connect back to them
        # to serve my newest shit.
        # Okay, get the connectToNode method
        Demeter.committer.newPostsToIncrement.append(PostFingerprint)
        Demeter.connectToLastConnected(10)
        print('I received a commit signal from main (gui) process.')
        return {}

    @killApp.responder
    def killAppResponder(self):
        # This is probably not needed in the new arch. TODO
        print('I received a kill signal. KTHXBAI')
        d = self.processPool.stop()
        def stopR(*a):
            self.reactor.stop()
        d.addCallback(stopR)
        return {}

    @connectWithIP.responder
    def respondToConnectButton(self, IP, Port):
        print('I received a connect request to %s:%s' %(IP, Port))
        try:
            assert(isinstance(IP, basestring))
            assert(isinstance(int(Port), int))
        except:
            return {}
        else:
            if globals.checkIPValidity(IP) and Port > 2 and Port < 65536:
                globals.userProfile.set('machineDetails', 'lastConnectedBareIP', IP)
                globals.userProfile.set('machineDetails', 'lastConnectedBarePort', Port)
                connWithIP(IP, Port)
            else:
                print('The IP or Port the user has provided is invalid. IP: %s, Port: %s' % (IP, Port))
            return {}


# The thing below is useful for testing the thawed version. childBootstrap will only be called in case the app is
# running thawed, and in that case it will just pipe the import as text over to the concerned library.

if not FROZEN:
    with open('networking_daemon.py') as daemon_as_text:
        childBootstrap = daemon_as_text.read()