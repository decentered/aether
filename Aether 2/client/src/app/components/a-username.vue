<template>
  <router-link class="username" :to="userLink" :class="{'disabled-link': isFingerprint}">
    {{ownerName}}
  </router-link>
</template>

<script lang="ts">
  var globalMethods = require('../services/globals/methods')
  export default {
    name: 'a-username',
    props: ['owner'],
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
    color: inherit;
    &.disabled-link {
      cursor: default;
      @extend %link-hover-ghost-extenders-disable;
      &:hover {
        background-color: unset;
      }
    }
  }
</style>