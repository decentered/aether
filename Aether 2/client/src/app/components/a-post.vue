<template>
  <div class="post" :class="{'inflight': isInflightEntity}">
    <div class="inflight-box" v-if="isVisibleInflightEntity">
      <a-inflight-info :status="visibleInflightStatus" :refresherFunc="refresh"></a-inflight-info>
    </div>
    <div class="expand-container" @click="toggleCollapse">
      <template v-if="!collapsed">
        <icon name="minus-circle"></icon>
      </template>
      <template v-if="collapsed">
        <icon name="plus-circle"></icon>
      </template>
    </div>
    <div class="meta-container">
      <div class="author">
        <a-username :owner="post.owner"></a-username>
      </div>
      <div class="post-datetime">
        <a-timestamp :creation="postCreation" :lastupdate="postLastUpdate"></a-timestamp>
      </div>
    </div>
    <template v-if="!collapsed">
      <div class="content">
        <div class="content-text-container">
          <a-composer :spec="postEditExistingSpec" v-if="editPaneOpen"></a-composer>
          <a-markdown class="content-text" v-else :content="visiblePost.body"></a-markdown>
          <div class="actions-rack" v-if="actionsVisible">
            <a class="action reply" @click="toggleReplyPane">
            Reply
          </a>
            <a class="action edit" v-show="post.selfcreated" @click="toggleEditPane">
            Edit
          </a>
            <a class="action report" v-show="!post.selfcreated" @click="toggleReportPane">
            Report
          </a>
          </div>
        </div>
        <a-ballot v-if="actionsVisible" :contentsignals="post.compiledcontentsignals" :boardfp="post.board" :threadfp="post.thread"></a-ballot>
      </div>

      <div class="report-container" v-if="reportPaneOpen">
        <a-composer :spec="reportSpec"></a-composer>
      </div>
      <a-composer :spec="postComposerSpec" v-if="replyPaneOpen"></a-composer>
      <a-post v-for="iflChild in inflightChildren" :post="iflChild.entity" :inflightStatus="iflChild.status"></a-post>
      <a-post v-for="child in post.children" :post="child"></a-post>
    </template>
  </div>
</template>

