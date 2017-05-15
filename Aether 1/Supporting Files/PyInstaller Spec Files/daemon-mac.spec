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
          name='Aether Networking Process',
          debug=False,
          strip=None,
          upx=True,
          console=True)
coll = COLLECT(exe,
               a.binaries,
               a.zipfiles,
               a.datas,
               # dict_tree4,
               strip=None,
               upx=True,
               name='networking_daemon')
