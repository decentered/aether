<template>
  <div class="actions">
    <div class="action">
      <img src="./Actions/assets/ic_keyboard_arrow_up_black_24px.svg" alt="" class="action-icon upvote">
    </div>
    <div class="action">
      <img src="./Actions/assets/ic_keyboard_arrow_down_black_24px.svg" alt="" class="action-icon downvote">
    </div>
    <div class="action">
      <img src="./Actions/assets/ic_mode_comment_black_24px.svg" alt="" class="action-icon comment">
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
  .actions
  {
    display:flex;
    margin-top:5px;
    .action
    {
      background-color:#ddd;
      border-radius:2px;
      margin-right:3px;
      height:24px;
      width:24px;
      cursor:pointer;
      &:hover
      {
        background-color:#bbb;
      }
      .action-icon
      {
        &.upvote
        {

        }
        &.downvote
        {

        }
        &.comment
        {
          width:16px;
          /*margin-right:1px;*/
          margin-left:4px;
          margin-bottom:2px;
        }
      }
    }
  }
</style>
