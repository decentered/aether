from __future__ import print_function
import globals
from ORM import Demeter
from InputOutput.aetherProtocol import connectWithNode, connectWithIP

if not globals.userProfile.get('debugDetails', 'debugLogging'):
    def print(*a, **kwargs):
        pass

maxOutboundCount = globals.userProfile.get('machineDetails', 'maxOutboundCount')
cooldown = globals.userProfile.get('machineDetails', 'cooldown')

def marduk(factoryInstance, committerInstance):
    if globals.appIsPaused: # If the app is explicitly paused, don't create active network events.
        return

    d3 = Demeter.getNodesToConnect(maxOutboundCount, cooldown)
    print('max outbound count:', maxOutboundCount, 'cooldown:', cooldown)
    def connectToNodes(nodes, openConnsList):
        # if not committerInstance.commitInProgress: # Don't connect outbound if a commit is in progress.
        if len(nodes) == 0:
            print('There are no nodes to connect to!')
            lastBareIP = globals.userProfile.get('machineDetails', 'lastConnectedBareIP')
            lastBarePort = globals.userProfile.get('machineDetails', 'lastConnectedBarePort')
            if isinstance(lastBarePort, int) and isinstance(lastBareIP, basestring) and \
                            lastBarePort != 0 and not Demeter.committer.commitInProgress:
            # Only connect if there is a bare connection request before and there is no commit in progress.
            # This is to prevent the first connection repeatedly calling over and wasting bandwidth while he has
            # everything in memory, just committing.
                print('Reaching out to the last bare '
                      'connect attempt at %s : %s' % (lastBareIP, lastBarePort))
                connectWithIP(lastBareIP, lastBarePort)
            else:
                print('This computer has never bare-connected to another node or a commit is in progress. '
                      'This usually means onboarding is not completed. lastBarePort = %s, commitInProgress = %s' %
                      (lastBarePort, Demeter.committer.commitInProgress))
        print('Number of open connections:', len(openConnsList))
        for n in nodes:
             print('I\'m attempting to connect to node %s at %s:%s'
                   %(n.NodeId,
                     n.LastConnectedIP if n.LastConnectedIP is not None else n.LastRetrievedIP,
                     n.LastConnectedPort if n.LastConnectedPort is not None else n.LastRetrievedPort))
             connectWithNode(n.asDict())
        # else:
        #     print('A commit is in progress. Skipping this round of marduk.')
    d3.addCallback(connectToNodes, factoryInstance.openConnections)


