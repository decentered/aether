"use strict";
// Client > ClientAPIServer
// This file is the grpc server we want to use to talk to the frontend.
Object.defineProperty(exports, "__esModule", { value: true });
// Imports
var grpc = require('grpc');
var resolve = require('path').resolve;
// let globals = require('../globals/globals')
var feapiconsumer = require('../feapiconsumer/feapiconsumer');
var ipc = require('../../../../node_modules/electron-better-ipc');
var vuexStore = require('../../store/index').default;
// Load the proto file
var proto = grpc.load({
    file: 'clapi/clapi.proto',
    root: resolve(__dirname, '../protos')
}).clapi;
/**
 Client-side GRPC server so that the frontend can talk to the client. This is useful at the first start where the Frontend needs to start its own GRPC server and return its address to the client.
 */
function StartClientAPIServer() {
    var server = new grpc.Server();
    server.addService(proto.ClientAPI.service, {
        FrontendReady: FrontendReady,
        DeliverAmbients: DeliverAmbients
    });
    var boundPort = server.bind('127.0.0.1:0', grpc.ServerCredentials.createInsecure());
    server.start();
    return boundPort;
}
exports.StartClientAPIServer = StartClientAPIServer;
function FrontendReady(req, callback) {
    console.log("frontend ready at: ", req.request.address, ":", req.request.port);
    // globals.FrontendReady = true
    ipc.callMain('SetFrontendReady', true);
    // globals.FrontendAPIPort = req.request.port
    ipc.callMain('SetFrontendAPIPort', req.request.port);
    feapiconsumer.Initialise();
    callback(null, {});
}
function DeliverAmbients(req, callback) {
    vuexStore.dispatch('setAmbientBoards', req.request.Boards);
    callback(null, {});
}
// /**
//  * Implements the SayHello RPC method.
//  */
// function sayHello(call: any, callback: any) {
//   console.log("sayhello called from the server!")
//   callback(null, { message: 'Hello ' + call.request.name })
// }
// I think the callback is used to respond to the request.
//# sourceMappingURL=clapiserver.js.map