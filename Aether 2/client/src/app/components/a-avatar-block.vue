<template>
  <router-link class="profile-header" tag="div" :to="link" :class="{'clickable': clickable}">
    <div class="profile-avatar">
      <div class="profile-img">
        <a-hashimage :hash="user.fingerprint" isUser="true" height="128px"></a-hashimage>
      </div>
    </div>
    <div class="profile-name">
      <div class="profile-name-text">
        {{getUserName()}}
      </div>
    </div>
    <div class="profile-fingerprint">
      <div class="profile-fingerprint-text">
        {{user.fingerprint}}
      </div>
    </div>
  </router-link>
</template>

<script lang="ts">
  var globalMethods = require('../services/globals/methods')
  export default {
    name: 'a-avatar-block',
    props: ['user', 'clickable'],
    data() {
      return {}
    },
    methods: {
      getUserName(this: any): string {
        // console.log(this.user)
        return globalMethods.GetUserName(this.user)
      }
    },
    computed: {
      // userName(this: any) {
      //   return this.user.NonCanonicalName
      // },
      link(this: any): string {
        if (this.clickable) {
          return '/user/' + this.user.fingerprint
        }
        return ''
      }
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/globals";
  .profile-header {
    width: 320px;
    display: flex;
    flex-direction: column;
    padding: 20px 15px 12px 15px;
    &.clickable {
      cursor: pointer;
    }
    .profile-avatar {
      display: flex;

      .profile-img {
        width: 128px;
        height: 128px;
        margin: auto;
      }
    }
    .profile-fingerprint {
      display: flex;

      .profile-fingerprint-text {
        font-family: "SCP Bold";
        word-wrap: break-word;
        width: 258px;
        font-size: 80%;
        margin: auto;
        margin-top: 20px;
        margin-bottom: 6px;
        color: $a-grey-600;
        background-color: rgba(255, 255, 255, 0.075);
        padding: 2px 6px;
      }
    }
    .profile-name {
      display: flex;
      margin-top: 20px;
      .profile-name-text {
        font-size: 125%;
        margin: auto;
      }
    }
  }
</style>