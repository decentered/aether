// GENERATED CODE -- DO NOT EDIT!

// Original file comments:
// Frontend API server Protobufs
//
'use strict';
var grpc = require('grpc');
var feapi_feapi_pb = require('../feapi/feapi_pb.js');
var feobjects_feobjects_pb = require('../feobjects/feobjects_pb.js');
var mimapi_structprotos_pb = require('../mimapi/structprotos_pb.js');

function serialize_feapi_AllBoardsRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.AllBoardsRequest)) {
    throw new Error('Expected argument of type feapi.AllBoardsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_AllBoardsRequest(buffer_arg) {
  return feapi_feapi_pb.AllBoardsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_AllBoardsResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.AllBoardsResponse)) {
    throw new Error('Expected argument of type feapi.AllBoardsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_AllBoardsResponse(buffer_arg) {
  return feapi_feapi_pb.AllBoardsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BEReadyRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.BEReadyRequest)) {
    throw new Error('Expected argument of type feapi.BEReadyRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BEReadyRequest(buffer_arg) {
  return feapi_feapi_pb.BEReadyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BEReadyResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.BEReadyResponse)) {
    throw new Error('Expected argument of type feapi.BEReadyResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BEReadyResponse(buffer_arg) {
  return feapi_feapi_pb.BEReadyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardAndThreadsRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardAndThreadsRequest)) {
    throw new Error('Expected argument of type feapi.BoardAndThreadsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardAndThreadsRequest(buffer_arg) {
  return feapi_feapi_pb.BoardAndThreadsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardAndThreadsResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardAndThreadsResponse)) {
    throw new Error('Expected argument of type feapi.BoardAndThreadsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardAndThreadsResponse(buffer_arg) {
  return feapi_feapi_pb.BoardAndThreadsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardSignalRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardSignalRequest)) {
    throw new Error('Expected argument of type feapi.BoardSignalRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardSignalRequest(buffer_arg) {
  return feapi_feapi_pb.BoardSignalRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardSignalResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardSignalResponse)) {
    throw new Error('Expected argument of type feapi.BoardSignalResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardSignalResponse(buffer_arg) {
  return feapi_feapi_pb.BoardSignalResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ContentEventPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.ContentEventPayload)) {
    throw new Error('Expected argument of type feapi.ContentEventPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ContentEventPayload(buffer_arg) {
  return feapi_feapi_pb.ContentEventPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ContentEventResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.ContentEventResponse)) {
    throw new Error('Expected argument of type feapi.ContentEventResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ContentEventResponse(buffer_arg) {
  return feapi_feapi_pb.ContentEventResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_InflightsPruneRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.InflightsPruneRequest)) {
    throw new Error('Expected argument of type feapi.InflightsPruneRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_InflightsPruneRequest(buffer_arg) {
  return feapi_feapi_pb.InflightsPruneRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_InflightsPruneResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.InflightsPruneResponse)) {
    throw new Error('Expected argument of type feapi.InflightsPruneResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_InflightsPruneResponse(buffer_arg) {
  return feapi_feapi_pb.InflightsPruneResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SetClientAPIServerPortRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.SetClientAPIServerPortRequest)) {
    throw new Error('Expected argument of type feapi.SetClientAPIServerPortRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SetClientAPIServerPortRequest(buffer_arg) {
  return feapi_feapi_pb.SetClientAPIServerPortRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SetClientAPIServerPortResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.SetClientAPIServerPortResponse)) {
    throw new Error('Expected argument of type feapi.SetClientAPIServerPortResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SetClientAPIServerPortResponse(buffer_arg) {
  return feapi_feapi_pb.SetClientAPIServerPortResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SignalEventPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.SignalEventPayload)) {
    throw new Error('Expected argument of type feapi.SignalEventPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SignalEventPayload(buffer_arg) {
  return feapi_feapi_pb.SignalEventPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SignalEventResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.SignalEventResponse)) {
    throw new Error('Expected argument of type feapi.SignalEventResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SignalEventResponse(buffer_arg) {
  return feapi_feapi_pb.SignalEventResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ThreadAndPostsRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.ThreadAndPostsRequest)) {
    throw new Error('Expected argument of type feapi.ThreadAndPostsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ThreadAndPostsRequest(buffer_arg) {
  return feapi_feapi_pb.ThreadAndPostsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ThreadAndPostsResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.ThreadAndPostsResponse)) {
    throw new Error('Expected argument of type feapi.ThreadAndPostsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ThreadAndPostsResponse(buffer_arg) {
  return feapi_feapi_pb.ThreadAndPostsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_UncompiledEntityByKeyRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.UncompiledEntityByKeyRequest)) {
    throw new Error('Expected argument of type feapi.UncompiledEntityByKeyRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_UncompiledEntityByKeyRequest(buffer_arg) {
  return feapi_feapi_pb.UncompiledEntityByKeyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_UncompiledEntityByKeyResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.UncompiledEntityByKeyResponse)) {
    throw new Error('Expected argument of type feapi.UncompiledEntityByKeyResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_UncompiledEntityByKeyResponse(buffer_arg) {
  return feapi_feapi_pb.UncompiledEntityByKeyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_UserAndGraphRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.UserAndGraphRequest)) {
    throw new Error('Expected argument of type feapi.UserAndGraphRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_UserAndGraphRequest(buffer_arg) {
  return feapi_feapi_pb.UserAndGraphRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_UserAndGraphResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.UserAndGraphResponse)) {
    throw new Error('Expected argument of type feapi.UserAndGraphResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_UserAndGraphResponse(buffer_arg) {
  return feapi_feapi_pb.UserAndGraphResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// These "Set", "Get" verbs are written from the viewpoint of the consumer of this api.
var FrontendAPIService = exports.FrontendAPIService = {
  backendReady: {
    path: '/feapi.FrontendAPI/BackendReady',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.BEReadyRequest,
    responseType: feapi_feapi_pb.BEReadyResponse,
    requestSerialize: serialize_feapi_BEReadyRequest,
    requestDeserialize: deserialize_feapi_BEReadyRequest,
    responseSerialize: serialize_feapi_BEReadyResponse,
    responseDeserialize: deserialize_feapi_BEReadyResponse,
  },
  setClientAPIServerPort: {
    path: '/feapi.FrontendAPI/SetClientAPIServerPort',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.SetClientAPIServerPortRequest,
    responseType: feapi_feapi_pb.SetClientAPIServerPortResponse,
    requestSerialize: serialize_feapi_SetClientAPIServerPortRequest,
    requestDeserialize: deserialize_feapi_SetClientAPIServerPortRequest,
    responseSerialize: serialize_feapi_SetClientAPIServerPortResponse,
    responseDeserialize: deserialize_feapi_SetClientAPIServerPortResponse,
  },
  getThreadAndPosts: {
    path: '/feapi.FrontendAPI/GetThreadAndPosts',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.ThreadAndPostsRequest,
    responseType: feapi_feapi_pb.ThreadAndPostsResponse,
    requestSerialize: serialize_feapi_ThreadAndPostsRequest,
    requestDeserialize: deserialize_feapi_ThreadAndPostsRequest,
    responseSerialize: serialize_feapi_ThreadAndPostsResponse,
    responseDeserialize: deserialize_feapi_ThreadAndPostsResponse,
  },
  getBoardAndThreads: {
    path: '/feapi.FrontendAPI/GetBoardAndThreads',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.BoardAndThreadsRequest,
    responseType: feapi_feapi_pb.BoardAndThreadsResponse,
    requestSerialize: serialize_feapi_BoardAndThreadsRequest,
    requestDeserialize: deserialize_feapi_BoardAndThreadsRequest,
    responseSerialize: serialize_feapi_BoardAndThreadsResponse,
    responseDeserialize: deserialize_feapi_BoardAndThreadsResponse,
  },
  getAllBoards: {
    path: '/feapi.FrontendAPI/GetAllBoards',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.AllBoardsRequest,
    responseType: feapi_feapi_pb.AllBoardsResponse,
    requestSerialize: serialize_feapi_AllBoardsRequest,
    requestDeserialize: deserialize_feapi_AllBoardsRequest,
    responseSerialize: serialize_feapi_AllBoardsResponse,
    responseDeserialize: deserialize_feapi_AllBoardsResponse,
  },
  setBoardSignal: {
    path: '/feapi.FrontendAPI/SetBoardSignal',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.BoardSignalRequest,
    responseType: feapi_feapi_pb.BoardSignalResponse,
    requestSerialize: serialize_feapi_BoardSignalRequest,
    requestDeserialize: deserialize_feapi_BoardSignalRequest,
    responseSerialize: serialize_feapi_BoardSignalResponse,
    responseDeserialize: deserialize_feapi_BoardSignalResponse,
  },
  getUserAndGraph: {
    path: '/feapi.FrontendAPI/GetUserAndGraph',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.UserAndGraphRequest,
    responseType: feapi_feapi_pb.UserAndGraphResponse,
    requestSerialize: serialize_feapi_UserAndGraphRequest,
    requestDeserialize: deserialize_feapi_UserAndGraphRequest,
    responseSerialize: serialize_feapi_UserAndGraphResponse,
    responseDeserialize: deserialize_feapi_UserAndGraphResponse,
  },
  sendContentEvent: {
    path: '/feapi.FrontendAPI/SendContentEvent',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.ContentEventPayload,
    responseType: feapi_feapi_pb.ContentEventResponse,
    requestSerialize: serialize_feapi_ContentEventPayload,
    requestDeserialize: deserialize_feapi_ContentEventPayload,
    responseSerialize: serialize_feapi_ContentEventResponse,
    responseDeserialize: deserialize_feapi_ContentEventResponse,
  },
  sendSignalEvent: {
    path: '/feapi.FrontendAPI/SendSignalEvent',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.SignalEventPayload,
    responseType: feapi_feapi_pb.SignalEventResponse,
    requestSerialize: serialize_feapi_SignalEventPayload,
    requestDeserialize: deserialize_feapi_SignalEventPayload,
    responseSerialize: serialize_feapi_SignalEventResponse,
    responseDeserialize: deserialize_feapi_SignalEventResponse,
  },
  getUncompiledEntityByKey: {
    path: '/feapi.FrontendAPI/GetUncompiledEntityByKey',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.UncompiledEntityByKeyRequest,
    responseType: feapi_feapi_pb.UncompiledEntityByKeyResponse,
    requestSerialize: serialize_feapi_UncompiledEntityByKeyRequest,
    requestDeserialize: deserialize_feapi_UncompiledEntityByKeyRequest,
    responseSerialize: serialize_feapi_UncompiledEntityByKeyResponse,
    responseDeserialize: deserialize_feapi_UncompiledEntityByKeyResponse,
  },
  sendInflightsPruneRequest: {
    path: '/feapi.FrontendAPI/SendInflightsPruneRequest',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.InflightsPruneRequest,
    responseType: feapi_feapi_pb.InflightsPruneResponse,
    requestSerialize: serialize_feapi_InflightsPruneRequest,
    requestDeserialize: deserialize_feapi_InflightsPruneRequest,
    responseSerialize: serialize_feapi_InflightsPruneResponse,
    responseDeserialize: deserialize_feapi_InflightsPruneResponse,
  },
};

exports.FrontendAPIClient = grpc.makeGenericClientConstructor(FrontendAPIService);
