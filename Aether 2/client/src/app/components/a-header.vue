<template>
  <div class="header-container" :class="{'sidebar-closed': !sidebarOpen}">
    <div class="history-container">
      <a-header-icon icon="chevron-left" @click.native="$router.go(-1)"></a-header-icon>
      <a-header-icon icon="chevron-right" @click.native="$router.go(+1)"></a-header-icon>
      <!-- <div @click="$router.go(-1)">test</div> -->
    </div>
    <div class="breadcrumbs-container">
      <a-breadcrumbs></a-breadcrumbs>
    </div>
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
                <a-info-marker header="You are in read only mode." text="You either haven't created a user yet, or have disabled it by going into read only mode. <br><br>You can read, but you won't be able to participate until you create or re-enable your user key."></a-info-marker>
              </div>
              <div class="user-name" :class="{'read-only': this.localUserReadOnly}">
                {{localUserName}}
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
                  <router-link class="button is-primary is-outlined is-small join-button" to="/newuser">
                    JOIN
                  </router-link>
                </div>
                <hr class="dropdown-divider">
              </template>
              <template v-if="!localUserReadOnly">
                <a-avatar-block :user="$store.state.localUser" :clickable="true"></a-avatar-block>
                <hr class="dropdown-divider">
                <router-link :to="'/user/'+$store.state.localUser.fingerprint" class="dropdown-item">
                  Profile
                </router-link>
              </template>
              <router-link to="/intro" class="dropdown-item">
                I'm new here
              </router-link>
              <router-link to="/settings" class="dropdown-item">
                Settings
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
        userMenuOpen: false
      };
    },
    computed: {
      ...Vuex.mapState(['sidebarOpen']),
      localUserName(this: any) {
        if (this.localUserReadOnly) {
          return "Read only"
        }
        return globalMethods.GetUserName(this.$store.state.localUser)
      }
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
      signOut(this: any) {}
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
    }

    .breadcrumbs-container {
      flex: 1;
      min-width: 0;
    }

    .profile-container {
      display: flex;
      padding: 0 15px;

      &:hover {
        background-color: rgba(255, 255, 255, 0.05);
      }

      .dropdown-container {
        display: flex;

        .dropdown-trigger {
          display: flex;
          cursor: pointer;

          .user-name {
            &.read-only {
              font-family: "SSP Semibold Italic"; // color: $a-grey-400;
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
    }
  }
</style>