<template>
  <div class="thread-entity" v-if="threadReadyToRender" :class="{'uncompiled': uncompiled}">
    <!-- Without above v-if it fails when currentThread===undefined -->
    <div class="signals-container" v-if="contentSignalsVisible">
      <div class="thread-actions">
      </div>
      <div class="thread-vote-count">
        <icon name="arrow-up"></icon> {{threadVoteCount}}
      </div>
      <div class="thread-comment-count">
        <icon name="comment-alt"></icon> {{currentThread.postscount}}
      </div>
    </div>
    <div class="image-container" v-show="imageLoadedSuccessfully">
      <div class="image-box" :style="'background-image: url('+ sanitisedLink +')'" @click="openLightbox">
      </div>
      <div class="lightbox" @click="closeLightbox" v-show="lightboxOpen">
        <img :src="sanitisedLink" alt="" @load="handleImageLoadSuccess">
      </div>
    </div>
    <div class="main-data-container">
      <div class="inflight-box" v-if="inflightBoxAtTopVisible">
        <a-inflight-info :status="visibleInflightStatus"></a-inflight-info>
      </div>
      <div class="thread-name">
        {{currentThread.name}}
      </div>
      <div class="thread-link">
        <a :href="currentThread.link" @click.stop>{{currentThread.link}}</a>
      </div>
      <template v-if="!editPaneOpen">
        <a-markdown class="thread-body" :content="visibleThread.body"></a-markdown>
        <div class="inflight-box" v-if="inflightBoxAtBottomVisible">
          <a-inflight-info :status="visibleInflightStatus"></a-inflight-info>
        </div>
      </template>
      <template v-else>
        <a-composer :spec="threadEditExistingSpec" v-if="editPaneOpen"></a-composer>
      </template>
      <div class="meta">
        <div class="thread-owner">
          <a-username :owner="currentThread.owner"></a-username>
        </div>
        <div class="thread-datetime">
          <a-timestamp :creation="threadCreation" :lastupdate="threadLastUpdate"></a-timestamp>
        </div>
        <div class="actions-rack" v-if="actionsVisible">
          <a class="action edit" v-show="currentThread.selfcreated" @click="toggleEditPane">
          Edit
        </a>
          <a class="action report" v-show="!currentThread.selfcreated" @click="toggleReportPane">
          Report
        </a>
        </div>
      </div>
      <div class="report-container" v-if="reportPaneOpen">
        <a-composer :spec="reportSpec"></a-composer>
      </div>
    </div>
    <a-ballot v-if="actionsVisible" :contentsignals="currentThread.compiledcontentsignals" :boardfp="currentThread.board" :threadfp="currentThread.fingerprint"></a-ballot>
  </div>
</template>

