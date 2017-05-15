from __future__ import print_function
import os
import sys
import imp
import sets
import itertools

from zope.interface import implements

from twisted.internet import reactor, protocol, defer, error
from twisted.python import log, util, reflect
from twisted.protocols import amp
from twisted.python import runtime

from ampoule_aether import iampoule

import globals
from globals import FROZEN, BASEDIR, PLATFORM

if not globals.userProfile.get('debugDetails', 'debugLogging'):
    def print(*a, **kwargs):
        pass

gen = itertools.count()

if runtime.platform.isWindows():
    IS_WINDOWS = True
    TO_CHILD = 0
    FROM_CHILD = 1
else:
    IS_WINDOWS = False
    TO_CHILD = 7
    FROM_CHILD = 8

class AMPConnector(protocol.ProcessProtocol):
    """
    A L{ProcessProtocol} subclass that can understand and speak AMP.

    @ivar amp: the children AMP process
    @type amp: L{amp.AMP}

    @ivar finished: a deferred triggered when the process dies.
    @type finished: L{defer.Deferred}

    @ivar name: Unique name for the connector, much like a pid.
    @type name: int
    """

    def __init__(self, proto, name=None):
        """
        @param proto: An instance or subclass of L{amp.AMP}
        @type proto: L{amp.AMP}

        @param name: optional name of the subprocess.
        @type name: int
        """
        self.finished = defer.Deferred()
        self.amp = proto
        self.name = name
        if name is None:
            self.name = gen.next()

    def signalProcess(self, signalID):
        """
        Send the signal signalID to the child process

        @param signalID: The signal ID that you want to send to the
                        corresponding child
        @type signalID: C{str} or C{int}
        """
        return self.transport.signalProcess(signalID)

    def connectionMade(self):
        print("Subprocess %s started." % (self.name,))
        self.amp.makeConnection(self)

    # Transport
    disconnecting = False

    def write(self, data):
        if IS_WINDOWS:
            self.transport.write(data)
        else:
            self.transport.writeToChild(TO_CHILD, data)

    def loseConnection(self):
        self.transport.closeChildFD(TO_CHILD)
        self.transport.closeChildFD(FROM_CHILD)
        self.transport.loseConnection()

    def getPeer(self):
        return ('subprocess',)

    def getHost(self):
        return ('no host',)

    def childDataReceived(self, childFD, data):
        if childFD == FROM_CHILD:
            self.amp.dataReceived(data)
            return
        self.errReceived(data)

    def errReceived(self, data):
        for line in data.strip().splitlines():
            print("FROM %s: %s" % (self.name, line))

    def processEnded(self, status):
        print("Process: %s ended" % (self.name,))
        self.amp.connectionLost(status)
        if status.check(error.ProcessDone):
            self.finished.callback('')
            return
        self.finished.errback(status)

BOOTSTRAP = """\
import sys

def main(reactor, ampChildPath):
    from twisted.application import reactors
    reactors.installReactor(reactor)

    from twisted.python import log
    log.startLogging(sys.stderr)

    from twisted.internet import reactor, stdio
    from twisted.python import reflect, runtime

    ampChild = reflect.namedAny(ampChildPath)
    if runtime.platform.isWindows():
        stdio.StandardIO(ampChild(*sys.argv[1:-2]))
    else:
        stdio.StandardIO(ampChild(*sys.argv[1:-2]), %s, %s)
    enter = getattr(ampChild, '__enter__', None)
    if enter is not None:
        enter()
    try:
        reactor.run()
    except:
        if enter is not None:
            info = sys.exc_info()
            if not ampChild.__exit__(*info):
                raise
        else:
            raise
    else:
        if enter is not None:
            ampChild.__exit__(None, None, None)

main(sys.argv[-2], sys.argv[-1])
""" % (TO_CHILD, FROM_CHILD)

