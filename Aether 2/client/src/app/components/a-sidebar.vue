<template>
  <div class="sidebar-container">
    <div class="sidebar-group subs">
      <router-link class="special-sidebar-item" to="/">
        <div class="sidebar-item-icon"></div>
        <div class="sidebar-item-text">Home</div>
      </router-link>
      <router-link class="special-sidebar-item" to="/popular">
        <div class="sidebar-item-icon"></div>
        <div class="sidebar-item-text">Popular</div>
      </router-link>
      <router-link tag="div" class="sidebar-group-header" to="/globalscope/subbed">
        <div class="header-icon"></div>
        <div class="header-text">
          SUBS
        </div>
      </router-link>
      <router-link class="sidebar-item iterable" v-for="board in ambientBoards" :to="'/board/'+board.fingerprint" :class="{'updated':board.lastUpdate>board.lastseen}">
        <div class="sidebar-item-icon"></div>
        <div class="sidebar-item-text">{{board.name}}</div>
        <div class="sidebar-item-notifier">
          <div class="notifier-dot" v-show="board.lastupdate>board.lastseen"></div>
        </div>
      </router-link>
      <router-link class="sidebar-item iterable browse-boards" :to="'/globalscope'">
        <div class="sidebar-item-icon">
          <icon name="plus"></icon>
        </div>
        <div class="sidebar-item-text">Browse communities</div>
        <!-- <div class="sidebar-item-notifier">
          <div class="notifier-dot"></div>
        </div> -->
      </router-link>
    </div>
    <div class="sidebar-group status">
      <router-link tag="div" class="sidebar-group-header" to="/status">
        <div class="header-icon"></div>
        <div class="header-text">
          STATUS
        </div>
      </router-link>
    </div>

  </div>
</template>

<script lang="ts">
  // var feapiconsumer = require('../services/feapiconsumer/feapiconsumer')
  var Vuex = require('../../../node_modules/vuex').default
  export default {
    name: 'a-sidebar',
    data() {
      return {
        allBoards: {}
      }
    },
    computed: {
      ...Vuex.mapState(['ambientBoards']),
    },
    mounted(this: any) {
      let vm = this
      console.log(vm)
      console.log("side bar is created")
      // feapiconsumer.GetAllBoards(function(result: any) {
      //   console.log("callback gets called")
      //   vm.allBoards = result
      //   console.log(result)
      //   console.log(vm)
      // })
      // setTimeout(function() {
      //   console.log("all boards length")
      //   console.log(vm.allBoards.length)
      // }, 10000)
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";

  .sidebar-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    @include generateScrollbar($a-grey-100);
    .sidebar-group {
      // box-shadow: $line-separator-shadow-v2;
      // box-shadow: 0 3px 1px -2px rgba(0, 0, 0, 0.25);
      padding: 5px 0 10px 0;
      &.subs {
        flex: 1;
        height: 0; // https://stackoverflow.com/a/14964944 or: I love CSS
        overflow-y: scroll;
      }
      &.global-locations {}

      &.status {
        height: 150px;
        background-color: $dark-base*0.8;
      }

      .sidebar-group-header {
        font-size: 80%;
        padding: 5px 10px;
        cursor: pointer;
        letter-spacing: 1.5px;
        color: $a-grey-500;

        &:hover {
          color: $a-grey-800;
        }
      }

      .special-sidebar-item {
        @extend .sidebar-item;
        &.router-link-exact-active {
          @extend .selected;
        }
      }

      .sidebar-item {
        width: 95%;
        margin: 1px 2% 1px 3%;
        padding: 5px 10px;
        border-radius: 3px;
        cursor: pointer;
        font-family: "SSP Semibold";
        display: flex;
        color: $a-grey-400;

        .sidebar-item-icon {
          // width: 18px;
          // height: 18px;
          display: flex;
          svg {
            margin: auto;
            width: 16px;
            height: 14px;
            margin-right: 2px;
          }
        }
        &:hover {
          background-color: rgba(255, 255, 255, 0.05);
          color: $a-grey-800;
          @extend %link-hover-ghost-extenders-disable;
        }

        &.selected {
          color: $a-grey-800;
          background-color: rgba(255, 255, 255, 0.1); // font-family: "SSP Bold";
          &:hover {
            background-color: rgba(255, 255, 255, 0.15)
          }
        } // &.router-link-exact-active {
        //   @extend .selected;
        // }
        .sidebar-item-text {
          flex: 1;
          line-height: 120%;
        }

        .sidebar-item-notifier {
          display: flex;
          padding-left: 10px;
          width: 18px;
          .notifier-dot {
            margin: auto;
            width: 8px;
            height: 8px;
            border-radius: 4px;
            background-color: $a-orange
          }
        }

        &.updated {
          color: $a-grey-500; // font-family: "SSP Bold";
        }
      }
    }
  }

  .iterable {
    &.router-link-active {
      @extend .selected;
    }
  }

  .special-sidebar-item {
    @extend .sidebar-item;

    &.router-link-exact-active {
      @extend .selected;
    }
  }

  .browse-boards {
    &.router-link-active {
      @extend .selected;
    }
  }
</style>