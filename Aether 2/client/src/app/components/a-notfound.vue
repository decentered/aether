<template>
  <div class="notfound-container">
    <div class="not-found">
      <div class="headline">
        <icon name="unlink" scale="4"></icon>
      </div>
      <div class="subheader">
        Entity not found
      </div>
      <div class="description">This should not happen often.</div>
      <div class="description">(If it does, please reach out, it's worth a look.)</div>
      <div class="goback">
        <a class="button is-warning is-outlined" @click="goBack">GO BACK</a>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'a-notfound',
    data() {
      return {}
    },
    methods: {
      goBack(this: any) {
        history.back()
      }
    }
  }
</script>

<style lang="scss" scoped>
  @import"../scss/bulmastyles";
  @import "../scss/globals";
  .notfound-container {
    width: 100%;
    display: flex;
    padding-top: 100px;

    .not-found {
      margin: auto; // padding-bottom: 7%; // To compensate for the line-height of the top thing.
      font-family: "SCP Semibold";
      color: $a-grey-600;
      text-align: center;
      .headline {
        font-size: 700%;
        line-height: 100%;
        font-family: "SCP ExtraLight";
        padding-bottom: 15px;
      }
      .subheader {
        font-size: 125%;
        font-family: "SCP Bold";
      }
      .description {
        font-family: "SCP Regular";
      }
      .goback {
        padding-top: 40px;
        font-family: "SSP Semibold"
      }
    }
  }
</style>