<template>
  <div class="aether-navigation-trail">
    <div class="trail-separator root">
      <img src="./AetherNavigationTrail/assets/Aether-Logo-bw-small.svg" alt="">
    </div>
    <div class="trail-separator">
      <img src="./AetherNavigationTrail/assets/ic_keyboard_arrow_right_white_24px.svg" alt="">
    </div>
    <div class="trail-item">Front page</div>
    <div class="trail-separator">
      <img src="./AetherNavigationTrail/assets/ic_keyboard_arrow_right_white_24px.svg" alt="">
    </div>
    <div class="trail-item">Tester</div>
    <div class="trail-separator">
      <img src="./AetherNavigationTrail/assets/ic_keyboard_arrow_right_white_24px.svg" alt="">
    </div>
    <div class="trail-item">Another tester</div>
  </div>
</template>

<script>
</script>

<style scoped lang="scss">
.aether-navigation-trail
{
  display:flex;
  padding:12px 20px;
  flex:1;
  .trail-separator
  {
    position:relative;
    bottom:1px;
    padding-top:3px;
    &.root
    {
      padding-right:6px;
    }
  }
  .trail-item, .trail-item a
  {
    color:white;
    cursor:pointer;
    padding: 3px 5px;
    height:20px;
    border-radius:2px;
    &:hover
    {
      background-color:#666;
    }
  }
}
</style>
