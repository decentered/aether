// Electron main.
/*
    This file is the main execution path of Electron. It starts up Electron, and loads our HTML main html file. In this HTML file is our app contained, its JS code, etc.

    In other words, anything here runs as an OS-privileged executable a la Node. Anything that takes root from the main HTML file runs effectively as a web page.
*/

require('./services/globals/globals') // Register globals
require('./services/eipc/eipc') // Register IPC events

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
  prepend: () => [{
    label: 'Rainbow',
    // Only show it when right-clicking images
    visible: true
  }]
})


// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win: any

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

  // This list should be editable. (TODO)
  const whitelist = [
    'https://*.imgur.com/*',
    'https://imgur.com/*',
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
    width: dm.width * 0.8,
    height: dm.height * 0.8,
    titleBarStyle: 'hiddenInset',
    // titleBarStyle: 'customButtonsOnHover',
    // frame: false,
    title: 'Aether',
    fullscreenWindowTitle: true,
    darkTheme: true, // Linux only for now, for GTK3+
    backgroundColor: '#162127',
    disableBlinkFeatures: "Auxclick" // disable middle click new window
    // webPreferences: {
    //   blinkFeatures: 'OverlayScrollbars'
    // },
  })
  // and load the index.html of the app.
  win.loadFile('index.html')
  // Open the DevTools.
  win.webContents.openDevTools({ mode: 'bottom' })
  // Emitted when the window is closed.
  win.on('closed', () => {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    win = null
  })

  // win.webContents.on('new-window', function(e, url) {
  //   e.preventDefault()
  //   elc.shell.openExternal(url)
  // })

  win.webContents.on('will-navigate', function(e: any, reqUrl: any) {
    let getHost = (url: any) => require('url').parse(url).host
    let reqHost = getHost(reqUrl)
    let isExternal = reqHost && reqHost != getHost(win.webContents.getURL())
    if (isExternal) {
      e.preventDefault()
      elc.shell.openExternal(reqUrl)
    }
  })

  win.webContents.on('new-window', function(e: any) {
    e.preventDefault()
  })
}

/**
  This is the main() of Electron. It starts the Client GRPC server, and kicks of the frontend and the backend daemons.
*/

function main() {
  console.log("mainmain reruns")
  EstablishExternalResourceAutoLoadFiltering()
  EstablishElectronWindow()

  // setTimeout(function() {
  //   feapiconsumer.GetAllBoards()
  // }, 10000)
}

