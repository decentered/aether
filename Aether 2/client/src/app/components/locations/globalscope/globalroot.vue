<template>
  <div class="global-sublocation">
    <div class="global-root">
      <template v-if="!loadingComplete">
        <div class="spinner-container">
          <a-spinner></a-spinner>
        </div>
      </template>
      <template v-else>
        <a-board-entity v-for="board in whitelistedBoards" :board="board"></a-board-entity>
        <div class="load-more-carrier" v-show="whitelistedBoards.length >= this.whitelistedLastLoadedItem">
          <a class="button is-warning is-outlined load-more-button" @click="loadMoreWhitelisted()">
        LOAD MORE
              </a>
        </div>
        <div class="non-whitelist-info">
          <p>
            You've reached to the end of the whitelisted boards list.
          </p>
          <p>
            Whitelisted boards are the suggested boards that are chosen to be interesting and work-safe, and they contain some of the higher-quality discussion and content the network has available.
          </p>
          <p>
            If you prefer, you can view the list of non-whitelisted boards. A board being non-whitelisted is not necessarily a negative signal. It might mean that the board was just formed, or that it is interesting and high quality but not work-safe, which is perfectly fine, or any combination of many other reasons.
          </p>
          <p>
            However, it also means that the quality will be much more varied. Proceed at your own risk. Here be dragons.
          </p>
          <a @click="toggleNonWhitelisted()" v-show="true">
          <span v-show="!nonWhitelistedVisible">I'm OK with that, show me</span>
          <span v-show="nonWhitelistedVisible">Hide</span>
        </a>
        </div>
        <template v-if="nonWhitelistedVisible">
          <a-board-entity v-for="board in nonWhitelistedBoards" :board="board"></a-board-entity>
          <div class="load-more-carrier" v-show="nonWhitelistedBoards.length >= this.nonWhitelistedLastLoadedItem">
            <a class="button is-warning is-outlined load-more-button" @click="loadMoreNonWhitelisted()">
          LOAD MORE
                </a>
          </div>
        </template>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
  var Vuex = require('../../../../../node_modules/vuex').default
  export default {
    name: 'globalroot',
    data() {
      return {
        nonWhitelistedVisible: false,
        nonWhitelistedFirstLoadedItem: 0,
        nonWhitelistedLastLoadedItem: 25,
        whitelistedFirstLoadedItem: 0,
        whitelistedLastLoadedItem: 25,
        batchSize: 25,
      }
    },
    computed: {
      ...Vuex.mapState(['allBoards', 'allBoardsLoadComplete']),
      loadingComplete(this: any) {
        return this.allBoardsLoadComplete
      },
      whitelistedBoards(this: any) {
        let whitelisted = []
        let vm = this
        for (var i = 0; i < this['allBoards'].length; i++) {
          (function(i) {
            if (vm.allBoards[i].whitelisted) {
              whitelisted.push(vm.allBoards[i])
            }
          })(i)
        }
        return whitelisted.slice(this.whitelistedFirstLoadedItem, this.whitelistedLastLoadedItem)
      },
      nonWhitelistedBoards(this: any) {
        let nonWhitelisted = []
        let vm = this
        for (var i = 0; i < this['allBoards'].length; i++) {
          (function(i) {
            if (!vm.allBoards[i].whitelisted) {
              nonWhitelisted.push(vm.allBoards[i])
            }
          })(i)
        }
        return nonWhitelisted.slice(this.nonWhitelistedFirstLoadedItem, this.nonWhitelistedLastLoadedItem)
      }
    },
    methods: {
      toggleNonWhitelisted(this: any) {
        this.nonWhitelistedVisible = !this.nonWhitelistedVisible
      },
      loadMoreNonWhitelisted(this: any) {
        this.nonWhitelistedLastLoadedItem = this.nonWhitelistedLastLoadedItem + this.batchSize
      },
      loadMoreWhitelisted(this: any) {
        this.whitelistedLastLoadedItem = this.whitelistedLastLoadedItem + this.batchSize
      }
    }
  }
</script>

<style lang="scss" scoped>
  @import "../../../scss/bulmastyles";
  @import "../../../scss/globals";
  .non-whitelist-info {
    font-family: "SCP Regular";
    margin: 20px;
    padding: 20px;
    background-color: rgba(255, 255, 255, 0.05);
    a {
      font-family: 'SCP Bold';
      &:hover {
        color: $a-grey-800;
      }
    }
  }

  .button {
    font-family: 'SSP Semibold'
  }

  .load-more-carrier {
    display: flex;
    padding: 20px 0 0 0;
    .load-more-button {
      margin: auto;
    }
  }

  .spinner-container {
    display: flex;
    .spinner {
      margin: auto;
      padding-top: 50px;
    }
  }
</style>