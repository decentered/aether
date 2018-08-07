// Services > Frontend Supervisor
// This service handles the boot process of the frontend and if there is any that needs to be booted, the backend, by proxy.
export { }

const { spawn } = require('child_process')
let globals = require('../globals/globals')

let clientAPIServerIP: string = '127.0.0.1'

export function StartFrontendDaemon(clientAPIServerPort: number): boolean {
  if (globals.FrontendDaemonStarted) {
    console.log("frontend daemon already running. skipping the start.")
    return false
  }
  globals.FrontendDaemonStarted = true
  // This is where we start the frontend binary.
  console.log("Frontend daemon starting")
  const child = spawn('go', ['run', '../frontend/main.go', 'run', '--logginglevel=1', `--clientip=${clientAPIServerIP}`, `--clientport=${clientAPIServerPort}`], {
    // env: {}, // no access to environment, enabled this in prod to make sure that the app can run w/out depending on anything
    detached: false,
    // stdio: 'ignore', // enable this in prod, we don't need any i/o in stdio
  })
  // child.unref() // Unreference = means it can continue running even when client shuts down. todo: figure out how to make best use of this, we want the frontend to shut down but maybe not the backend? do we want client to have code that searches for an existing fe?
  child.on('exit', function(code: any, signal: any) {
    console.log(globals)
    globals.FrontendDaemonStarted = false
    console.log('Frontend process exited with ' +
      `code ${code} and signal ${signal}`)
    console.log('We will reattempt to start the frontend daemon in 10 seconds.')
    setTimeout(function() {
      console.log('Attempting to restart the frontend now.')
      console.log(globals.ClientAPIServerPort)
      console.log(globals)
      StartFrontendDaemon(globals.ClientAPIServerPort)
    }, 10000)
  })
  child.stdout.on('data', (data: any) => {
    console.log(`${data}`)
  })
  child.stderr.on('data', (data: any) => {
    console.error(`${data}`)
  })
  return true
}