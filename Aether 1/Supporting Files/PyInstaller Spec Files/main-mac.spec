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
          name='main',
          debug=False,
          strip=None,
          upx=False,
          console=False )

dict_tree3 = Tree('PATH TO GUI FOLDER', prefix = 'GUI')
dict_tree5 = Tree('PATH TO ASSETS FOLDER', prefix = 'Assets')
dict_tree6 = Tree('PATH TO FOLDER OF NETWORKING DAEMON', prefix = 'networking_daemon')

coll = COLLECT(exe,
               a.binaries,
               a.zipfiles,
               a.datas,
               dict_tree3,
               # dict_tree4,
               dict_tree5,
               dict_tree6,
               strip=None,
               upx=False,
               name='main')

app = BUNDLE(coll,
            name='Aether.app',
            icon='icon-windowed.icns')
