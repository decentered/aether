"use strict";
// Client > ClientAPIServer
// This file is the grpc server we want to use to talk to the frontend.
Object.defineProperty(exports, "__esModule", { value: true });
// Imports
var grpc = require('grpc');
// const resolve = require('path').resolve
// let globals = require('../globals/globals')
var feapiconsumer = require('../feapiconsumer/feapiconsumer');
var ipc = require('../../../../node_modules/electron-better-ipc');
var vuexStore = require('../../store/index').default;
// // Load the proto file
// const proto = grpc.load({
//   file: 'clapi/clapi.proto',
//   root: resolve(__dirname, '../protos')
// }).clapi
var messages = require('../../../../../protos/clapi/clapi_pb.js');
var services = require('../../../../../protos/clapi/clapi_grpc_pb');
/**
 Client-side GRPC server so that the frontend can talk to the client. This is useful at the first start where the Frontend needs to start its own GRPC server and return its address to the client.
 */
function StartClientAPIServer() {
    var server = new grpc.Server();
    server.addService(services.ClientAPIService, {
        frontendReady: FrontendReady,
        deliverAmbients: DeliverAmbients,
        sendAmbientStatus: SendAmbientStatus,
        sendAmbientLocalUserEntity: SendAmbientLocalUserEntity,
        sendHomeView: SendHomeView,
        sendPopularView: SendPopularView,
        sendNotifications: SendNotifications,
        sendOnboardCompleteStatus: SendOnboardCompleteStatus,
        sendModModeEnabledStatus: SendModModeEnabledStatus,
    });
    var boundPort = server.bind('127.0.0.1:0', grpc.ServerCredentials.createInsecure());
    server.start();
    return boundPort;
}
exports.StartClientAPIServer = StartClientAPIServer;
function FrontendReady(req, callback) {
    var r = req.request.toObject();
    console.log("frontend ready at: ", r.address, ":", r.port);
    // globals.FrontendReady = true
    ipc.callMain('SetFrontendReady', true);
    // globals.FrontendAPIPort = req.request.port
    ipc.callMain('SetFrontendAPIPort', r.port);
    feapiconsumer.Initialise();
    var resp = new messages.FEReadyResponse;
    callback(null, resp);
}
function DeliverAmbients(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setAmbientBoards', r.boardsList);
    var resp = new messages.AmbientsResponse;
    callback(null, resp);
}
function SendAmbientStatus(req, callback) {
    var r = req.request.toObject();
    // console.log(r)
    vuexStore.dispatch('setAmbientStatus', r);
    var resp = new messages.AmbientStatusResponse;
    callback(null, resp);
}
function SendAmbientLocalUserEntity(req, callback) {
    var r = req.request.toObject();
    // console.log(r)
    vuexStore.dispatch('setAmbientLocalUserEntity', r);
    var resp = new messages.AmbientLocalUserEntityResponse;
    callback(null, resp);
}
function SendHomeView(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setHomeView', r.threadsList);
    var resp = new messages.HomeViewResponse;
    callback(null, resp);
}
function SendPopularView(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setPopularView', r.threadsList);
    var resp = new messages.PopularViewResponse;
    callback(null, resp);
}
function SendNotifications(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setNotifications', r);
    var resp = new messages.NotificationsResponse;
    callback(null, resp);
}
function SendOnboardCompleteStatus(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setOnboardCompleteStatus', r.onboardcomplete);
    var resp = new messages.OnboardCompleteStatusResponse;
    callback(null, resp);
}
function SendModModeEnabledStatus(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setModModeEnabledStatus', r.modmodeenabled);
    var resp = new messages.ModModeEnabledStatusResponse;
    callback(null, resp);
}
//# sourceMappingURL=clapiserver.js.map