<script lang="ts">
  var globalMethods = require('../services/globals/methods')
  var mixins = require('../mixins/mixins')
  var fe = require('../services/feapiconsumer/feapiconsumer')
  var mimobjs = require('../../../../protos/mimapi/structprotos_pb.js')
  var vuexStore = require('../store/index').default
  export default {
    name: 'a-post',
    mixins: [mixins.localUserMixin],
    props: ['post', 'inflightStatus', 'uncompiled'],
    data(this: any): any {
      return {
        hasUpvoted: false,
        hasDownvoted: false,
        replyPaneOpen: false,
        editPaneOpen: false,
        reportPaneOpen: false,
        collapsed: false,
        postComposerSpec: {
          fields: [{
            id: "postBody",
            visibleName: "",
            description: "",
            placeholder: "Post a reply",
            maxCharCount: 20480,
            heightRows: 5,
            previewDisabled: false,
            content: "",
            optional: false,
          }],
          commitAction: this.submitPost,
          commitActionName: "SUBMIT",
          cancelAction: this.toggleReplyPane,
          cancelActionName: "CANCEL",
          autofocus: true,
        },
        postEditExistingSpec: {
          fields: [{
            id: "postBody",
            visibleName: "",
            description: "",
            placeholder: "",
            maxCharCount: 20480,
            heightRows: 5,
            previewDisabled: false,
            content: "",
            optional: false,
          }],
          commitAction: this.submitEditExistingPost,
          commitActionName: "SAVE",
          cancelAction: this.toggleEditPane,
          cancelActionName: "CANCEL",
          autofocus: false,
        },
        reportSpec: {
          fields: [{
            id: 'reportReason',
            emptyWarningDisabled: true,
            visibleName: "Report to mods",
            description: "This report will go to the mods. (Heads up - the reports are publicly visible to everyone.)",
            placeholder: "What's the reason?",
            maxCharCount: 256,
            heightRows: 1,
            previewDisabled: true,
            content: '',
            optional: false,
          }],
          commitAction: this.submitReport,
          commitActionName: "REPORT",
          cancelAction: this.toggleReportPane,
          cancelActionName: "CANCEL",
          autofocus: true,
        },
      }
    },
    computed: {
      /*
        This requires a little bit of an explanation, since there is a few layers of stuff here.

        First of all, visibleInflightEntity, or visible[anything] is for edits. Effectively, 'post' is the unedited underlying post, and 'visiblePost' is the edit applied on it, if any. Why is this necessary? Because we want to show the user the edits that it has made immediately, alongside a progress bar for showing when the edit will actually commit in.

        This gets a little confusing to read at first because the post itself can also show recursive children, and each of those children also use the same template. But fundamentally it's simple: what we do is if it's an inflight entity, we disable the further content generation anchors (upvote, downvote, report, create, edit ...) until that entity is fully generated and is in place.

        One trick to keep in mind that I'm doing is I'm only replacing the post's edited body, and I keep all other data using the original post. This simplifies a lot of things because our inflight posts don't actually carry most of the stuff the post has, just the delta (so that it's more efficient), so showing the inflight post fully in exchange to the old post would leave a lot of fields missing.
      */
      postCreation(this: any) {
        // These are necessary because in uncompiled entities, these are in thread.provable.creation, but in compiled ones it's in thread.creation.
        if (this.uncompiled || this.isInflightEntity) {
          if (typeof this.post.provable === 'undefined') {
            return 0
          }
          return this.post.provable.creation
        }
        return this.post.creation
      },
      postLastUpdate(this: any) {
        if (this.uncompiled || this.isInflightEntity) {
          if (typeof this.post.updateable === 'undefined') {
            return 0
          }
          return this.post.updateable.lastupdate
        }
        return this.post.lastupdate
      },
      postFingerprint(this: any) {
        if (this.uncompiled) {
          return this.post.provable.fingerprint
        }
        return this.post.fingerprint
      },
      isVisibleInflightEntity(this: any): boolean {
        if (typeof this.visibleInflightStatus !== 'undefined') {
          return true
        }
        return false
      },
      isInflightEntity(this: any): boolean {
        if (typeof this.inflightStatus !== 'undefined') {
          return true
        }
        return false
      },
      inflightChildren(this: any): any[] {
        let iflChildren: any[] = []
        for (let val of this.$store.state.ambientStatus.inflights.postsList) {
          if (this.post.fingerprint !== val.entity.parent) {
            continue
          }
          iflChildren.push(val)
        }
        return iflChildren
      },
      inflightEdits(this: any): any[] {
        /*This one looks for whether the current post itself was edited, and if so, lists all of them. We'll be using the most recent one off of this list.*/
        let iflEdits: any[] = []
        for (let val of this.$store.state.ambientStatus.inflights.postsList) {
          if (this.postFingerprint !== val.entity.provable.fingerprint) {
            continue
          }
          iflEdits.push(val)
        }
        return iflEdits
      },
      mostRecentInflightEdit(this: any) {
        let mostRecentTs = 0
        let mostRecent = undefined
        for (let val of this.inflightEdits) {
          if (val.status.requestedtimestamp >= mostRecentTs) {
            mostRecentTs = val.status.requestedtimestamp
            mostRecent = val
          }
        }
        return mostRecent
      },
      visiblePost(this: any) {
        if (typeof this.mostRecentInflightEdit !== 'undefined') {
          return this.mostRecentInflightEdit.entity
        }
        return this.post
      },
      visibleInflightStatus(this: any) {
        if (typeof this.mostRecentInflightEdit !== 'undefined') {
          return this.mostRecentInflightEdit.status
        }
        return this.inflightStatus
      },
      /*----------  Visibility checks  ----------*/
      /*
        These are useful because our visible / nonvisible logic is getting too complex to retain in the template itself.
      */
      actionsVisible(this: any) {
        if (this.localUserReadOnly) {
          return false
        }
        if (this.uncompiled) {
          return false
        }
        if (this.isInflightEntity) {
          return false
        }
        return true
      },
    },
    mounted(this: any) {
      // console.log("post self created")
      // console.log(this.post.selfcreated)
    },
    methods: {
      getOwnerName(owner: any): string {
        return globalMethods.GetOwnerName(owner)
      },
      /*----------  Collapse / open state  ----------*/
      toggleCollapse(this: any) {
        if (this.collapsed) {
          this.collapsed = false
          return
        }
        this.collapsed = true
      },
      /*----------  Reply actions  ----------*/
      toggleReplyPane(this: any) {
        console.log('toggle reply pane runs')
        if (this.replyPaneOpen) {
          this.replyPaneOpen = false
        } else {
          this.replyPaneOpen = true
          // this.editPaneOpen = false
        }
      },
      submitPost(this: any, fields: any) {
        let postBody = ""
        for (let val of fields) {
          if (val.id === 'postBody') {
            postBody = val.content
          }
        }
        let post = new mimobjs.Post
        // Set board, thread, parent, body fields
        post.setBoard(this.post.board)
        post.setThread(this.post.thread)
        post.setParent(this.post.fingerprint)
        post.setBody(postBody)
        let vm = this
        fe.SendPostContent('', post, function(resp: any) {
          console.log(resp.toObject())
          vm.toggleReplyPane()
        })
      },
      /*----------  Edit actions  ----------*/
      toggleEditPane(this: any) {
        if (this.editPaneOpen) {
          this.editPaneOpen = false
        } else {
          // this.replyPaneOpen = false
          this.postEditExistingSpec.fields[0].content = this.visiblePost.body
          this.editPaneOpen = true
        }
      },
      submitEditExistingPost(this: any, fields: any) {
        let postBody = ""
        for (let val of fields) {
          if (val.id === 'postBody') {
            postBody = val.content
          }
        }
        let post = new mimobjs.Post
        // Set board, thread, parent, body fields
        let pv = new mimobjs.Provable
        pv.setFingerprint(this.post.fingerprint)
        post.setProvable(pv)
        post.setBoard(this.post.board)
        post.setThread(this.post.thread)
        post.setParent(this.post.parent) // < heads up, this is different
        post.setBody(postBody)
        console.log(this.post.fingerprint)
        let vm = this
        fe.SendPostContent(this.post.fingerprint, post, function(resp: any) {
          console.log(resp.toObject())
          vm.toggleEditPane()
        })
      },
      /*----------  Report actions  ----------*/
      toggleReportPane(this: any) {
        if (this.reportPaneOpen) {
          this.reportPaneOpen = false
        } else {
          this.reportPaneOpen = true
        }
      },
      submitReport(this: any, fields: any) {
        let reportReason = ""
        for (let val of fields) {
          if (val.id === 'reportReason') {
            reportReason = val.content
          }
        }
        fe.ReportToMod(this.post.fingerprint, '', reportReason, function(resp: any) {
          console.log(resp.toObject())
        })
      },
      /*----------  Refresh func inflight info  ----------*/
      refresh(this: any) {
        console.log('post refresher is called')
        vuexStore.dispatch('refreshCurrentThreadAndPosts', {
          boardfp: this.$store.state.route.params.boardfp,
          threadfp: this.$store.state.route.params.threadfp
        })
      }
    }
  }
