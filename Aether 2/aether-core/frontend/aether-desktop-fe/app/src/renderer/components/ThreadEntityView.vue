<template>
  <div id="vertical-flex-container">
    <AetherHeader></AetherHeader>
    <div id="horizontal-flex-container">
      <AetherSidebar></AetherSidebar>
      <div id="main-content">
        <ThreadHeader></ThreadHeader>
        <Post depth=1></Post>
        <Post depth=1></Post>
        <Post depth=2></Post>
        <Post depth=1></Post>
        <Post depth=2></Post>
        <Post depth=3></Post>
        <Post depth=3></Post>
      </div>
    </div>
  </div>
</template>

<script>
  import AetherHeader from './Shared/AetherHeader'
  import AetherSidebar from './Shared/AetherSidebar'
  import ThreadHeader from './ThreadEntityView/ThreadHeader'
  import Post from './Shared/Post'
  export default {
    components: {
      AetherHeader,
      AetherSidebar,
      ThreadHeader,
      Post
    },
    name: 'thread-entity-view'
  }
</script>

<style scoped lang="scss">
  #main-content
  {
    flex:1;
    background-color:#fff;
    padding: 25px;
    overflow:scroll;
  }
</style>
