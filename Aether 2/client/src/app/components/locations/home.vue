<template>
  <div class="location">
    <div class="home">
      This is home.
    </div>
  </div>
</template>

<script lang="ts">
  var Vuex = require('../../../../node_modules/vuex').default
  export default {
    name: 'home',
    ata() {
      return {
        breadcrumbState: []
      }
    },
    methods: {
      ...Vuex.mapActions(['setBreadcrumbs']),
    },
    mounted() {
      console.log('home mounted routine gets called')
      this['setBreadcrumbs']([])
    },
    updated() {
      console.log('home updated routine gets called')
      this['setBreadcrumbs']([])
    }
  }
</script>

<style lang="scss" scoped>

</style>