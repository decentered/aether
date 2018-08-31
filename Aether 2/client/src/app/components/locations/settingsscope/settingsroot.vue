<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content1"></a-markdown>
    <a-markdown class="settings-bottom-markdown" :content="content2"></a-markdown>
    <div class="config-location">
      <div class="config-location-header">
        User profile and frontend preferences file
      </div>
      <div class="config-location-text">
        {{$store.state.ambientStatus.frontendambientstatus.frontendconfiglocation}}
      </div>
    </div>
    <div class="config-location">
      <div class="config-location-header">
        Backend preferences file
      </div>
      <div class="config-location-text">
        {{$store.state.ambientStatus.backendambientstatus.backendconfiglocation}}
      </div>
    </div>
    <a-markdown :content="content3"></a-markdown>
    <router-link to="/settings/advanced" class="advanced-button button is-warning is-outlined">
      ADVANCED
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
        content1: content1,
        content2: content2,
        content3: content3,
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = '# Preferences'
  var intro =
    `**You can modify the important main settings here.**

*More accurately, you will be able to, whenever the GUI for modifying these is ready.*
In the meanwhile, if you need to change them, the [Advanced](#/settings/advanced) section describes how you can do so.
  `
  var content1 =
    `
## Defaults

| Setting | Value |
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
  var content2 = `
## Your user profile and preferences

* **These are the locations your settings and profile are stored in.**
The \`\`\`frontend-config\`\`\` file also contains your user profile and public / private key pair.

* **You should back up a copy of your \`\`\`frontend-config\`\`\` in a safe location.**
This file contains your *private key*, and anybody who has this key can become you on the network.

* **You can move the \`\`\`frontend-config\`\`\` file between machines.**
This will carry over your subscribed communities, your user profile and other details.
  `
  var content3 = `
* You can get a list of the major configuration options in the files pointed above in the Advanced tab.
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

      &.settings-bottom-markdown {
        margin-bottom: 0;
      }
    }
    .advanced-button {
      font-family: "SSP Bold"
    }
    .config-location {
      padding: 16px 20px;
      color: $a-grey-800;
      background-color: $mid-base*0.5;
      margin-bottom: 15px;
      border-radius: 3px;
      .config-location-header {
        padding-bottom: 10px;
        font-family: "SSP Bold";
        font-size: 17.6px;
      }
      .config-location-text {
        font-family: "SCP Regular";
        word-break: break-all;
      }
    }
  }
</style>