<template>
  <div class="board-sublocation">
    <div class="board-root">
      <template v-if="$store.state.route.name === 'Board>ThreadsNewList'">
        <a-thread-entity v-for="thr in inflightNewThreads.slice().reverse()" :thread="thr.entity" :inflightStatus="thr.status" :key="thr.Fingerprint"></a-thread-entity>
      </template>
      <a-thread-entity v-for="thr in threadsList" :thread="thr" :key="thr.Fingerprint"></a-thread-entity>
      <a-no-content no-content-text="There are no threads yet. You should write something." v-if="hasNoContent"></a-no-content>
    </div>
  </div>
</template>

<script lang="ts">
  var Vuex = require('../../../../../node_modules/vuex').default
  export default {
    name: 'boardroot',
    data() {
      return {
        currentBoardsThreadsNew: []
      }
    },
    computed: {
      ...Vuex.mapState(['currentBoardsThreads']),
      threadsList(this: any) {
        // return this.currentBoardsThreads
        if (this.$store.state.route.name === 'Board' || this.$store.state.route.name === undefined) {
          return this.currentBoardsThreads
        }
        if (this.$store.state.route.name === 'Board>ThreadsNewList') {
          return this.currentBoardsThreadsNew
        }
      },
      inflightNewThreads(this: any) {
        let inflightNewThreads = []
        for (let val of this.$store.state.ambientStatus.inflights.threadsList) {
          if (val.status.eventtype !== 'CREATE') {
            continue
          }
          if (this.$store.state.currentBoard.fingerprint !== val.entity.board) {
            continue
          }
          inflightNewThreads.push(val)
        }
        console.log('inflight new threads')
        console.log(inflightNewThreads)
        return inflightNewThreads
      },
      hasNoContent(this: any) {
        if (this.$store.state.route.name === 'Board' || this.$store.state.route.name === undefined) {
          if (this.threadsList.length === 0) {
            return true
          }
        }
        if (this.$store.state.route.name === 'Board>ThreadsNewList') {
          if (this.inflightNewThreads.length === 0) {
            return true
          }
        }
        return false
      }
    },
    methods: {},
    mounted(this: any) {
      // console.log('this is currentboardsthreads')
      // console.log(this.currentBoardsThreads)
    },
    updated() {}
  }
</script>

<style lang="scss" scoped>

</style>