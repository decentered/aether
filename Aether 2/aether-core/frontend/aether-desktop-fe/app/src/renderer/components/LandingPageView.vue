<template>
  <div>
    <img src="./LandingPageView/assets/logo.png" alt="electron-vue">
    <h1>O yea.</h1>
    <a href="#/tester">tester</a>
    <current-page></current-page>
    <versions></versions>
    <links></links>
    <div class="draggable-box">
      drag me!
    </div>
  </div>
</template>

<script>
  import CurrentPage from './LandingPageView/CurrentPage'
  import Links from './LandingPageView/Links'
  import Versions from './LandingPageView/Versions'
  export default {
    components: {
      CurrentPage,
      Links,
      Versions
    },
    name: 'landing-page'
  }
</script>

<style scoped>
.draggable-box
{
  width:300px;
  height:300px;
  background-color:aqua;
  -webkit-app-region: drag;
  app-region: drag;
  -webkit-user-select: none;
  user-select: none;
}
  img {
    margin-top: -25px;
    width: 450px;
  }
</style>
