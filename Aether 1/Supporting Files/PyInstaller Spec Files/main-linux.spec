# -*- mode: python -*-
a = Analysis(['main.py'],
             pathex=['PATH TO DIRECTORY OF MAIN.PY'],
             hiddenimports=[],
             hookspath=None,
             runtime_hooks=None)
pyz = PYZ(a.pure)
exe = EXE(pyz,
          a.scripts,
          exclude_binaries=True,
          name='Aether',
          debug=False,
          strip=None,
          upx=True,
          console=False )
gui_folder = Tree('PATH TO GUI FOLDER', prefix = 'GUI')
assets_folder = Tree('PATH TO ASSETS FOLDER', prefix = 'Assets')
daemon_folder = Tree('PATH TO FOLDER OF NETWORKING DAEMON', prefix = 'networking_daemon_linux')

coll = COLLECT(exe,
           a.binaries,
           a.zipfiles,
           a.datas,
           gui_folder,
           assets_folder,
           daemon_folder,
           strip=None,
           upx=True,
           name='aether')
