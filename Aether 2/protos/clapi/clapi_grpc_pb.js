// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var clapi_clapi_pb = require('../clapi/clapi_pb.js');
var feobjects_feobjects_pb = require('../feobjects/feobjects_pb.js');
var mimapi_structprotos_pb = require('../mimapi/structprotos_pb.js');

function serialize_clapi_AmbientLocalUserEntityPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientLocalUserEntityPayload)) {
    throw new Error('Expected argument of type clapi.AmbientLocalUserEntityPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientLocalUserEntityPayload(buffer_arg) {
  return clapi_clapi_pb.AmbientLocalUserEntityPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientLocalUserEntityResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientLocalUserEntityResponse)) {
    throw new Error('Expected argument of type clapi.AmbientLocalUserEntityResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientLocalUserEntityResponse(buffer_arg) {
  return clapi_clapi_pb.AmbientLocalUserEntityResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientStatusPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientStatusPayload)) {
    throw new Error('Expected argument of type clapi.AmbientStatusPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientStatusPayload(buffer_arg) {
  return clapi_clapi_pb.AmbientStatusPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientStatusResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientStatusResponse)) {
    throw new Error('Expected argument of type clapi.AmbientStatusResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientStatusResponse(buffer_arg) {
  return clapi_clapi_pb.AmbientStatusResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientsRequest(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientsRequest)) {
    throw new Error('Expected argument of type clapi.AmbientsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientsRequest(buffer_arg) {
  return clapi_clapi_pb.AmbientsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientsResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientsResponse)) {
    throw new Error('Expected argument of type clapi.AmbientsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientsResponse(buffer_arg) {
  return clapi_clapi_pb.AmbientsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_FEReadyRequest(arg) {
  if (!(arg instanceof clapi_clapi_pb.FEReadyRequest)) {
    throw new Error('Expected argument of type clapi.FEReadyRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_FEReadyRequest(buffer_arg) {
  return clapi_clapi_pb.FEReadyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_FEReadyResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.FEReadyResponse)) {
    throw new Error('Expected argument of type clapi.FEReadyResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_FEReadyResponse(buffer_arg) {
  return clapi_clapi_pb.FEReadyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// These "Set", "Get" verbs are written from the viewpoint of the consumer of this api.
var ClientAPIService = exports.ClientAPIService = {
  frontendReady: {
    path: '/clapi.ClientAPI/FrontendReady',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.FEReadyRequest,
    responseType: clapi_clapi_pb.FEReadyResponse,
    requestSerialize: serialize_clapi_FEReadyRequest,
    requestDeserialize: deserialize_clapi_FEReadyRequest,
    responseSerialize: serialize_clapi_FEReadyResponse,
    responseDeserialize: deserialize_clapi_FEReadyResponse,
  },
  deliverAmbients: {
    path: '/clapi.ClientAPI/DeliverAmbients',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.AmbientsRequest,
    responseType: clapi_clapi_pb.AmbientsResponse,
    requestSerialize: serialize_clapi_AmbientsRequest,
    requestDeserialize: deserialize_clapi_AmbientsRequest,
    responseSerialize: serialize_clapi_AmbientsResponse,
    responseDeserialize: deserialize_clapi_AmbientsResponse,
  },
  sendAmbientStatus: {
    path: '/clapi.ClientAPI/SendAmbientStatus',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.AmbientStatusPayload,
    responseType: clapi_clapi_pb.AmbientStatusResponse,
    requestSerialize: serialize_clapi_AmbientStatusPayload,
    requestDeserialize: deserialize_clapi_AmbientStatusPayload,
    responseSerialize: serialize_clapi_AmbientStatusResponse,
    responseDeserialize: deserialize_clapi_AmbientStatusResponse,
  },
  sendAmbientLocalUserEntity: {
    path: '/clapi.ClientAPI/SendAmbientLocalUserEntity',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.AmbientLocalUserEntityPayload,
    responseType: clapi_clapi_pb.AmbientLocalUserEntityResponse,
    requestSerialize: serialize_clapi_AmbientLocalUserEntityPayload,
    requestDeserialize: deserialize_clapi_AmbientLocalUserEntityPayload,
    responseSerialize: serialize_clapi_AmbientLocalUserEntityResponse,
    responseDeserialize: deserialize_clapi_AmbientLocalUserEntityResponse,
  },
};

exports.ClientAPIClient = grpc.makeGenericClientConstructor(ClientAPIService);
