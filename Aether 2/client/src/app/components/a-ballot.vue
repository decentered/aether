<template>
  <div class="ballot" :class="{ 'voted': voteState !== 0, 'locked': locked, 'downhidden': downhidden }">
    <!-- todo >0 -->
    <div class="upvote" :class="{ 'voted': voteState === 1 }" @click.prevent="toggleUpvoteState">
      <icon name="arrow-up"></icon>
    </div>
    <div class="downvote" v-show="!downhidden" :class="{ 'voted': voteState === -1 }" @click.prevent="toggleDownvoteState">
      <icon name="arrow-down"></icon>
    </div>
  </div>
</template>

<script lang="ts">
  var globalMethods = require('../services/globals/methods')
  var fe = require('../services/feapiconsumer/feapiconsumer')
  var beObj = require('../../../../protos/mimapi/structprotos_pb.js')
  var clObj = require('../../../../protos/clapi/clapi_pb.js')
  export default {
    name: 'a-ballot',
    props: ['contentsignals', 'downhidden', 'boardfp', 'threadfp'],
    data() {
      return {
        // upvoteCast: false,
        // downvoteCast: false,
      }
    },
    computed: {
      locked(this: any) {
        // If the baseline already has a change saved
        if (this.contentsignals.selfatdlastupdate > 0) {
          return true
        }
        // If there is already one vote cast prior, and there is at least one inflight that has started the processing.
        if ((this.contentsignals.selfupvoted || this.contentsignals.selfdownvoted) && this.inflightATD.status.completionpercent !== 0) {
          return true
        }
        // If there are two in-progress inflight atds (this would be very rare, if at all possible - you'd have to flip the vote an instant after it enters minting, and somehow mint multi-threaded. Regardless, it's good to be defensive.)
        let count = 0
        for (let val of this.validInflightATDs) {
          if (val.status.completionpercent > 0) {
            count++
          }
          if (count >= 2) {
            return true
          }
        }
        return false
      },
      voteState(this: any): number {
        if (this.inflightATD.entity.type === 1) {
          return 1
        }
        if (this.inflightATD.entity.type === 2) {
          return -1
        }
        if (this.contentsignals.selfupvoted) {
          return 1
        }
        if (this.contentsignals.selfdownvoted) {
          return -1
        }
        return 0
      },
      inflightATD(this: any) {
        let ev = new beObj.Vote // empty vote
        let newestEntity = ev.toObject()
        let iv = new clObj.InflightVote
        let inflightVote = iv.toObject()
        let st = new clObj.InflightStatus
        let inflightStatus = st.toObject()
        inflightVote.status = inflightStatus
        inflightVote.entity = newestEntity
        // console.log(inflightVote)
        let newestReqTs = 0
        for (let val of this.validInflightATDs) {
          if (val.status.requestedtimestamp > newestReqTs) {
            newestReqTs = val.status.requestedtimestamp
            inflightVote.entity = val.entity
            inflightVote.status = val.status
          }
        }
        return inflightVote
      },
      validInflightATDs(this: any) {
        let validInflights = []
        for (let val of this.$store.state.ambientStatus.inflights.votesList) {
          if (this.contentsignals.targetfingerprint !== val.entity.target) {
            continue
          }
          if (val.entity.typeclass !== 1) {
            continue
          }
          if (val.status.completionpercent === -1) {
            // failed request
            continue
          }
          validInflights.push(val)
        }
        return validInflights
      },
      /*----------  fe upvote/downvote metadata  ----------*/
      // boardfp(this: any) {
      //   if (!globalMethods.IsUndefined(this.board)) {
      //     return this.board
      //   }
      //   return ""
      // },
      // threadfp(this: any) {
      //   if (!globalMethods.IsUndefined(this.thread)) {
      //     return this.thread
      //   }
      //   return ""
      // }
      voteVisibleFingerprint(this: any) {
        // This is the visible fingerprint of the vote entity.
        /*
          This can be sourced from two places. it might either be, a) from an actual base entity, or from a completed inflight that hasn't yet been fully refreshed into the page. In either of the cases, the underlying entity exists, it's just that they put their fingerprint in different places.
        */
        if (!globalMethods.IsUndefined(this.inflightATD.entity.provable)) {
          if (this.inflightATD.status.completionpercent !== -1) {
            // -1 = failed
            return this.inflightATD.entity.provable.fingerprint
          }
        }
        return this.contentsignals.selfatdfingerprint
      }
    },
    methods: {
      toggleUpvoteState(this: any) {
        if (this.locked) { return }
        fe.Upvote(this.contentsignals.targetfingerprint, this.voteVisibleFingerprint, this.boardfp, this.threadfp, function(resp: any) {
          console.log('we got a callback: ', resp)
        })
      },
      toggleDownvoteState(this: any) {
        if (this.locked) { return }
        fe.Downvote(this.contentsignals.targetfingerprint, this.voteVisibleFingerprint, this.boardfp, this.threadfp, function(resp: any) {
          console.log('we got a callback: ', resp)
        })
      },
      getOwnerName(owner: any): string {
        return globalMethods.GetOwnerName(owner)
      },
    },
    // mounted(this: any) {
    //   console.log('this.contentsignals')
    //   console.log(this.contentsignals)
    // }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .ballot {
    display: flex;
    visibility: hidden;
    cursor: default;
    width: 130px;
    margin-left: 20px;
    height: fit-content;
    padding: 15px 0px;

    &.locked {
      // border: 1px solid $a-grey-600;
      // border-radius: 5px;
      .upvote {
        cursor: default;
        &::before {
          box-shadow: 0 0 0 2px $a-transparent inset;
        }
        &.voted svg {
          fill: $mid-base;
        }
        svg {
          fill: $a-transparent;
        }
        &:hover {
          &::before {
            box-shadow: 0 0 0 2px $a-transparent inset;
          }
          &.voted svg {
            fill: $mid-base;
          }
          svg {
            fill: $a-transparent;
          }
        }
      }
    }

    &.downhidden {
      width: 75px;
      margin-left: 10px;
      height: unset;
    }

    .voted {
      visibility: visible;
    }


    .upvote {
      // visibility: hidden;
      width: 50px;
      height: 50px;
      display: flex;
      cursor: pointer;
      position: relative;
      margin: auto;
      margin-left: 0;

      &::before {
        // nonvoted normal state
        content: '';
        width: 54px;
        height: 54px;
        left: -2px;
        top: -2px;
        display: inline-block;
        border-radius: 27px;
        position: absolute;
        opacity: 1;
        box-shadow: 0 0 0 2px $a-grey-500 inset;
      }

      svg {
        padding: 10px;
        margin: auto;
        width: 100%;
        height: 100%;
        fill: $a-grey-500;
        position: relative; // z-index: 2;
      }

      &:hover {
        // non-voted hover state
        &::before {
          box-shadow: 0 0 0 2px $a-grey-800 inset; // background-color: $a-cerulean;
        }

        svg {
          fill: $a-grey-800;
        }
      }
    }
  }

  .downvote {
    @extend .upvote;
    margin-left: 5px;
    margin-right: 10px;
  }


  .upvote {
    &.voted::before {
      background-color: $a-cerulean-100;
    }
    svg {
      // fill: $a-grey-800;
    }
  }

  .downvote {
    &.voted::before {
      background-color: $a-red-100;
    }
    svg {
      // fill: $a-grey-800;
    }
  }

  .upvote.voted {
    &::before {
      box-shadow: 0 0 0 2px $a-transparent inset; // background-color: $a-cerulean-100;
    }
    svg {
      fill: $mid-base; // fill: $a-transparent;
    }
    &:hover {
      cursor: default;
      svg {
        fill: $mid-base;
      }
      &::before {
        box-shadow: none;
      }
    }
  }

  .downvote.voted {
    &::before {
      box-shadow: 0 0 0 2px $a-transparent inset; // background-color: $a-red-100;
    }
    svg {
      // fill: $a-transparent;
      fill: $mid-base;
    }
    &:hover {
      cursor: default;
      svg {
        fill: $mid-base;
      }
      &::before {
        box-shadow: none;
      }
    }
  }
</style>