class ProcessStarter(object):

    implements(iampoule.IStarter)

    connectorFactory = AMPConnector
    def __init__(self, bootstrap=BOOTSTRAP, args=(), env={},
                 path=None, uid=None, gid=None, usePTY=0,
                 packages=(), childReactor="select", isProcess=False):
        """
        @param bootstrap: Startup code for the child process
        @type  bootstrap: C{str}

        @param args: Arguments that should be supplied to every child
                     created.
        @type args: C{tuple} of C{str}

        @param env: Environment variables that should be present in the
                    child environment
        @type env: C{dict}

        @param path: Path in which to run the child
        @type path: C{str}

        @param uid: if defined, the uid used to run the new process.
        @type uid: C{int}

        @param gid: if defined, the gid used to run the new process.
        @type gid: C{int}

        @param usePTY: Should the child processes use PTY processes
        @type usePTY: 0 or 1

        @param packages: A tuple of packages that should be guaranteed
                         to be importable in the child processes
        @type packages: C{tuple} of C{str}

        @param childReactor: a string that sets the reactor for child
                             processes
        @type childReactor: C{str}
        """
        self.bootstrap = bootstrap
        self.args = args
        self.env = env
        self.path = path
        self.uid = uid
        self.gid = gid
        self.usePTY = usePTY
        self.packages = packages
        self.childReactor = childReactor
        self.isProcess = isProcess

    def __repr__(self):
        """
        Represent the ProcessStarter with a string.
        """
        return """ProcessStarter(bootstrap=%r,
                                 args=%r,
                                 env=%r,
                                 path=%r,
                                 uid=%r,
                                 gid=%r,
                                 usePTY=%r,
                                 packages=%r,
                                 childReactor=%r)""" % (self.bootstrap,
                                                        self.args,
                                                        self.env,
                                                        self.path,
                                                        self.uid,
                                                        self.gid,
                                                        self.usePTY,
                                                        self.packages,
                                                        self.childReactor)

    def _checkRoundTrip(self, obj):
        """
        Make sure that an object will properly round-trip through 'qual' and
        'namedAny'.

        Raise a L{RuntimeError} if they aren't.
        """
        tripped = reflect.namedAny(reflect.qual(obj))
        if tripped is not obj:
            raise RuntimeError("importing %r is not the same as %r" %
                               (reflect.qual(obj), obj))

    def startAMPProcess(self, ampChild, ampParent=None, ampChildArgs=()):
        """
        @param ampChild: a L{ampoule.child.AMPChild} subclass.
        @type ampChild: L{ampoule.child.AMPChild}

        @param ampParent: an L{amp.AMP} subclass that implements the parent
                          protocol for this process pool
        @type ampParent: L{amp.AMP}
        """
        self._checkRoundTrip(ampChild)
        fullPath = reflect.qual(ampChild)
        print(fullPath)
        if ampParent is None:
            ampParent = amp.AMP # UNTIL HERE it looks good. check sanity of rest.
        prot = self.connectorFactory(ampParent())
        args = ampChildArgs + (self.childReactor, fullPath)
        if self.isProcess:
            return self.startExternalProcess(location)
        else:
            return self.startPythonProcess(prot, *args)

    def startExternalProcess(self, location):
        startProcess(location)

    def startPythonProcess(self, prot, *args):
        """
        @param prot: a L{protocol.ProcessProtocol} subclass
        @type prot: L{protocol.ProcessProtocol}

        @param args: a tuple of arguments that will be added after the
                     ones in L{self.args} to start the child process.

        @return: a tuple of the child process and the deferred finished.
                 finished triggers when the subprocess dies for any reason.
        """
        spawnProcess(prot, self.bootstrap, self.args+args, env=self.env,
                     path=self.path, uid=self.uid, gid=self.gid,
                     usePTY=self.usePTY, packages=self.packages)

        # XXX: we could wait for startup here, but ... is there really any
        # reason to?  the pipe should be ready for writing.  The subprocess
        # might not start up properly, but then, a subprocess might shut down
        # at any point too. So we just return amp and have this piece to be
        # synchronous.
        return prot.amp, prot.finished

def startProcess(location):
    import subprocess

