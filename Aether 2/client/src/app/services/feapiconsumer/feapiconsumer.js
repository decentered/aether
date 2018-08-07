"use strict";
// Services > Frontend API Consumer
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = y[op[0] & 2 ? "return" : op[0] ? "throw" : "next"]) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [0, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
// Imports
var grpc = require('grpc');
// const resolve = require('path').resolve
// let globals = require('../globals/globals')
var ipc = require('../../../../node_modules/electron-better-ipc');
// Consts
// const proto = grpc.load({
//   file: 'feapi/feapi.proto',
//   root: resolve(__dirname, '../protos')
// }).feapi
var pmessages = require('../../../../../protos/feapi/feapi_pb.js');
// var feobjmessages = require('../../../../../protos/feobjects/feobjects_pb.js');
var proto = require('../../../../../protos/feapi/feapi_grpc_pb');
var feAPIConsumer;
var Initialised;
var ExportedMethods = {
    Initialise: function () {
        return __awaiter(this, void 0, void 0, function () {
            var feapiport, clapiserverport;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        console.log('init is called');
                        return [4 /*yield*/, ipc.callMain('GetFrontendAPIPort')];
                    case 1:
                        feapiport = _a.sent();
                        feAPIConsumer = new proto.FrontendAPIClient('127.0.0.1:' + feapiport, grpc.credentials.createInsecure());
                        console.log(feAPIConsumer);
                        return [4 /*yield*/, ipc.callMain('GetClientAPIServerPort')];
                    case 2:
                        clapiserverport = _a.sent();
                        return [4 /*yield*/, ExportedMethods.SetClientAPIServerPort(clapiserverport)];
                    case 3:
                        _a.sent();
                        ipc.callMain('SetFrontendClientConnInitialised', true);
                        Initialised = true;
                        return [2 /*return*/];
                }
            });
        });
    },
    GetAllBoards: function (callback) {
        WaitUntilFrontendReady(function () {
            return __awaiter(this, void 0, void 0, function () {
                var req;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            console.log("get all boards is making a call");
                            console.log('initstate: ', Initialised);
                            if (!!Initialised) return [3 /*break*/, 2];
                            return [4 /*yield*/, ExportedMethods.Initialise()];
                        case 1:
                            _a.sent();
                            _a.label = 2;
                        case 2:
                            req = new pmessages.AllBoardsRequest;
                            feAPIConsumer.getAllBoards(req, function (err, response) {
                                if (err) {
                                    console.log(err);
                                }
                                else {
                                    callback(response.toObject().allboardsList);
                                }
                            });
                            return [2 /*return*/];
                    }
                });
            });
        });
    },
    SetClientAPIServerPort: function (clientAPIServerPort) {
        WaitUntilFrontendReady(function () {
            return __awaiter(this, void 0, void 0, function () {
                var req;
                return __generator(this, function (_a) {
                    console.log('clapiserverport mapping is triggered. initstate: ', Initialised);
                    req = new pmessages.SetClientAPIServerPortRequest();
                    req.setPort(clientAPIServerPort);
                    console.log(req);
                    feAPIConsumer.setClientAPIServerPort(req, function (err, response) {
                        if (err) {
                            console.log(err);
                        }
                        else {
                            console.log(response);
                        }
                    });
                    return [2 /*return*/];
                });
            });
        });
    },
    GetBoardAndThreads: function (boardfp, callback) {
        WaitUntilFrontendReady(function () {
            return __awaiter(this, void 0, void 0, function () {
                var req;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            if (!!Initialised) return [3 /*break*/, 2];
                            return [4 /*yield*/, ExportedMethods.Initialise()];
                        case 1:
                            _a.sent();
                            _a.label = 2;
                        case 2:
                            console.log('GetBoardsAndThread triggered.');
                            req = new pmessages.BoardAndThreadsRequest;
                            req.setBoardfingerprint(boardfp);
                            feAPIConsumer.getBoardAndThreads(req, function (err, resp) {
                                if (err) {
                                    console.log(err);
                                }
                                else {
                                    console.log(resp.toObject());
                                    callback(resp.toObject());
                                }
                            });
                            return [2 /*return*/];
                    }
                });
            });
        });
    },
    GetThreadAndPosts: function (boardfp, threadfp, callback) {
        WaitUntilFrontendReady(function () {
            return __awaiter(this, void 0, void 0, function () {
                var req;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            if (!!Initialised) return [3 /*break*/, 2];
                            return [4 /*yield*/, ExportedMethods.Initialise()];
                        case 1:
                            _a.sent();
                            _a.label = 2;
                        case 2:
                            console.log('GetThreadAndPosts triggered.');
                            req = new pmessages.ThreadAndPostsRequest;
                            req.setBoardfingerprint(boardfp);
                            req.setThreadfingerprint(threadfp);
                            feAPIConsumer.getThreadAndPosts(req, function (err, resp) {
                                if (err) {
                                    console.log(err);
                                }
                                else {
                                    callback(resp.toObject());
                                }
                            });
                            return [2 /*return*/];
                    }
                });
            });
        });
    },
    SetBoardSignal: function (fp, subbed, notify, lastseen, lastSeenOnly, callback) {
        WaitUntilFrontendReady(function () {
            return __awaiter(this, void 0, void 0, function () {
                var req;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            if (!!Initialised) return [3 /*break*/, 2];
                            return [4 /*yield*/, ExportedMethods.Initialise()];
                        case 1:
                            _a.sent();
                            _a.label = 2;
                        case 2:
                            console.log('SetBoardSignal triggered.');
                            req = new pmessages.BoardSignalRequest;
                            req.setFingerprint(fp);
                            req.setSubscribed(subbed);
                            req.setNotify(notify);
                            req.setLastseen(lastseen);
                            req.setLastseenonly(lastSeenOnly);
                            feAPIConsumer.setBoardSignal(req, function (err, resp) {
                                if (err) {
                                    console.log(err);
                                }
                                else {
                                    callback(resp.toObject());
                                }
                            });
                            return [2 /*return*/];
                    }
                });
            });
        });
    },
    GetUserAndGraph: function (fp, userEntityRequested, boardsRequested, threadsRequested, postsRequested, callback) {
        WaitUntilFrontendReady(function () {
            return __awaiter(this, void 0, void 0, function () {
                var req;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            if (!!Initialised) return [3 /*break*/, 2];
                            return [4 /*yield*/, ExportedMethods.Initialise()];
                        case 1:
                            _a.sent();
                            _a.label = 2;
                        case 2:
                            console.log('GetUserAndGraph triggered.');
                            req = new pmessages.UserAndGraphRequest;
                            req.setFingerprint(fp);
                            req.setUserentityrequested(userEntityRequested);
                            req.setUserboardsrequested(boardsRequested);
                            req.setUserthreadsrequested(threadsRequested);
                            req.setUserpostsrequested(postsRequested);
                            feAPIConsumer.getUserAndGraph(req, function (err, resp) {
                                if (err) {
                                    console.log(err);
                                }
                                else {
                                    callback(resp.toObject());
                                }
                            });
                            return [2 /*return*/];
                    }
                });
            });
        });
    },
    IsInitialised: function () {
        return Initialised;
    },
    /*----------  Methods for user signal actions  ----------*/
    /*
      Important thing here. We do have a RETRACT defined, but this is not defined for anything that goes into a bloom filter. Which means, you cannot retract an upvote, but you can downvote to reverse it. The reason why is that upvotes and downvotes (and elects) are aggregated, therefore after they get added to the bloom filter, we only know probabilistically that they're there. We have two bloom filters for each so we can have a +1 and -1, but adding 0 means adding another bloom filter in. Depending on the demand for a retract we can add a third bloom to the implementation to keep tracking of that, but bloom filters are very expensive because they're per-entity, and we have a lot of entities.
  
      This does not apply to non-aggregated signals like reporting to mod, those are kept instact and individual, and they can be retracted.
    */
    Upvote: function (targetfp, priorfp, callback) {
        this.sendSignalEvent(targetfp, priorfp, 'ADDS_TO_DISCUSSION', 'UPVOTE', 'CONTENT', callback);
    },
    Downvote: function (targetfp, priorfp, callback) {
        this.sendSignalEvent(targetfp, priorfp, 'ADDS_TO_DISCUSSION', 'DOWNVOTE', 'CONTENT', callback);
    },
    ReportToMod: function (targetfp, priorfp, callback) {
        this.sendSignalEvent(targetfp, priorfp, 'FOLLOWS_GUIDELINES', 'REPORT_TO_MOD', 'CONTENT', callback);
    },
    ModBlock: function (targetfp, priorfp, callback) {
        this.sendSignalEvent(targetfp, priorfp, 'MOD_ACTIONS', 'MODBLOCK', 'CONTENT', callback);
    },
    ModApprove: function (targetfp, priorfp, callback) {
        this.sendSignalEvent(targetfp, priorfp, 'MOD_ACTIONS', 'MODAPPROVE', 'CONTENT', callback);
    },
    Follow: function (targetfp, priorfp, callback) {
        this.sendSignalEvent(targetfp, priorfp, 'PUBLIC_TRUST', 'FOLLOW', 'USER', callback);
    },
    Block: function (targetfp, priorfp, callback) {
        this.sendSignalEvent(targetfp, priorfp, 'PUBLIC_TRUST', 'BLOCK', 'USER', callback);
    },
    Elect: function (targetfp, priorfp, callback) {
        this.sendSignalEvent(targetfp, priorfp, 'PUBLIC_ELECT', 'ELECT', 'USER', callback);
    },
    Disqualify: function (targetfp, priorfp, callback) {
        this.sendSignalEvent(targetfp, priorfp, 'PUBLIC_ELECT', 'DISQUALIFY', 'USER', callback);
    },
    /*----------  Base signal event action.  ----------*/
    sendSignalEvent: function (targetfp, priorfp, typeclass, typ, targettype, callback) {
        WaitUntilFrontendReady(function () {
            return __awaiter(this, void 0, void 0, function () {
                var now, req, e;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            if (!!Initialised) return [3 /*break*/, 2];
                            return [4 /*yield*/, ExportedMethods.Initialise()];
                        case 1:
                            _a.sent();
                            _a.label = 2;
                        case 2:
                            // let boarddata2 = new feobjmessages.CompiledBoardEntity
                            console.log('Send Signal Event base triggered.');
                            now = Math.floor(Date.now() / 1000);
                            req = new pmessages.SignalEventRequest;
                            e = new pmessages.Event;
                            e.setOwnerfingerprint("TODO HERE COMES LOCAL USER FP");
                            e.setPriorfingerprint(priorfp);
                            e.setEventtype(priorfp.length === 0 ? pmessages.EventType.CREATE : pmessages.EventType.EDIT);
                            e.setTimestamp(now);
                            req.setEvent(e);
                            req.setSignaltargettype(pmessages.SignalTargetType[targettype]);
                            req.setTargetfingerprint(targetfp);
                            req.setSignaltypeclass(pmessages.SignalTypeClass[typeclass]);
                            console.log(pmessages.SignalType[typ]);
                            req.setSignaltype(pmessages.SignalType[typ]);
                            feAPIConsumer.sendSignalEvent(req, function (err, resp) {
                                if (err) {
                                    console.log(err);
                                }
                                else {
                                    callback(resp);
                                }
                            });
                            return [2 /*return*/];
                    }
                });
            });
        });
    },
    /*----------  Methods for content event actions  ----------*/
    /*
      These are things like creating or editing entities that the user has created. If a priorfp is provided, it is an update. If not, it is a create.
    */
    BoardContentEvent: function (priorfp, boarddata, callback) {
        this.sendContentEvent(priorfp, boarddata, undefined, undefined, undefined, callback);
    },
    ThreadContentEvent: function (priorfp, threaddata, callback) {
        this.sendContentEvent(priorfp, undefined, threaddata, undefined, undefined, callback);
    },
    PostContentEvent: function (priorfp, postdata, callback) {
        this.sendContentEvent(priorfp, undefined, undefined, postdata, undefined, callback);
    },
    UserContentEvent: function (priorfp, userdata, callback) {
        this.sendContentEvent(priorfp, undefined, undefined, undefined, userdata, callback);
    },
    /*----------  Base content event action.  ----------*/
    sendContentEvent: function (priorfp, boarddata, threaddata, postdata, userdata, callback) {
        WaitUntilFrontendReady(function () {
            return __awaiter(this, void 0, void 0, function () {
                var now, req, e;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            if (!!Initialised) return [3 /*break*/, 2];
                            return [4 /*yield*/, ExportedMethods.Initialise()];
                        case 1:
                            _a.sent();
                            _a.label = 2;
                        case 2:
                            console.log('Send Content Event base triggered.');
                            now = Math.floor(Date.now() / 1000);
                            req = new pmessages.ContentEventRequest;
                            e = new pmessages.Event;
                            e.setOwnerfingerprint("TODO HERE COMES LOCAL USER FP");
                            e.setPriorfingerprint(priorfp);
                            e.setEventtype(priorfp.length === 0 ? pmessages.EventType.CREATE : pmessages.EventType.EDIT);
                            e.setTimestamp(now);
                            req.setEvent(e);
                            req.setBoarddata(boarddata);
                            req.setThreaddata(threaddata);
                            req.setPostdata(postdata);
                            req.setUserdata(userdata);
                            feAPIConsumer.sendContentEvent(req, function (err, resp) {
                                if (err) {
                                    console.log(err);
                                }
                                else {
                                    callback(resp);
                                }
                            });
                            return [2 /*return*/];
                    }
                });
            });
        });
    },
};
module.exports = ExportedMethods;
function WaitUntilFrontendReady(cb) {
    function check() {
        return __awaiter(this, void 0, void 0, function () {
            var initialised;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, ipc.callMain('GetFrontendClientConnInitialised')
                        // console.log(initialised)
                    ];
                    case 1:
                        initialised = _a.sent();
                        // console.log(initialised)
                        if (!initialised) {
                            // console.log("Frontend still not ready, waiting a little more...")
                            return [2 /*return*/, setTimeout(check, 333)];
                        }
                        else {
                            return [2 /*return*/, cb()];
                        }
                        return [2 /*return*/];
                }
            });
        });
    }
    return check();
}
//# sourceMappingURL=feapiconsumer.js.map