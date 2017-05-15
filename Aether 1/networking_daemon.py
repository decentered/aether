# Python stdlib imports
from __future__ import print_function
import globals
from globals import FROZEN, PLATFORM, PROFILE_DIR
import datetime

# Twisted imports
from twisted.internet import reactor
from twisted.internet.task import LoopingCall
from twisted.internet.endpoints import SSL4ServerEndpoint

# Aether imports
from InputOutput import aetherProtocol
from DecisionEngine import eventLoop
from ORM import Demeter


# Without this line, the networking process can't communicate with main on Windows when frozen.
# Possibly also valid for OS X. Ignore PyCharm 'unused import' warning.
from InputOutput.interprocessChildProt import NetworkingProcessProtocol

if globals.userProfile.get('debugDetails', 'debugLogging'):  # Debug enabled. Keep print enabled, and route it to the logs.
    from twisted.python import log
    from twisted.python.logfile import DailyLogFile
    log.startLogging(DailyLogFile.fromFullPath(PROFILE_DIR+'/Logs/network.log'))
    globals.logSystemDetails()

else:  # Debug not enabled. Disable print
    def print(*a, **kwargs):
        pass


if FROZEN:
    print('Networking Daemon: I am frozen.')
else:
    print('Networking Daemon: I am thawed.')

def main():
    ampChildPath = 'InputOutput.interprocessChildProt.NetworkingProcessProtocol'
    from twisted.internet import stdio
    from twisted.python import reflect
    protInstance = reflect.namedAny(ampChildPath)(reactor)  # Invoke what you found.
    Demeter.committer.receiveInterprocessProtocolInstance(protInstance)
    d = Demeter.checkUPNPStatus(2000)

    def maybeCommit():
        print('Attempting to commit. Commit In Progress flag is %s, last commit was at %s' %
              (Demeter.committer.commitInProgress, Demeter.committer.lastCommit))
        if Demeter.committer.commitInProgress is False:
            print('Commit loop decided to commit. #COMMITSTART')
            Demeter.committer.commit()
        else:
            print('Commit loop decided to NOT commit, variables are: commitInProgress = %s, lastCommit = %s '
                  '#COMMITDECLINE' % (Demeter.committer.commitInProgress, Demeter.committer.lastCommit))

    persephone = LoopingCall(maybeCommit)
    persephone.start(30)  # this should be 60 under normal circumstances.
    marduk = LoopingCall(eventLoop.marduk, aetherProtocol.aetherProtocolFactoryInstance, Demeter.committer)
    marduk.start(30)  # Should this be 5 minutes? (300)
    listenerEndpoint = SSL4ServerEndpoint(reactor, globals.userProfile.get('machineDetails', 'listeningPort'), globals.AetherContextFactory())
    listenerEndpoint.listen(aetherProtocol.aetherProtocolFactoryInstance)

    if FROZEN:
        if PLATFORM == 'WIN':
            stdio.StandardIO(protInstance)
        elif PLATFORM == 'OSX':
            stdio.StandardIO(protInstance, 7, 8)
        elif PLATFORM == 'LNX':
            stdio.StandardIO(protInstance, 7, 8)
    else:
        if PLATFORM == 'WIN':
            stdio.StandardIO(protInstance)
        elif PLATFORM == 'OSX':
            stdio.StandardIO(protInstance, 7, 8)
        elif PLATFORM == 'LNX':
            stdio.StandardIO(protInstance, 7, 8)

    reactor.run()

if __name__ == "__main__":
    main()