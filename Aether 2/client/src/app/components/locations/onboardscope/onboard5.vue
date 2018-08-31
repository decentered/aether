<template>
  <div class="onboard-sublocation">
    <div class="onboard-carrier">
      <transition name="fade" appear>
        <a-markdown class="onboard-markdown" :content="content"></a-markdown>
      </transition>
      <div class="continue-box">
        <router-link class="button is-outlined next-button" to="/onboard/6">
          ACCEPT
        </router-link>
        <div class="button is-outlined next-button" @click="quitApp">
          DECLINE AND QUIT
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  var ipc = require('../../../../../node_modules/electron-better-ipc')
  export default {
    name: 'onboard5',
    data() {
      return {
        content: textContent
      }
    },
    methods: {
      quitApp(this: any) {
        ipc.callMain('QuitApp')
      }
    }
  }
  var textContent = `
## Information wants to be free

If you post any text on Aether, you agree that it will be licensed under [*Creative Commons BY-SA*](https://creativecommons.org/licenses/by-sa/4.0/) license.

This is so that everyone can share content you create.

This only applies to text directly put in, though, not what links point to. If you make a video and post a link to it, it doesn't change the license of that video.
  `
</script>

<style lang="scss" scoped>
  @import"../../../scss/bulmastyles";
  @import"../../../scss/globals";
  .onboard-sublocation {
    display: flex;
    flex: 1;
    min-height: 100%;
    .onboard-carrier {
      margin: auto;
      display: flex;
      width: 30%;
      margin-top: 200px;
      flex-direction: column;
    }
  }

  .continue-box {
    display: flex;
    margin-top: 15px;
    .next-button {
      @extend %link-hover-ghost-extenders-disable;
      background-color: $a-transparent;
      color: $a-grey-800;
      &:hover {
        background-color: $a-grey-800;
        color: $mid-base;
      }
      animation-duration: 1.5s;
      animation-name: DELAY_VISIBLE;
      margin-right:10px;
    }
  }

  @keyframes DELAY_VISIBLE {
    0% {
      opacity: 0;
    }
    60% {
      opacity: 0;
    }
    100% {
      opacity: 1;
    }
  }

  .fade-enter-active,
  .fade-leave-active {
    transition-property: opacity;
    transition-duration: .25s;
  }

  .fade-enter-active {
    transition-delay: .25s;
  }

  .fade-enter,
  .fade-leave-active {
    opacity: 0
  }
</style>

<style lang="scss">
  .onboard-markdown {
    font-size: 120%;
    p {
      font-family: "SSP Regular"
    }
  }
</style>