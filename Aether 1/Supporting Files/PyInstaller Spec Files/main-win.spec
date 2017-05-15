# -*- mode: python -*-
a = Analysis(['PATH TO DIRECTORY OF MAIN.PY'],
             hiddenimports=[],
             hookspath=None,
             runtime_hooks=None)
pyz = PYZ(a.pure)
exe = EXE(pyz,
          a.scripts, # You need to have these files. These are provided by Visual Studio. Express version is fine.
          [('msvcp100.dll', 'C:\\Windows\\System32\\msvcp100.dll', 'BINARY'),
          ('msvcr100.dll', 'C:\\Windows\\System32\\msvcr100.dll', 'BINARY')],
          exclude_binaries=True,
          name='Aether.exe',
          debug=False,
          strip=None,
          upx=True,
          console=False,
          icon='PATH TO AETHER.ICO')

dict_tree3 = Tree('PATH TO GUI FOLDER', prefix = 'GUI')
dict_tree5 = Tree('PATH TO ASSETS FOLDER', prefix = 'Assets')
dict_tree6 = Tree('PATH TO FOLDER OF NETWORKING DAEMON', prefix = 'networking_daemon')

coll = COLLECT(exe,
               a.binaries +
               a.zipfiles,
               a.datas,
               dict_tree3,
               dict_tree5,
               dict_tree6,
               strip=None,
               upx=True,
               name='Aether'
               )
app = BUNDLE(coll,
             name='Aether',
             icon='PATH TO AETHER.ICO')