</script>

<style lang="scss">
  .post.inflight {
    .markdowned p:last-child {
      margin-bottom: 0;
    }
  }
</style>

<style lang="scss" scoped>
  @import "../scss/globals";
  .post {
    display: flex;
    flex-direction: column;
    padding-left: 20px;
    margin: 20px 20px;
    border-left: 3px solid rgba(255, 255, 255, 0.15);
    position: relative;

    &.inflight {}

    &>.post {
      margin-left: 0;
      margin-right: 0;
      padding-left: 15px;
    }
    .meta-container {
      display: flex;
      .post-datetime {
        margin-left: 10px;
        font-family: "SSP Regular Italic";
        color: $a-grey-600;
      }
    }
    .content {
      font-family: "SSP Regular";
      display: flex;
      font-size: 110%;

      .content-text-container {
        flex: 1;
      }

      .content-text {
        flex: 1;
        p:last-of-type {
          margin-bottom: 5px;
        }
      }
    }
    .ballot {
      padding-top: 0px;
    }
    &:hover {
      .ballot {
        visibility: visible;
      }
    }
  }

  .actions-rack {
    display: flex;
    .action {
      margin-right: 10px;
      font-family: "SSP Semibold";
      font-size: 14.4px;
      line-height: 21.6px;
    }
  }

  .report-container {
    padding-top: 35px;
  }

  .expand-container {
    position: absolute;
    top: -8px;
    left: -14px;
    background-color: $mid-base;
    border-radius: 10px;
    display: flex;
    cursor: pointer;
    height: 35px;
    width: 25px;
    &:hover {
      svg {
        fill: $a-grey-600;
      }
    }
    svg {
      fill: rgba(90, 90, 90, 1);
      width: 12px;
      height: 12px;
      margin: 0 auto;
      margin-top: 15px;
    }
  }
</style>