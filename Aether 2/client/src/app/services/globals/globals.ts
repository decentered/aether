// Services > Global Constants


// var ipc = require('../../../../node_modules/electron-better-ipc')
// const fesupervisor = require('../fesupervisor/fesupervisor')

interface Globals {
  FrontendReady: boolean;
  FrontendAPIPort: number;
  FrontendClientConnInitialised: boolean;
  ClientAPIServerPort: number;
  FrontendDaemonStarted: boolean;
}

let glob: Globals = {
  FrontendReady: false,
  FrontendAPIPort: 0,
  FrontendClientConnInitialised: false,
  ClientAPIServerPort: 0,
  FrontendDaemonStarted: false,
}

module.exports = glob


