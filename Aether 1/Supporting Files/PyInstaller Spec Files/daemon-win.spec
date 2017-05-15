# -*- mode: python -*-
a = Analysis(['PATH TO DIRECTORY OF NETWORKING DAEMON'],
             hiddenimports=[],
             hookspath=None,
             runtime_hooks=None)
pyz = PYZ(a.pure)
exe = EXE(pyz,
          a.scripts, # You need to have these files. These are provided by Visual Studio. Express version is fine.
          [('msvcp100.dll', 'C:\\Windows\\System32\\msvcp100.dll', 'BINARY'),
          ('msvcr100.dll', 'C:\\Windows\\System32\\msvcr100.dll', 'BINARY')],
          exclude_binaries=True,
          name='Aether_Networking_Process.exe',
          debug=False,
          strip=None,
          upx=True,
          console=False,
          icon='PATH TO AETHER.ICO')

coll = COLLECT(exe,
               a.binaries +
               a.zipfiles,
               a.datas,
               strip=None,
               upx=True,
               name='networking_daemon'
               )
