// Services > Frontend API Consumer

// This service is the API through which the client accesses data in the frontend.

export { } // This says this file is a module, not a script.

// Imports
const grpc = require('grpc')
// const resolve = require('path').resolve
var ipc = require('../../../../node_modules/electron-better-ipc')

// Consts
// const proto = grpc.load({
//   file: 'feapi/feapi.proto',
//   root: resolve(__dirname, '../protos')
// }).feapi

var pmessages = require('../../../../../protos/feapi/feapi_pb.js')
// var feobjmessages = require('../../../../../protos/feobjects/feobjects_pb.js');
var proto = require('../../../../../protos/feapi/feapi_grpc_pb')

let feAPIConsumer: any
let Initialised: boolean

let ExportedMethods = {
  async Initialise() {
    console.log('init is called')
    let feapiport = await ipc.callMain('GetFrontendAPIPort')
    feAPIConsumer = new proto.FrontendAPIClient('127.0.0.1:' + feapiport, grpc.credentials.createInsecure())
    console.log(feAPIConsumer)
    let clapiserverport = await ipc.callMain('GetClientAPIServerPort')
    await ExportedMethods.SetClientAPIServerPort(clapiserverport)
    ipc.callMain('SetFrontendClientConnInitialised', true)
    Initialised = true
  },
  GetAllBoards(callback: any) {
    WaitUntilFrontendReady(async function() {
      console.log("get all boards is making a call")
      console.log('initstate: ', Initialised)
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      let req = new pmessages.AllBoardsRequest
      feAPIConsumer.getAllBoards(req, function(err: any, response: any) {
        if (err) {
          console.log(err)
        } else {
          callback(response.toObject().allboardsList)
        }
      })
    })
  },
  SetClientAPIServerPort(clientAPIServerPort: number) {
    WaitUntilFrontendReady(async function() {
      console.log('clapiserverport mapping is triggered. initstate: ', Initialised)
      // if (!Initialised) {
      //   await ExportedMethods.Initialise()
      // }
      let req = new pmessages.SetClientAPIServerPortRequest()
      req.setPort(clientAPIServerPort)
      console.log(req)
      feAPIConsumer.setClientAPIServerPort(req, function(err: any, response: any) {
        if (err) {
          console.log(err)
        } else {
          console.log(response)
        }
      })
    })
  },
  GetBoardAndThreads(boardfp: string, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('GetBoardsAndThread triggered.')
      let req = new pmessages.BoardAndThreadsRequest
      req.setBoardfingerprint(boardfp)
      feAPIConsumer.getBoardAndThreads(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          console.log(resp.toObject())
          callback(resp.toObject())
        }
      })
    })
  },
  GetThreadAndPosts(boardfp: string, threadfp: string, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('GetThreadAndPosts triggered.')
      let req = new pmessages.ThreadAndPostsRequest
      req.setBoardfingerprint(boardfp)
      req.setThreadfingerprint(threadfp)
      feAPIConsumer.getThreadAndPosts(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  SetBoardSignal(fp: string, subbed: boolean, notify: boolean, lastseen: number, lastSeenOnly: boolean, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SetBoardSignal triggered.')
      let req = new pmessages.BoardSignalRequest
      req.setFingerprint(fp)
      req.setSubscribed(subbed)
      req.setNotify(notify)
      req.setLastseen(lastseen)
      req.setLastseenonly(lastSeenOnly)
      feAPIConsumer.setBoardSignal(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  GetUserAndGraph(fp: string, userEntityRequested: boolean, boardsRequested: boolean, threadsRequested: boolean, postsRequested: boolean, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('GetUserAndGraph triggered.')
      let req = new pmessages.UserAndGraphRequest
      req.setFingerprint(fp)
      req.setUserentityrequested(userEntityRequested)
      req.setUserboardsrequested(boardsRequested)
      req.setUserthreadsrequested(threadsRequested)
      req.setUserpostsrequested(postsRequested)
      console.log(req)
      feAPIConsumer.getUserAndGraph(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  GetUncompiledEntityByKey(entityType: string, ownerfp: string, limit: number, offset: number, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('GetUncompiledEntityByKey triggered.')
      let req = new pmessages.UncompiledEntityByKeyRequest
      if (entityType === 'Board') {
        req.setEntitytype(pmessages.UncompiledEntityType.BOARD)
      }
      if (entityType === 'Thread') {
        req.setEntitytype(pmessages.UncompiledEntityType.THREAD)
      }
      if (entityType === 'Post') {
        req.setEntitytype(pmessages.UncompiledEntityType.POST)
      }
      if (entityType === 'Vote') {
        req.setEntitytype(pmessages.UncompiledEntityType.VOTE)
      }
      if (entityType === 'Key') {
        req.setEntitytype(pmessages.UncompiledEntityType.KEY)
      }
      if (entityType === 'Truststate') {
        req.setEntitytype(pmessages.UncompiledEntityType.TRUSTSTATE)
      }
      req.setLimit(limit)
      req.setOffset(offset)
      req.setOwnerfingerprint(ownerfp)
      console.log(req)
      feAPIConsumer.getUncompiledEntityByKey(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  SendInflightsPruneRequest(callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SendInflightsPruneRequest triggered.')
      let req = new pmessages.InflightsPruneRequest
      feAPIConsumer.sendInflightsPruneRequest(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  IsInitialised(): boolean {
    return Initialised
  },

  /*----------  Methods for user signal actions  ----------*/

  /*
    Important thing here. We do have a RETRACT defined, but this is not defined for anything that goes into a bloom filter. Which means, you cannot retract an upvote, but you can downvote to reverse it. The reason why is that upvotes and downvotes (and elects) are aggregated, therefore after they get added to the bloom filter, we only know probabilistically that they're there. We have two bloom filters for each so we can have a +1 and -1, but adding 0 means adding another bloom filter in. Depending on the demand for a retract we can add a third bloom to the implementation to keep tracking of that, but bloom filters are very expensive because they're per-entity, and we have a lot of entities.

    This does not apply to non-aggregated signals like reporting to mod, those are kept instact and individual, and they can be retracted.
  */
  Upvote(this: any, targetfp: string, priorfp: string, boardfp: string, threadfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp, priorfp,
      'ADDS_TO_DISCUSSION', 'UPVOTE',
      'CONTENT', '', boardfp, threadfp, callback)
  },

  Downvote(this: any, targetfp: string, priorfp: string, boardfp: string, threadfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp, priorfp,
      'ADDS_TO_DISCUSSION', 'DOWNVOTE',
      'CONTENT', '', boardfp, threadfp, callback)
  },

  ReportToMod(this: any, targetfp: string, priorfp: string, reason: string, boardfp: string, threadfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp, priorfp,
      'FOLLOWS_GUIDELINES', 'REPORT_TO_MOD',
      'CONTENT', reason, boardfp, threadfp, callback)
  },

  ModBlock(this: any, targetfp: string, priorfp: string, reason: string, boardfp: string, threadfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp, priorfp,
      'MOD_ACTIONS', 'MODBLOCK',
      'CONTENT', reason, boardfp, threadfp, callback)
  },

  ModApprove(this: any, targetfp: string, priorfp: string, reason: string, boardfp: string, threadfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp, priorfp,
      'MOD_ACTIONS', 'MODAPPROVE',
      'CONTENT', reason, boardfp, threadfp, callback)
  },

  Follow(this: any, targetfp: string, priorfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp, priorfp,
      'PUBLIC_TRUST', 'FOLLOW',
      'USER', '', '', '', callback)
  },

  Block(this: any, targetfp: string, priorfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp, priorfp,
      'PUBLIC_TRUST', 'BLOCK',
      'USER', '', '', '', callback)
  },

  Elect(this: any, targetfp: string, priorfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp, priorfp,
      'PUBLIC_ELECT', 'ELECT',
      'USER', '', '', '', callback)
  },
  Disqualify(this: any, targetfp: string, priorfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp, priorfp,
      'PUBLIC_ELECT', 'DISQUALIFY',
      'USER', '', '', '', callback)
  },

  /*----------  Base signal event action.  ----------*/

  sendSignalEvent(targetfp: string, priorfp: string, typeclass: string, typ: string, targettype: string, signaltext: string, boardfp: string, threadfp: string, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('Send Signal Event base triggered.')
      let now = Math.floor(Date.now() / 1000)
      let req = new pmessages.SignalEventPayload
      let e = new pmessages.Event
      var localUser = require('../../store/index').default.state.localUser
      // ^ Only import when needed and only the specific part. Because vuexstore is also importing this feapi - we don't want it being imported at the beginning to prevent vuexstore from loading feapi.
      e.setOwnerfingerprint(localUser.fingerprint)
      e.setPriorfingerprint(priorfp)
      e.setEventtype(priorfp.length === 0 ? pmessages.EventType.CREATE : pmessages.EventType.UPDATE)
      e.setTimestamp(now)
      req.setEvent(e)
      req.setSignaltargettype(pmessages.SignalTargetType[targettype])
      if (targettype === 'CONTENT') {
        req.setTargetboard(boardfp)
        req.setTargetthread(threadfp)
      }
      if (targettype === 'USER') {
        req.setDomain() // todo
      }
      req.setTargetfingerprint(targetfp)
      req.setSignaltypeclass(pmessages.SignalTypeClass[typeclass])
      req.setSignaltext(signaltext)
      console.log('signal type:')
      console.log(pmessages.SignalType[typ])
      req.setSignaltype(pmessages.SignalType[typ])
      feAPIConsumer.sendSignalEvent(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp)
        }
      })
    })
  },

  /*----------  Methods for content event actions  ----------*/

  /*
    These are things like creating or editing entities that the user has created. If a priorfp is provided, it is an update. If not, it is a create.
  */

  SendBoardContent(this: any, priorfp: string, boarddata: any, callback: any) {
    this.sendContentEvent(priorfp, boarddata, undefined, undefined, undefined, callback)
  },

  SendThreadContent(this: any, priorfp: string, threaddata: any, callback: any) {
    this.sendContentEvent(priorfp, undefined, threaddata, undefined, undefined, callback)
  },

  SendPostContent(this: any, priorfp: string, postdata: any, callback: any) {
    this.sendContentEvent(priorfp, undefined, undefined, postdata, undefined, callback)
  },

  SendUserContent(this: any, priorfp: string, userdata: any, callback: any) {
    this.sendContentEvent(priorfp, undefined, undefined, undefined, userdata, callback)
  },

  /*----------  Base content event action.  ----------*/

  sendContentEvent(priorfp: string, boarddata: any, threaddata: any, postdata: any, userdata: any, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('Send Content Event base triggered.')
      let now = Math.floor(Date.now() / 1000)
      let req = new pmessages.ContentEventPayload
      let e = new pmessages.Event
      var localUser = require('../../store/index').default.state.localUser
      // ^ Only import when needed and only the specific part. Because vuexstore is also importing this feapi - we don't want it being imported at the beginning to prevent vuexstore from loading feapi.
      let globalMethods = require('../globals/methods')
      if (!globalMethods.IsUndefined(localUser)) {
        e.setOwnerfingerprint(localUser.fingerprint)
      }
      e.setPriorfingerprint(priorfp)
      e.setEventtype(priorfp.length === 0 ? pmessages.EventType.CREATE : pmessages.EventType.UPDATE)
      e.setTimestamp(now)
      req.setEvent(e)
      req.setBoarddata(boarddata)
      req.setThreaddata(threaddata)
      req.setPostdata(postdata)
      req.setKeydata(userdata)
      feAPIConsumer.sendContentEvent(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp)
        }
      })
    })
  },
}
module.exports = ExportedMethods

function WaitUntilFrontendReady(cb: any): any {
  async function check() {
    let initialised = await ipc.callMain('GetFrontendClientConnInitialised')
    // console.log(initialised)
    if (!initialised) {
      // console.log("Frontend still not ready, waiting a little more...")
      return setTimeout(check, 333)
    } else {
      return cb()
    }
  }
  return check()
}