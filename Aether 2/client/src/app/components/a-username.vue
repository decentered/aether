<template>
  <div class="user-name-container">
    <template v-if="isCanonical">
      <router-link class="username canonical" :to="userLink" :class="{'disabled-link': isFingerprint}" hasUsernameTooltip title="<div class='username-internal-container'><div class='usernametooltip-header'>This is a registered name. There can be no other registered names that are the same.</div> <div class='usernametooltip-body'>You can get registered names by funding the development of the <em>Aether Open Source Project.</em> You can <a href='http://bit.ly/supportaether'><b>fund it here.</b></a> <br><br>You can get more information in the <a href='#/membership'><b>Membership</b></a> page.
        <br><br><p> If this is you: Thank you. ❤️</p> </div></div> ">
        {{ownerName}}
        <icon class="canonical-icon" name="check-circle"></icon>
      </router-link>
    </template>
    <template v-else>
      <router-link class="username" :to="userLink" :class="{'disabled-link': isFingerprint}">
        {{ownerName}}
      </router-link>
    </template>
    <icon v-show="isop" class="original-poster-icon" hasInfomark title="<div class='infomark-body'>This user is the <b><em>original poster</em></b> (OP) of this thread.</div>" name="pencil-alt"></icon>
  </div>
</template>

<script lang="ts">
  var Tooltips = require('../services/tooltips/tooltips')
  var globalMethods = require('../services/globals/methods')
  export default {
    name: 'a-username',
    props: ['owner', 'isop'],
    data() {
      return {}
    },
    computed: {
      userLink(this: any) {
        if (this.isFingerprint) {
          return ""
        }
        return '/user/' + this.owner.fingerprint
      },
      isFingerprint(this: any) {
        if (typeof this.owner === 'string') {
          // This is a fingerprint (i.e. this entity is uncompiled, ergo, it's in the user view of that current user.)
          return true
        }
        return false
      },
      ownerName(this: any) {
        return globalMethods.GetUserName(this.owner)
      },
      isCanonical(this: any) {
        if (this.isFingerprint) {
          return false
        }
        if (this.owner.compiledusersignals.canonicalname.length > 0) {
          return true
        }
        return false
      }
    },
    mounted() {
      Tooltips.MountUsernameTooltip()
      Tooltips.MountInfomark()
    },
    updated() {
      Tooltips.MountUsernameTooltip()
      Tooltips.MountInfomark()
    },
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .username {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 20ch;
    color: inherit; // font-family: "SSP Semibold";
    border-radius: 5px; // color: $a-grey-600;
    &.disabled-link {
      cursor: default;
      @extend %link-hover-ghost-extenders-disable;
      &:hover {
        background-color: unset;
      }
    }
    &.canonical {
      // font-family: "SSP Bold";
      // background-color: $a-cerulean-20; // color: $a-grey-800; // font-family: "SSP Bold";
      background-image: linear-gradient(to right top, #e64440, #ec5e3d, #f0743d, #f28a40, #f49e47);
      padding: 2px 3px;
      padding-left: 5px;
      color: white;
      text-shadow: 0 0px 1px rgba(0, 0, 0, 0.4);
      @extend %link-hover-ghost-extenders-disable;
    }
    .canonical-icon {
      width: 12px;
      height: 12px;
      margin-left: 2px;
      margin-right: 2px;
    }
  }

  .original-poster-icon {
    width: 12px;
    height: 12px;
    margin-left: 7px;
  }
</style>