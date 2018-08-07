<template>
  <div class="settings-sublocation create-user">
    <a-markdown :content="headline"></a-markdown>
    <template v-if="localUserExists && !mintingStarted">
      <!-- ^ Minting started conditional because we don't want to show this the first time the user actually completes the process.  -->
      <a-markdown :content="userAlreadyExistsIntro"></a-markdown>
      <a-markdown :content="userAlreadyExistsContent"></a-markdown>
      <a class="button is-warning is-outlined" @click="goBack">GO BACK</a>
    </template>
    <template v-else>
      <template v-if="initialFormVisible">
        <a-markdown :content="intro"></a-markdown>
        <a-composer id="userComposer" :spec="createNewUserSpec"></a-composer>
      </template>
      <template v-if="inProgressVisible">
        <a-markdown :content="intro"></a-markdown>
        <a-markdown :content="mintingInProgressContent"></a-markdown>
        <a-spinner :hideText="true"></a-spinner>

      </template>
      <template v-if="completionVisible">
        <a-markdown :content="intro"></a-markdown>
        <a-markdown :content="completionContent"></a-markdown>
        <div class="video-container">
          <video class="words-of-wisdom" src="src/app/ext_dep/videos/wise_doggo.mp4" loop="true" muted="true" autoplay="true"></video>
        </div>
        <router-link to="/popular" class="button is-success is-outlined">
          GO TO POPULAR
        </router-link>
      </template>
    </template>

  </div>
</template>

<script lang="ts">
  // var globalMethods = require('../../../services/globals/methods')
  var mixins = require('../../../mixins/mixins')
  var fe = require('../../../services/feapiconsumer/feapiconsumer')
  var mimobjs = require('../../../../../../protos/mimapi/structprotos_pb.js')
  export default {
    name: 'newuser',
    mixins: [mixins.localUserMixin],
    data(this: any): any {
      return {
        headline: headline,
        intro: intro,
        content: content,
        userAlreadyExistsIntro: userAlreadyExistsIntro,
        userAlreadyExistsContent: userAlreadyExistsContent,
        mintingInProgressContent: mintingInProgressContent,
        completionContent: completionContent,
        mintingStarted: false,
        createNewUserSpec: {
          fields: [{
              id: "userName",
              visibleName: "Pick a name",
              description: `These names are <p class="em">not unique</p>, there can be multiple users with the same name. However, the blockavatars of two different users won't ever be the same. When in doubt, check name <i>and</i> the picture. <div id="postscript">(BTW, funders of the work on this app can get unique  names and flair in recognition of their support. Check the Membership tab to the left if you're interested in that.)</div>`,
              placeholder: "deanmoriarty",
              maxCharCount: 24,
              heightRows: 1,
              previewDisabled: true,
              content: "",
              optional: false,
            },
            {
              id: "userInfo",
              visibleName: "Info",
              description: "Optional, can be changed later. Markdown is available.",
              placeholder: "rebel without a cause / new york - san francisco",
              maxCharCount: 20480,
              heightRows: 5,
              previewDisabled: false,
              content: "",
              optional: true,
            }
          ],
          commitActionName: "CREATE",
          commitAction: this.submitNewUser,
          cancelActionName: "",
          cancelAction: function() {},
          fixToBottom: true,
          autofocus: true,
        }
      }
    },
    computed: {
      initialFormVisible(this: any) {
        if (this.mintingStarted) {
          return false
        }
        return true
      },
      inProgressVisible(this: any) {
        if (this.$store.state.localUserExists) {
          return false
        }
        if (!this.mintingStarted) {
          return false
        }
        return true
      },
      completionVisible(this: any) {
        if (this.$store.state.localUserExists) {
          return true
        }
        return false
        // if (globalMethods.IsUndefined(this.$store.state.localUser)) {
        //   return false
        // }
        // if (globalMethods.IsEmptyObject(this.$store.state.localUser)) {
        //   return false
        // }
        // return true
      },
    },
    methods: {
      goBack(this: any) {
        history.back()
      },
      submitNewUser(this: any, fields: any) {
        this.mintingStarted = true
        let userName = ""
        let userInfo = ""
        for (let val of fields) {
          if (val.id === 'userName') {
            userName = val.content
            continue
          }
          if (val.id === 'userInfo') {
            userInfo = val.content
            continue
          }
        }
        let user = new mimobjs.Key
        user.setName(userName)
        user.setInfo(userInfo)
        // let vm = this
        fe.SendUserContent('', user, function(resp: any) {
          console.log('user create request sent in.')
          console.log(resp.toObject())
        })
      }
    }
  }
  /*<br><br>(PS. Supporters of the work on this app can get unique  names and flair in recognition of their support. Check the Membership tab to the left if you're interested in that.)*/
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = '# Create new user'
  var intro =
    `**Hey there! 👋 &nbsp; Let's get you set up.**`
  var content =
    `Content text`
  var userAlreadyExistsIntro = `**There's already a user present on this app.**`
  var userAlreadyExistsContent = `You can enable the existing user by opening the menu at the top right, and choosing \`\`\`Exit read-only mode\`\`\` at the bottom.`
  var mintingInProgressContent = `
### Minting in progress...
Minting the proof-of-work for your user key. This can take from around 10 seconds to a minute depending on your computer.`
  var completionContent = `
### Successfully created.
  Your user is now ready. You can write and edit posts, upvote and downvote content, elect and impeach mods, and create & moderate communities.`
</script>

<style lang="scss">
  /* <<--global, not scoped */

  @import "../../../scss/globals";
  #userComposer {
    font-family: "SSP Bold";
    p.em {
      font-family: "SSP Black";
      display: inline;
    }
    #postscript {
      font-family: "SSP Regular Italic";
      padding-top: 10px;
      letter-spacing: 0.3px;
    }
    .visible-name {
      color: $a-grey-600;
    }
  }
</style>

<style lang="scss" scoped>
  @import "../../../scss/globals";
  @import"../../../scss/bulmastyles";

  .settings-sublocation {
    color: $a-grey-600;

    &.create-user {
      // font-size: 16px;
    }
    .markdowned {
      &:first-of-type {
        margin-bottom: 0;
      }
      margin-bottom: 40px;
    }
  }

  .button {
    font-family: "SSP Semibold"
  }

  .words-of-wisdom {
    width: 350px;
    margin-bottom: 30px;
    border-radius: 5px;
  }
</style>