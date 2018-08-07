<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content"></a-markdown>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'adminsquickstart',
    data() {
      return {
        headline: headline,
        intro: intro,
        content: content,
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = "# Admin's Quickstart"
  var intro =
    `**This is a quick introduction to founding your own community on Aether.**`
  var content =
    `
Thanks for being here! Aether is better because of you. The more communities with diverse views inhabiting this place, the more interesting, insightful, interesting it ends up being. We have space for cats, corgis, humans, motorcycles, \`\`\`python\`\`\`, pythons, cryptocurrencies, **Ć̢̘̳̺͉̯̘̣̠̂̓̔̓͗ͭ͊̿̿́͜ţ̛̫̗̱̗̖̬̱̜͖̬̐͊̌̃̀̅̒̋͛͝͡ͅh̷̤̦̝̟̜̓ͦͥ̿ͬ̐̔͘͜ụ̺̻͍̱̟̦͇̗̦̳̤͑̌͐̆̔͊ͧ̔̈͊̀͞l̸̴͔̣͚͇̯̭ͦ̋̈͌͊͂̇͂̚͟͝h̤̰͕͓̪͓͗ͭ̆̾ͮ̄ͪ͒́̆̉͊̍͜ͅͅȗ̸̮̦̦̣̠̝̫̹̟͈͇̯͚̏ͯ͂̎̚͞ͅ**, and anything in between.

## The prime directive

bork bork bork
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