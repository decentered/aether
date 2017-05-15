# -*- mode: python -*-
a = Analysis(['networking_daemon.py'],
             pathex=['PATH TO DIRECTORY OF NETWORKING DAEMON'],
             hiddenimports=[],
             hookspath=None,
             runtime_hooks=None)
pyz = PYZ(a.pure)
exe = EXE(pyz,
          a.scripts,
          exclude_binaries=True,
          name='Aether_Networking_Process',
          debug=False,
          strip=None,
          upx=False,
          console=True)
coll = COLLECT(exe,
               a.binaries,
               a.zipfiles,
               a.datas,
               strip=None,
               upx=False,
               name='networking_daemon_linux')
