// Client > ClientAPIServer
// This file is the grpc server we want to use to talk to the frontend.

export { } // This says this file is a module, not a script.

// Imports
const grpc = require('grpc')
// const resolve = require('path').resolve
// let globals = require('../globals/globals')
const feapiconsumer = require('../feapiconsumer/feapiconsumer')
var ipc = require('../../../../node_modules/electron-better-ipc')
var vuexStore = require('../../store/index').default


// // Load the proto file
// const proto = grpc.load({
//   file: 'clapi/clapi.proto',
//   root: resolve(__dirname, '../protos')
// }).clapi

var messages = require('../../../../../protos/clapi/clapi_pb.js')
var services = require('../../../../../protos/clapi/clapi_grpc_pb')


/**
 Client-side GRPC server so that the frontend can talk to the client. This is useful at the first start where the Frontend needs to start its own GRPC server and return its address to the client.
 */

export function StartClientAPIServer(): number {
  let server = new grpc.Server()
  server.addService(
    services.ClientAPIService, {
      frontendReady: FrontendReady,
      deliverAmbients: DeliverAmbients,
      sendAmbientStatus: SendAmbientStatus,
      sendAmbientLocalUserEntity: SendAmbientLocalUserEntity,
      sendHomeView: SendHomeView,
      sendPopularView: SendPopularView,
      sendNotifications: SendNotifications,
      sendOnboardCompleteStatus: SendOnboardCompleteStatus,
      sendModModeEnabledStatus: SendModModeEnabledStatus,
    }
  )
  let boundPort: number = server.bind(
    '127.0.0.1:0', grpc.ServerCredentials.createInsecure()
  )
  server.start()
  return boundPort
}

function FrontendReady(req: any, callback: any) {
  let r = req.request.toObject()
  console.log("frontend ready at: ", r.address, ":", r.port)
  // globals.FrontendReady = true
  ipc.callMain('SetFrontendReady', true)
  // globals.FrontendAPIPort = req.request.port
  ipc.callMain('SetFrontendAPIPort', r.port)
  feapiconsumer.Initialise()
  let resp = new messages.FEReadyResponse
  callback(null, resp)
}

function DeliverAmbients(req: any, callback: any) {
  let r = req.request.toObject()
  vuexStore.dispatch('setAmbientBoards', r.boardsList)
  let resp = new messages.AmbientsResponse
  callback(null, resp)
}

function SendAmbientStatus(req: any, callback: any) {
  let r = req.request.toObject()
  // console.log(r)
  vuexStore.dispatch('setAmbientStatus', r)
  let resp = new messages.AmbientStatusResponse
  callback(null, resp)
}

function SendAmbientLocalUserEntity(req: any, callback: any) {
  let r = req.request.toObject()
  // console.log(r)
  vuexStore.dispatch('setAmbientLocalUserEntity', r)
  let resp = new messages.AmbientLocalUserEntityResponse
  callback(null, resp)
}

function SendHomeView(req: any, callback: any) {
  let r = req.request.toObject()
  vuexStore.dispatch('setHomeView', r.threadsList)
  let resp = new messages.HomeViewResponse
  callback(null, resp)
}

function SendPopularView(req: any, callback: any) {
  let r = req.request.toObject()
  vuexStore.dispatch('setPopularView', r.threadsList)
  let resp = new messages.PopularViewResponse
  callback(null, resp)
}

function SendNotifications(req: any, callback: any) {
  let r = req.request.toObject()
  vuexStore.dispatch('setNotifications', r)
  let resp = new messages.NotificationsResponse
  callback(null, resp)
}

function SendOnboardCompleteStatus(req: any, callback: any) {
  let r = req.request.toObject()
  vuexStore.dispatch('setOnboardCompleteStatus', r.onboardcomplete)
  let resp = new messages.OnboardCompleteStatusResponse
  callback(null, resp)
}

function SendModModeEnabledStatus(req: any, callback: any) {
  let r = req.request.toObject()
  vuexStore.dispatch('setModModeEnabledStatus', r.modmodeenabled)
  let resp = new messages.ModModeEnabledStatusResponse
  callback(null, resp)
}
