<template>
  <div class="tickers">
    <div class="ticker vote-balance">
      <img class="ticker-icon up-arrow" src="./Tickers/assets/ic_keyboard_arrow_up_black_24px.svg" alt="">
      <div class="value">
        17
      </div>
    </div>
    <div class="ticker comments">
      <img class="ticker-icon comment" src="./Tickers/assets/ic_mode_comment_black_24px.svg" alt="">
      <div class="value">
        4
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    components: {},
    data () {
      return {}
    },
    methods: {}
  }
</script>

<style scoped lang="scss">
  .tickers
  {
    display:flex;
    .ticker
    {
      display:flex;
      margin-right:7px;
      &:last-of-type
      {
        margin-right:0;
      }
      img.ticker-icon
      {
        &.comment
        {
          width:16px;
          margin-right:4px;
          margin-left:2px;
        }
        &.up-arrow
        {
          width:28px;
          margin-right:-3px;
          margin-left:-5px;
        }
      }
    }
  }
</style>
