<template>
  <div class="onboard-sublocation">
    <div class="onboard-carrier">
      <transition name="fade" appear>
        <a-markdown class="onboard-markdown" :content="content"></a-markdown>
      </transition>
      <div class="continue-box">
        <router-link class="button is-outlined next-button" to="/onboard/5">
          CONTINUE
        </router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'onboard4',
    data() {
      return {
        content: textContent
      }
    }
  }
  var textContent = `
## ... and it is *public*.

Aether is a public discussion forum. This means everything you put into it will be public.

Eventually, you'll be able to create private communities. That feature is not ready yet, though, we're still working on&nbsp;it.
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