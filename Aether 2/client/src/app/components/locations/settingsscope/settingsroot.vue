<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content"></a-markdown>
    <router-link to="/settings/advanced" class="button is-warning is-outlined">
      ADVANCED SETTINGS
    </router-link>

  </div>
</template>

<script lang="ts">
  export default {
    name: 'settingsroot',
    data() {
      return {
        headline: headline,
        intro: intro,
        content: content,
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = '# Settings'
  var intro =
    `**You can modify the more important main settings here.**

Or more accurately, you will be able to, whenever the GUI for modifying these is ready. In the meanwhile, if you need to change them, the [Advanced](#/settings/advanced) section describes how you can do so.
  `
  var content =
    `
## Default settings

These are the default settings of the app. You can find a description of these knobs and how to change them at the [Advanced](#/settings/advanced) section.

| Knob | Value |
|---|---|
| Maximum disk space for database | 10 Gb
| Local memory | 6 months |
| Neighbourhood size | 10
| Tick duration | 60 seconds
| Reverse open | Enabled
| Maximum address table size | 1000
| Maximum simultaneous inbound connections | 5
| Maximum simultaneous outbound connections | 1
  `
</script>

<style lang="scss" scoped>
  @import"../../../scss/bulmastyles";
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