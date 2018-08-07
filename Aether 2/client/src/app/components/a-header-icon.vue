<template>
  <div class="header-icon-container">
    <icon :name="icon"></icon>
  </div>
</template>
<style lang="scss">
  @import "../scss/globals";
  .header-icon-container {
    -webkit-app-region: no-drag;
    height: $top-bar-height;
    width: $top-bar-height;
    display: flex;
    cursor: pointer;
    svg {
      margin: auto;
    }
    &:hover {
      background-color: rgba(255, 255, 255, 0.05);
    }
  }
</style>
<script lang="ts">
  var Vue = require('../../../node_modules/vue/dist/vue.js')
  export default Vue.extend({
    props: ['icon']
  })
</script>