<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content"></a-markdown>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'about',
    data() {
      return {
        headline: headline,
        intro: intro,
        content: content,
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = '# About'
  var intro =
    `**Aether v2.0.0**`
  var content =
    `
### What is Aether?

Aether is a client for the c0 subprotocol of Mim protocol, which establishes a protocol for efficient diffing and distribution of entity graphs over peer-to-peer networks.

### License

The source code for Aether is licensed under GNU AGPLv3 and can be found [here](https://github.com/nehbit/aether).



    `
</script>

<style lang="scss" scoped>
  @import "../../../scss/globals";
  .settings-sublocation {
    color: $a-grey-600;
    .markdowned {
      &:first-of-type {
        margin-bottom: 0;
      }
      margin-bottom: 40px;
    }
  }
</style>