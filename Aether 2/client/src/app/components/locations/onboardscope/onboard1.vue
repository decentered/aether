<template>
  <div class="onboard-sublocation">
    <div class="onboard-carrier">
      <transition name="fade" appear>
        <a-markdown class="onboard-markdown" :content="content"></a-markdown>
      </transition>
      <div class="continue-box">
        <router-link class="button is-outlined next-button" to="/onboard/2">
          CONTINUE
        </router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'onboard1',
    data() {
      return {
        content: textContent,
      }
    },
  }
  var textContent = `
## Welcome in!

Aether is a network of independent, self-moderated communities. Some are friend groups, some are based on shared interests, and some are for cats.
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
    em {
      font-family: "SSP Regular Italic";
      strong {
        font-family: "SSP Bold Italic"
      }
    }
    h1,
    h2,
    h3,
    h4,
    h5,
    h6 {
      em {
        font-family: "SSP Bold Italic";
        strong {
          font-family: "SSP Black Italic"
        }
      }
    }
  }
</style>