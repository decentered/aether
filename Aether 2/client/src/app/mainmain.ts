// Electron main.
/*
    This file is the main execution path of Electron. It starts up Electron, and loads our HTML main html file. In this HTML file is our app contained, its JS code, etc.

    In other words, anything here runs as an OS-privileged executable a la Node. Anything that takes root from the main HTML file runs effectively as a web page.
*/

var globals = require('./services/globals/globals') // Register globals
require('./services/eipc/eipc-main') // Register IPC events
var ipc = require('../../node_modules/electron-better-ipc') // Register IPC caller
const elc = require('electron')
// const starters = require('./starters')
// const feapiconsumer = require('./services/feapiconsumer/feapiconsumer')
const minimatch = require('../../node_modules/minimatch')
// var ipc = require('../../node_modules/electron-better-ipc')
// const fesupervisor = require('./services/fesupervisor/fesupervisor')

// Enable live reload. This should be disabled in production. TODO
const path = require('path')
const maindir = path.dirname(__dirname)
// require('electron-reload')(maindir)


require('electron-context-menu')({
  // prepend: (params, browserWindow) => [{
  // prepend: () => [{
  //   label: 'Rainbow',
  //   // Only show it when right-clicking images
  //   visible: true
  // }]
})


// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win: any
let tray: any = null
var DOM_READY: boolean = false


// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
elc.app.on('ready', main)

// Quit when all windows are closed.
elc.app.on('window-all-closed', () => {
  // On macOS it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
    elc.app.quit()
  }
})

elc.app.on('activate', () => {
  // On macOS it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (win === null) {
    main()
  }
})

function EstablishExternalResourceAutoLoadFiltering() {
  // This is the list of allowed URLs that can be auto-loaded in electron. (This does not prevent links that open in your browser, just ones that fetch data within the app. You can link anywhere, but only links from the whitelist below will have a chance to auto-load.)

  /*
    Why does this list even exist? Shouldn't people be able to link to everywhere?

    You *can* link to everywhere, this list is just for auto-loading previews. Why does this matter? Because when an asset is auto-loaded, the entity on the other end (i.e. the sites below) will see your IP address as a normal user. That means, if there was no whitelist and all links were allowed to auto-load, then anybody could link to a site they control, and by listening to hits from IP addresses, it could figure out which IP addresses are using Aether. It wouldn't be able to figure out who is who, but the fact that IP is using Aether would be revealed.

    If you'd like to make things a little tighter in exchange to not being able to preview, replace this list with an empty one, and all auto-loads will be blocked.
  */

  // This list should be editable. (TODO)
  const whitelist = [
    'https://*.imgur.com/*',
    'https://imgur.com/*',
    'https://gfycat.com/*',
    'https://*.gfycat.com/*',
    'https://giphy.com/*',
    'https://*.giphy.com/*',
    'file://**', // So that we can load the local client app itself
    'chrome-devtools://**',
    'chrome-extension://**' // for vue devtools
  ]
  // Allow any auto-load request that's in the whitelist. Deny autoload requests to all other domains.
  elc.session.defaultSession.webRequest.onBeforeRequest(function(details: any, cb: any) {
    // console.log(details.url) // Uncomment this to see all attempted outbound network requests from the client app.
    for (let i = 0; i < whitelist.length; i++) {
      if (minimatch(details.url, whitelist[i], { matchBase: true })) {
        cb({ cancel: false })
        return
      }
    }
    cb({ cancel: true })
  })
}



function EstablishElectronWindow() {
  // If not prod, install Vue devtools.
  if (process.env.NODE_ENV !== 'production') {
    require('vue-devtools').install()
  }

  // Create the browser window.
  let dm = elc.screen.getPrimaryDisplay().size
  win = new elc.BrowserWindow({
    show: false,
    width: dm.width * 0.8,
    height: dm.height * 0.8,
    titleBarStyle: 'hiddenInset',
    // titleBarStyle: 'customButtonsOnHover',
    // frame: false,
    title: 'Aether',
    fullscreenWindowTitle: true,
    darkTheme: true, // Linux only for now, for GTK3+
    backgroundColor: '#292b2f',
    disableBlinkFeatures: "Auxclick" // disable middle click new window
    // webPreferences: {
    //   blinkFeatures: 'OverlayScrollbars'
    // },
  })
  win.once('ready-to-show', function() {
    // We want to show the window only after Electron is done readying itself.
    setTimeout(function() {
      win.show()
    }, 100)
    // Unfortunately, there's a race condition from the Electron side here (I might be making a mistake also, but it is simple enough to reproduce that there is not much space for me to make a mistake). If the setTimeout is 0 or is not present, there's about 1/10 chance the window is painted but completely frozen. Having 100ms seems to make it go away, but it's a little icky, because that basically is my guess. Not great. Hopefully they'll fix this in upcoming Electron 3.
  })
  win.loadFile('index.html')
  // Open the DevTools.
  // win.webContents.openDevTools({ mode: 'bottom' })
  // Emitted when the window is closed.
  win.on('closed', () => {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    DOM_READY = false
    win = null
  })

  // win.webContents.on('new-window', function(e, url) {
  //   e.preventDefault()
  //   elc.shell.openExternal(url)
  // })
  win.webContents.on('dom-ready', function() {
    DOM_READY = true
    // This is needed because the renderer process won't be able to respond to IPC requests before this event happens.
  })
  win.webContents.on('will-navigate', function(e: any, reqUrl: any) {
    e.preventDefault()
    elc.shell.openExternal(reqUrl)
    // return
    // let getHost = function(url: any) { require('url').parse(url).host }
    // let reqHost = getHost(reqUrl)
    // let isExternal = reqHost && reqHost != getHost(win.webContents.getURL())
    // if (isExternal) {
    //   e.preventDefault()
    //   elc.shell.openExternal(reqUrl)
    // }
  })

  win.webContents.on('new-window', function(e: any) {
    e.preventDefault()
  })
}


function EstablishTray() {
  if (tray !== null) { return }
  /*----------  Tray functions  ----------*/
  let openApp = function() {
    if (win === null) {
      main()
    }
    win.focus()
  }
  let openPreferences = function() {
    openApp()
    let rendererReadyChecker = function() {
      if (!(globals.RendererReady && DOM_READY)) {
        return setTimeout(rendererReadyChecker, 100)
      }
      return ipc.callRenderer(win, 'RouteTo', '/settings')
    }
    setTimeout(rendererReadyChecker, 100)
  }
  let openSupport = function() {
    elc.shell.openExternal('https://meta.getaether.net')
  }

  let quitApp = function() {
    elc.app.quit()
  }
  /*----------  Tray functions END  ----------*/
  tray = new elc.Tray(path.join(__dirname, 'ext_dep/images/TrayTemplate.png'))
  const contextMenu = elc.Menu.buildFromTemplate([
    { label: 'Online', enabled: false },
    { type: 'separator' },
    { label: 'Open Aether', click: openApp },
    { type: 'separator' },
    { label: 'Preferences...', click: openPreferences },
    { label: 'Community support', click: openSupport },
    { type: 'separator' },
    { label: 'Quit Aether', click: quitApp }
  ])
  tray.setToolTip('Aether: the front page of /dev/null.')
  tray.setContextMenu(contextMenu)
}

ipc.answerRenderer('QuitApp', function() {
  return elc.app.quit()
})

/**
  This is the main() of Electron. It starts the Client GRPC server, and kicks of the frontend and the backend daemons.
*/

function main() {
  console.log("mainmain reruns")
  EstablishExternalResourceAutoLoadFiltering()
  EstablishElectronWindow()
  EstablishTray()
}

