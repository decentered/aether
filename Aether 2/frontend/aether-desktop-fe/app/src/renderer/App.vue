<template>
  <div id="#app" class="app-container">
    <router-view></router-view>
  </div>
</template>

<script>
  import store from 'renderer/vuex/store'
  export default {
    store
  }
</script>

<style lang="scss">
  /*@import url(https://fonts.googleapis.com/css?family=Lato:300);*/
  *
  {
    margin: 0;
    padding: 0;
  }
  html,body
  {
    height: 100%;
  }
  body
  {
    display:flex;
    margin:0;
    padding:0;
    font-family:"Work Sans";

    .app-container
    {
      flex:1;
    }
    #app
    {
      flex:1;
    }
    #vertical-flex-container
    {
      display:flex;
      flex-direction:column;
      height:100%;
    }
    #horizontal-flex-container
    {
      flex:1;
      display:flex;
      flex-direction:row;
    }
  }
</style>
