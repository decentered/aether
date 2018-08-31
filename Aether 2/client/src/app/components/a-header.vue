<template>
  <div class="header-container" :class="{'sidebar-closed': !sidebarOpen}">
    <div class="history-container">
      <a-header-icon icon="chevron-left" :class="{'disabled': !hasPrevious}" @click.native="goBackward"></a-header-icon>
      <a-header-icon icon="chevron-right" :class="{'disabled': !hasForward}" @click.native="goForward"></a-header-icon>
    </div>
    <div class="breadcrumbs-container">
      <a-breadcrumbs></a-breadcrumbs>
    </div>
    <a-notifications-icon v-if="localUserExists"></a-notifications-icon>
    <div class="profile-container" @click="toggleUserMenu" v-click-outside="onClickOutside">
      <div class="dropdown-container">
        <div class="dropdown is-right" :class="{ 'is-active': userMenuOpen }">
          <div class="dropdown-trigger">
            <template v-if="!$store.state.localUserArrived">
              <div class="user-name" :class="{'read-only': this.localUserReadOnly}">
                Refreshing...
              </div>
            </template>
            <template v-if="$store.state.localUserArrived">
              <div class="info-marker-container" v-if="localUserReadOnly">
                <a-info-marker header="You are in read only mode." text="<p>You haven't created an user. You can read, but you won't be able to post or vote until you create one.</p>"></a-info-marker>
              </div>
              <div class="user-name" :class="{'read-only': this.localUserReadOnly}">
                {{localUserName}}
              </div>
              <div class="mod-puck-container" v-show="isMod">
                <div class="mod-puck">
                  mod
                </div>
              </div>
            </template>
            <div class="profile-caret-icon">
              <icon name="chevron-down"></icon>
            </div>
          </div>
          <div class="dropdown-menu" id="dropdown-menu" role="menu">
            <div class="dropdown-content">
              <template v-if="!localUserExists">
                <!-- Local user doesn't exist, not created yet -->
                <!--    <router-link to="/settings/newuser" class="dropdown-item">
                  Create user
                </router-link> -->
                <div class="button-container">
                  <router-link class="button is-link is-outlined is-small join-button" to="/newuser">
                    JOIN AETHER
                  </router-link>
                </div>
                <hr class="dropdown-divider">
              </template>
              <template v-if="!localUserReadOnly">
                <a-avatar-block nofingerprint="true" :user="$store.state.localUser" :clickable="true"></a-avatar-block>
                <hr class="dropdown-divider">
                <router-link :to="'/user/'+$store.state.localUser.fingerprint" class="dropdown-item">
                  Profile
                </router-link>
              </template>
              <router-link to="/intro" class="dropdown-item">
                Beginner's guide
              </router-link>
              <router-link to="/settings" class="dropdown-item">
                Preferences
              </router-link>
              <hr class="dropdown-divider">
              <router-link to="/about" class="dropdown-item">
                About
              </router-link>
              <router-link to="/membership" class="dropdown-item">
                Membership
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
  var Vue = require("../../../node_modules/vue/dist/vue.js");
  var Vuex = require('../../../node_modules/vuex').default
  var mixins = require('../mixins/mixins')
  var globalMethods = require('../services/globals/methods')
  export default Vue.extend({
    name: 'a-header',
    mixins: [mixins.localUserMixin],
    data() {
      return {
        userMenuOpen: false,
      };
    },
    computed: {
      ...Vuex.mapState(['sidebarOpen']),
      localUserName(this: any) {
        if (this.localUserReadOnly) {
          return "Menu"
        }
        return globalMethods.GetUserName(this.$store.state.localUser)
      },
      isMod(this: any) {
        if (this.$store.state.modModeEnabledArrived && this.$store.state.modModeEnabled) {
          return true
        }
        return false
      },
      hasPrevious(this: any) {
        return this.$store.state.historyMaxCaret > 1
      },
      hasForward(this: any) {
        return this.$store.state.historyCurrentCaret < this.$store.state.historyMaxCaret
      },
    },
    methods: {
      ...Vuex.mapActions(['setSidebarState']),
      toggleUserMenu(): void {
        ( < any > this)['userMenuOpen'] ?
        ( < any > this)['userMenuOpen'] = false:
          ( < any > this)['userMenuOpen'] = true
      },
      onClickOutside() {
        ( < any > this)['userMenuOpen'] = false
      },
      toggleSidebarState() {
        this.$store.state.sidebarOpen === true ?
          this.setSidebarState(false) : this.setSidebarState(true)
      },
      signIn(this: any) {},
      signOut(this: any) {},
      goForward(this: any) {
        // if (!this.hasForward) {
        //   return
        // }
        this.$store.dispatch('registerNextActionIsHistoryMoveForward')
        this.$router.go(+1)
        // this.$router.push({ path: this.$routerHistory.next().path })
      },
      goBackward(this: any) {
        // if (!this.hasPrevious) {
        //   return
        // }
        this.$store.dispatch('registerNextActionIsHistoryMoveBack')
        this.$router.go(-1)
        // this.$router.push({ path: this.$routerHistory.previous().path })
      }
    }
  });
</script>
<style lang="scss" scoped>
  @import "../scss/bulmastyles";
  @import "../scss/globals";
  .header-container {
    -webkit-app-region: drag;
    /* ^ You must mark all interactive objects within this as NON-draggable to make this work, otherwise they will just be effectively unclickable. */
    height: $top-bar-height;
    background-color: $mid-base * 0.9;
    /*border-bottom: 1px solid #111;*/
    box-shadow: $line-separator-shadow, $line-separator-shadow-castleft-light;
    position: relative; // z-index: 3;
    display: flex;
    flex: 1;
    border-radius: 10px 0 0 0;
    min-width: 0;

    &.sidebar-closed {
      border-radius: 0;
    }
    .history-container {
      display: flex;
      .disabled {
        opacity: 0.15;
        cursor: default;
        &:hover {
          background-color: $a-transparent;
        }
      }
    }

    .breadcrumbs-container {
      flex: 1;
      min-width: 0;
    }

    .profile-container {
      display: flex; // padding: 0 15px;
      padding-right: 15px;
      padding-left: 7px;

      &:hover {
        background-color: rgba(255, 255, 255, 0.05);
      }

      .dropdown-container {
        display: flex;

        .profile-header {
          width: 225px;
        }

        .dropdown-trigger {
          display: flex;
          cursor: pointer;

          .user-name {
            &.read-only {
              // font-family: "SSP Semibold Italic"; // color: $a-grey-400;
              font-family: "SSP Semibold";
            }
          }
          .info-marker-container {
            padding: 0 3px;
            padding-right: 7px;
            margin-top: 1px;
            fill: $a-grey-400;
          }

          .profile-caret-icon {
            display: flex;
            padding-left: 8px;
            padding-top: 2px;
            svg {
              margin: auto;
            }
          }
        }
        .dropdown {
          margin: auto;
          .dropdown-divider {
            background-color: $dropdown-divider-color;
          }
        }
      }
    }
  }

  .button-container {
    width: 100%;
    display: flex;
    padding: 5px 15px;
    .join-button {
      flex: 1;
      font-family: "SCP Bold";
      font-size: 14px;
    }
  }

  .mod-puck-container {
    display: flex;
    .mod-puck {
      font-family: "SCP Bold";
      letter-spacing: 2px;
      font-size: 90%;
      border-radius: 5px;
      padding: 0 3px 0 6px;
      margin: auto 0px auto 8px;
      border: 1px solid $a-purple;
      color: $a-purple;
      line-height: 135%;
    }
  }
</style>