def spawnProcess(processProtocol, bootstrap, args=(), env={},
                 path=None, uid=None, gid=None, usePTY=0,
                 packages=()):
    env = env.copy()

    pythonpath = []
    # app_is_frozen = False
    for pkg in packages:
        if FROZEN:
            p = BASEDIR
        else:
            if PLATFORM == 'WIN' and pkg == 'ampoule':
                pkg = 'ampoule_aether'
            p = os.path.split(imp.find_module(pkg)[1])[0]
        # if getattr(sys, 'frozen', False): # If this is a frozen PyInstaller bundle.
        #     p = sys._MEIPASS+'/'
        #     print('from ampoule: we are frozen. setting p as static')
        #     app_is_frozen = True
        # else:
        #     print('from ampoule: we not frozen. setting p as normal')
        #     if IS_WINDOWS and pkg == 'ampoule':
        #         pkg = 'ampoule_aether'
        #     p = os.path.split(imp.find_module(pkg)[1])[0]
        #     app_is_frozen = False
        if p.startswith(os.path.join(sys.prefix, 'lib')):
            continue
        pythonpath.append(p)
    pythonpath = list(sets.Set(pythonpath))
    pythonpath.extend(env.get('PYTHONPATH', '').split(os.pathsep))
    env['PYTHONPATH'] = os.pathsep.join(pythonpath)
    args = (sys.executable, '-c', bootstrap) + args
    if FROZEN and PLATFORM == 'LNX':
        args = []
        # This is weird. Under Linux, if args is not empty and frozen, it tries to use the sys.executable, even if
        # there is a process app is given. This is required for Linux, but also might have been creating problems in
        # Windows and OS X too. Possibly. Probably not.
    # childFDs variable is needed because sometimes child processes
    # misbehave and use stdout to output stuff that should really go
    # to stderr. Of course child process might even use the wrong FDs
    # that I'm using here, 3 and 4, so we are going to fix all these
    # issues when I add support for the configuration object that can
    # fix this stuff in a more configurable way.
    print('sys executable is:',sys.executable)

    if FROZEN:
        print('THIS IS BASEDIR')
        print(BASEDIR)
        if PLATFORM == 'WIN': # Frozen and Windows.
            print('This is Windows and frozen.')
            print(BASEDIR + '/networking_daemon/Aether_Networking_Process.exe')
            return reactor.spawnProcess(processProtocol, BASEDIR + 'networking_daemon/Aether_Networking_Process.exe', args,
                                        env, path, uid, gid, usePTY)
        elif PLATFORM == 'OSX': # Frozen and sane.
            print('This is OSX and frozen.')
            print(BASEDIR + 'networking_daemon/networking_daemon')
            return reactor.spawnProcess(processProtocol, BASEDIR + 'networking_daemon/Aether Networking Process', args,
                                            env, path, uid, gid, usePTY,
                                            childFDs={0:"w", 1:"r", 2:"r", 7:"w", 8:"r"})
        elif PLATFORM == 'LNX': # Frozen and sane.
            print('This is Linux and frozen.')
            print(BASEDIR + 'networking_daemon_linux/Aether_Networking_Process')
            return reactor.spawnProcess(processProtocol, BASEDIR + 'networking_daemon_linux/Aether_Networking_Process', args,
                                            env, path, uid, gid, usePTY,
                                            childFDs={0:"w", 1:"r", 2:"r", 7:"w", 8:"r"})

    else:
        if PLATFORM == 'WIN': # Not frozen and Windows.
            print('This is Windows and not frozen.')
            return reactor.spawnProcess(processProtocol, sys.executable, args,
                                        env, path, uid, gid, usePTY)

        elif PLATFORM == 'OSX': # Not frozen and sane.
            print('This is OSX and not frozen.')
            return reactor.spawnProcess(processProtocol, sys.executable, args,
                                            env, path, uid, gid, usePTY,
                                            childFDs={0:"w", 1:"r", 2:"r", 7:"w", 8:"r"})

        elif PLATFORM == 'LNX': # Not frozen and sane.
            print('This is Linux and not frozen.')
            return reactor.spawnProcess(processProtocol, sys.executable, args,
                                            env, path, uid, gid, usePTY,
                                            childFDs={0:"w", 1:"r", 2:"r", 7:"w", 8:"r"})

