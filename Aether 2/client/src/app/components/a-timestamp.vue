<template>
  <div class="timestamp">
    {{computeTimestamp()}}
  </div>
</template>

<script lang="ts">
  var globalMethods = require('../services/globals/methods')
  export default {
    name: 'a-timestamp',
    props: ['creation', 'lastupdate'],
    data() {
      return {}
    },
    computed: {
      nowTs: function(): number {
        return Math.floor(Date.now() / 1000)
      }
    },
    methods: {
      computeTimestamp(this: any): string {
        let ts = this.creation > this.lastupdate ? this.creation : this.lastupdate
        return globalMethods.TimeSince(ts)
      },
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .timestamp {
    font-family: "SSP Regular Italic";
    color: $a-grey-600;
  }
</style>