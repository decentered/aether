<template>
  <div class="header-component">
    <div class="draggable-zone"></div>
    <div class="header-main-zone">
      <div class="history-navigation">
        <div class="header-button" @click="browserBack()">
          <img src="./AetherHeader/assets/ic_keyboard_arrow_left_white_24px.svg" alt="">
        </div>
        <div class="header-button"@click="browserForward()" >
          <img src="./AetherHeader/assets/ic_keyboard_arrow_right_white_24px.svg" alt="">
        </div>
      </div>
      <AetherNavigationTrail></AetherNavigationTrail>
      <div class="core-components">
        <div class="header-button">
          <img src="./AetherHeader/assets/placeholder.svg" alt="">
        </div>
        <div class="header-button">
          <img src="./AetherHeader/assets/placeholder.svg" alt="">
        </div>
        <div class="header-button">
          <img src="./AetherHeader/assets/placeholder.svg" alt="">
        </div>
        <div class="header-button">
          <img src="./AetherHeader/assets/placeholder.svg" alt="">
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import AetherNavigationTrail from './AetherNavigationTrail'
  export default {
    components: {
      AetherNavigationTrail
    },
    data () {
      return {
        // platform: os.platform(),
        // versions: process.versions
      }
    },
    methods: {
      browserBack: function () {
        window.history.back()
      },
      browserForward: function () {
        window.history.forward()
      }
    }
  }
</script>

<style scoped lang="scss">
  $header-background-color: #222;
  .header-component
  {
    height:70px;
    display:flex;
    flex-direction:column;
    .draggable-zone
    {
      height:22px;
      -webkit-app-region: drag;
      app-region: drag;
      -webkit-user-select: none;
      user-select: none;
      background-color:#111;
    }
    .header-main-zone
    {
      flex:1;
      /*background-color:#333;*/
      background-color: $header-background-color;
      display:flex;
    }
    .history-navigation
    {
      display:flex;
    }
    .header-button
    {
      display:flex;
      width:48px;
      cursor:pointer;
      padding:4px;
      padding-left:3px;
      img
      {
        flex:1;
      }
      &:hover
      {
        background-color:$header-background-color*1.5;
      }
    }
    .core-components
    {
      display:flex;
    }
  }
</style>