<script lang="ts">
  // var Vuex = require('../../../node_modules/vuex').default
  var globalMethods = require('../services/globals/methods')
  var mixins = require('../mixins/mixins')
  var fe = require('../services/feapiconsumer/feapiconsumer')
  var mimobjs = require('../../../../protos/mimapi/structprotos_pb.js')
  export default {
    name: 'a-thread-header-entity',
    mixins: [mixins.localUserMixin],
    props: ['inflightStatus', 'thread', 'uncompiled'],
    // ^ Unused as of now, since we do not allow inflight threads to be opened, but could be useful in the future.
    data(this: any): any {
      return {
        imageLoadedSuccessfully: false,
        hasUpvoted: false,
        hasDownvoted: false,
        editPaneOpen: false,
        reportPaneOpen: false,
        lightboxOpen: false,
        threadEditExistingSpec: {
          fields: [{
            id: "threadBody",
            visibleName: "",
            description: "",
            placeholder: "",
            maxCharCount: 20480,
            heightRows: 5,
            previewDisabled: false,
            content: '',
            optional: false,
          }],
          commitAction: this.submitEditExistingThread,
          commitActionName: "SAVE",
          cancelAction: this.toggleEditPane,
          cancelActionName: "CANCEL",
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
        },
      }
    },
    computed: {
      // ...Vuex.mapState(['currentThread']),
      currentThread(this: any) {
        return this.thread
      },
      threadCreation(this: any) {
        // These are necessary because in uncompiled entities, these are in thread.provable.creation, but in compiled ones it's in thread.creation.
        if (this.uncompiled) {
          if (typeof this.thread.provable === 'undefined') {
            return 0
          }
          return this.thread.provable.creation
        }
        if (typeof this.mostRecentInflightEdit !== 'undefined') {
          if (this.isVisibleInflightEntity) {
            if (typeof this.mostRecentInflightEdit.entity.provable === 'undefined') {
              return 0
            }
            return this.mostRecentInflightEdit.entity.provable.creation
          }
        }
        return this.thread.creation
      },
      threadLastUpdate(this: any) {
        if (this.uncompiled) {
          if (typeof this.thread.updateable === 'undefined') {
            return 0
          }
          return this.thread.updateable.lastupdate
        }
        if (typeof this.mostRecentInflightEdit !== 'undefined') {
          if (this.isVisibleInflightEntity) {
            if (typeof this.mostRecentInflightEdit.entity.updateable === 'undefined') {
              return 0
            }
            return this.mostRecentInflightEdit.entity.updateable.lastupdate
          }
        }
        return this.thread.lastupdate
      },
      threadFingerprint(this: any) {
        if (this.uncompiled) {
          return this.thread.provable.fingerprint
        }
        return this.thread.fingerprint
      },
      threadVoteCount(this: any) {
        if (typeof this.currentThread.compiledcontentsignals === 'undefined') {
          return 0
        }
        return this.currentThread.compiledcontentsignals.upvotes - this.currentThread.compiledcontentsignals.downvotes
      },
      threadReadyToRender(this: any) {
        if (this.uncompiled) {
          // If uncompiled and this code is running, the thread data is already there, and the compiled data won't ever be there since it's uncompiled, which means we're ready to go.
          return true
        }
        if (this.isVisibleInflightEntity) {
          return true
        }
        if (typeof this.currentThread.compiledcontentsignals === 'undefined') {
          return false
        }
        return true
      },
      sanitisedLink(this: any) {
        if (typeof this.currentThread.link === 'undefined') {
          return ""
        }
        if (this.currentThread.link.substring(0, 8) === 'https://' ||
          this.currentThread.link.substring(0, 7) === 'http://') {
          return this.currentThread.link
        }
        return 'http://' + this.currentThread.link
      },
      /*----------  Inflight computeds  ----------*/
      inflightEdits(this: any) {
        let iflEdits = []
        for (let val of this.$store.state.ambientStatus.inflights.threadsList) {
          if (this.threadFingerprint !== val.entity.provable.fingerprint) {
            continue
          }
          if (val.status.eventtype !== 'UPDATE') {
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
      visibleThread(this: any) {
        if (typeof this.mostRecentInflightEdit !== 'undefined') {
          return this.mostRecentInflightEdit.entity
        }
        return this.currentThread
      },
      visibleInflightStatus(this: any) {
        if (typeof this.mostRecentInflightEdit !== 'undefined') {
          return this.mostRecentInflightEdit.status
        }
        return this.inflightStatus
      },
      isVisibleInflightEntity(this: any) {
        if (typeof this.visibleInflightStatus !== 'undefined') {
          return true
        }
        return false
      },
      /*----------  Visibility calculations  ----------*/
      /*
        This matters because we're starting to have more states. One state that we have is inflight, and in this state, we don't have results of the compilation, nor do we have things such as a fingerprint. As a result, both signals and actions are disabled.

        The second state we have is the uncompiled state. This is how we show entities that we have pulled directly from the backend. A major use case for this is to show a user's profile, to show posts, threads, boards that the user has created, and so on.
      */
      contentSignalsVisible(this: any) {
        if (this.localUserReadOnly) {
          return false
        }
        if (this.uncompiled || this.isVisibleInflightEntity) {
          return false
        }
        return true
      },
      actionsVisible(this: any) {
        if (this.localUserReadOnly) {
          return false
        }
        // if inflight and compiled =  true
        // if inflight and uncompiled = false
        // if just uncompiled = false
        if (this.isVisibleInflightEntity) {
          if (!this.uncompiled) {
            return true
          }
          return false
        }
        if (this.uncompiled) {
          return false
        }
        return true
      },
      inflightBoxAtTopVisible(this: any) {
        // We use this in two different places. Box at top is suitable for lists, and box at bottom is suitable for full page views.
        if (this.uncompiled && this.isVisibleInflightEntity) {
          return true
        }
        return false
      },
      inflightBoxAtBottomVisible(this: any) {
        // We use this in two different places. Box at top is suitable for lists, and box at bottom is suitable for full page views.
        if (this.isVisibleInflightEntity && !this.uncompiled) {
          return true
        }
        return false
      }
    },
    methods: {
      /*----------  Lightbox open/close  ----------*/
      openLightbox(this: any) {
        this.lightboxOpen = true
      },
      closeLightbox(this: any) {
        this.lightboxOpen = false
      },
      getOwnerName(owner: any): string {
        return globalMethods.GetOwnerName(owner)
      },
      /*----------  Edit thread actions  ----------*/
      toggleEditPane(this: any) {
        if (this.editPaneOpen) {
          this.editPaneOpen = false
        } else {
          this.threadEditExistingSpec.fields[0].content = this.visibleThread.body
          this.editPaneOpen = true
        }
      },
      submitEditExistingThread(this: any, fields: any) {
        let threadBody = ""
        for (let val of fields) {
          if (val.id === 'threadBody') {
            threadBody = val.content
          }
        }
        let thread = new mimobjs.Thread
        // Set board, thread, parent, body fields
        thread.setBody(threadBody)
        thread.setBoard(this.currentThread.board)
        let vm = this
        vm.toggleEditPane()
        fe.SendThreadContent(this.currentThread.fingerprint, thread, function(resp: any) {
          console.log(resp.toObject())
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
        fe.ReportToMod(this.currentThread.fingerprint, '', reportReason, function(resp: any) {
          console.log(resp.toObject())
        })
      },
      handleImageLoadSuccess(this: any) {
        this.imageLoadedSuccessfully = true
      },
    },
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";

  .lightbox {
    display: flex;
    outline: none;
    position: fixed;
    z-index: 999;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    background-color: rgba(0, 0, 0, 0.5);

    img {
      margin: auto;
      max-height: 90%;
      max-width: 90%;
    }
  }

  .thread-entity {
    display: block;
    display: flex;
    padding: 15px 5px;
    margin: 0 20px;
    color: $a-grey-800;
    .ballot {
      padding-top: 75px;
    }
    &:hover {
      .ballot {
        visibility: visible;
      }
    }

    &.uncompiled {
      .main-data-container .thread-name {
        cursor: default;
      }
    }

    .signals-container {
      width: 64px;
      font-size: 110%;
      margin: auto;
      margin-top: 0;
      display: flex;
      flex-direction: column;
      color: $a-grey-600;

      .thread-vote-count {
        flex: 1;
        margin: auto;
        font-family: 'SSP Black';
        display: flex;
        svg {
          margin: auto;
          margin-right: 3px;
        }
      }

      .thread-comment-count {
        @extend .thread-vote-count;
        svg {
          margin-right: 5px;
          height: 13px;
          width: 13px;
        }
      }
    }

    .image-container {
      width: 80px;
      padding: 6px 8px;

      .image-box {
        height: 100%;
        max-height: 100px;
        overflow: hidden;
        border-radius: 2px;
        background-color: $a-cerulean;
        background-size: cover;
        cursor: pointer;

        .thread-image {
          object-fit: cover;
          height: inherit;
        }
      }
    }

    .main-data-container {
      flex: 1;
      padding-right: 15px;
      .thread-name {
        font-size: 120%;
        cursor: pointer;
      }

      .thread-link {
        font-size: 95%;
      }

      .thread-body {
        margin-top: 25px;
        font-family: 'SSP Regular';
        font-size: 110%;
      }
    }
    .meta {
      display: flex;
      .thread-owner {}
      .thread-datetime {
        margin-left: 10px;
        font-family: "SSP Regular Italic";
        color: $a-grey-600;
      }
    }
  }

  .actions-rack {
    display: flex;
    margin-left: 12px;
    .action {
      margin: auto;
      margin-bottom: 0;
      margin-right: 10px;
      font-family: "SSP Semibold";
      font-size: 90%;
    }
  }

  .report-container {
    padding-top: 35px;
  }
</style>