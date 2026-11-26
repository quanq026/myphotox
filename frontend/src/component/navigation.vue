<template>
  <div id="p-navigation" :class="{ 'sidenav-visible': drawer }">
    <template v-if="visible">
      <template v-if="$vuetify.display.smAndDown">
        <v-toolbar position="fixed" flat density="compact" color="navigation" class="nav-small elevation-2" @click.stop.prevent>
          <v-btn icon variant="text" class="bg-transparent nav-logo" @click.stop.prevent="toggleDrawer">
            <img :src="appIcon" :alt="appName" :class="{ 'animate-hue': indexing }" />
          </v-btn>
          <v-toolbar-title class="nav-toolbar-title">
            <span :class="{ clickable: auth }" @click.stop.prevent.self="toggleDrawer">{{ page.title }}</span>
          </v-toolbar-title>
          <v-btn
            icon="mdi-dots-vertical"
            variant="text"
            class="nav-mobile-menu-trigger elevation-0"
            :ripple="false"
            @click.stop.prevent="speedDial = true"
          ></v-btn>
        </v-toolbar>
      </template>
      <template v-else-if="!auth">
        <v-toolbar flat :density="$vuetify.display.smAndDown ? 'compact' : 'default'" color="navigation" class="nav-small">
          <v-btn icon variant="text" class="bg-transparent nav-logo">
            <img :src="appIcon" :alt="appName" />
          </v-btn>
          <v-toolbar-title class="nav-toolbar-title">
            {{ page.title }}
          </v-toolbar-title>
          <v-btn
            icon="mdi-dots-vertical"
            variant="text"
            class="nav-mobile-menu-trigger elevation-0"
            :ripple="false"
            @click.stop.prevent="speedDial = true"
          ></v-btn>
        </v-toolbar>
      </template>
      <v-navigation-drawer v-if="auth" v-model="drawer" :rail="isMini" color="navigation" class="nav-sidebar navigation">
        <div class="nav-container">
          <v-toolbar flat :density="$vuetify.display.smAndDown ? 'compact' : 'default'">
            <v-list class="navigation-home elevation-0" bg-color="navigation-home" width="100%" density="compact" tabindex="-1">
              <v-list-item class="px-3" :elevation="0" :ripple="false" @click.stop.prevent="onHome">
                <template #prepend>
                  <div class="v-avatar bg-transparent nav-logo">
                    <a :href="siteUrl" tabindex="-1" @click.stop.prevent="onHome">
                      <img :src="appIcon" :alt="appName" :class="{ 'animate-hue': indexing }" />
                    </a>
                  </div>
                </template>
                <template #append>
                  <v-btn
                    icon
                    variant="text"
                    :elevation="0"
                    tabindex="-1"
                    class="nav-minimize hidden-sm-and-down"
                    :ripple="false"
                    :title="$gettext('Minimize')"
                    @click.stop.prevent="toggleIsMini()"
                  >
                    <v-icon :icon="rtl ? 'mdi-chevron-right' : 'mdi-chevron-left'"></v-icon>
                  </v-btn>
                </template>
                <v-list-item-title class="nav-toolbar-title">
                  <a :href="siteUrl">{{ appName }}</a>
                </v-list-item-title>
              </v-list-item>
            </v-list>
          </v-toolbar>

          <v-list
            nav
            class="nav-menu"
            bg-color="navigation"
            color="primary"
            open-strategy="single"
            :density="$vuetify.display.smAndDown ? 'compact' : 'comfortable'"
            tabindex="-1"
          >
            <v-list-item v-if="isMini && !isRestricted" class="nav-expand" @click.stop="toggleIsMini()">
              <v-icon :icon="rtl ? 'mdi-chevron-left' : 'mdi-chevron-right'" class="ma-auto"></v-icon>
            </v-list-item>

            <v-list-item v-if="isMini && $config.feature('search')" to="/browse" variant="text" class="nav-browse" :ripple="false" @click.stop="">
              <v-icon class="ma-auto">mdi-magnify</v-icon>
            </v-list-item>
            <v-list-item
              v-else-if="!isMini && $config.feature('search')"
              to="/browse"
              variant="text"
              class="nav-browse"
              :ripple="false"
              @click.stop=""
            >
              <v-list-item-title class="nav-menu-item">
                <v-icon>mdi-magnify</v-icon>
                <p class="nav-item-title">
                  {{ $gettext(`Search`) }}
                </p>
              </v-list-item-title>
              <span v-if="config.count.all > 0" class="nav-count-item">{{ config.count.all }}</span>
            </v-list-item>

            <v-list-item v-if="isMini" v-show="$config.feature('albums')" to="/albums" variant="text" class="nav-albums" :ripple="false" @click.stop="">
              <v-icon class="ma-auto">mdi-bookmark</v-icon>
            </v-list-item>
            <v-list-item
              v-else-if="!isMini"
              v-show="$config.feature('albums')"
              to="/albums"
              variant="text"
              class="nav-albums"
              :ripple="false"
              @click.stop=""
            >
              <v-list-item-title class="nav-menu-item">
                <v-icon>mdi-bookmark</v-icon>
                <p class="nav-item-title">
                  {{ $gettext(`Albums`) }}
                </p>
              </v-list-item-title>
              <span v-if="config.count.albums > 0" class="nav-count-item">{{ config.count.albums }}</span>
            </v-list-item>

            <v-list-item
              v-if="isMini && $config.feature('videos')"
              :to="{ name: 'videos' }"
              variant="text"
              class="nav-video nav-videos"
              :ripple="false"
              @click.stop=""
            >
              <v-icon class="ma-auto">mdi-play-circle</v-icon>
            </v-list-item>
            <v-list-item
              v-else-if="!isMini && $config.feature('videos')"
              :to="{ name: 'videos' }"
              variant="text"
              class="nav-video nav-videos"
              :ripple="false"
              @click.stop=""
            >
              <v-list-item-title class="nav-menu-item">
                <v-icon>mdi-play-circle</v-icon>
                <p class="nav-item-title">
                  {{ $gettext(`Videos`) }}
                </p>
              </v-list-item-title>
              <span v-show="config.count.videos > 0" class="nav-count-item">{{ config.count.videos }}</span>
            </v-list-item>

            <v-list-item
              v-if="isMini && $config.feature('favorites')"
              :to="{ name: 'favorites' }"
              variant="text"
              class="nav-favorites"
              :ripple="false"
              @click.stop=""
            >
              <v-icon class="ma-auto">mdi-star</v-icon>
            </v-list-item>
            <v-list-item
              v-else-if="!isMini && $config.feature('favorites')"
              :to="{ name: 'favorites' }"
              variant="text"
              class="nav-favorites"
              :ripple="false"
              @click.stop=""
            >
              <v-list-item-title class="nav-menu-item">
                <v-icon>mdi-star</v-icon>
                <p class="nav-item-title">
                  {{ $gettext(`Favorites`) }}
                </p>
              </v-list-item-title>
              <span v-show="config.count.favorites > 0" class="nav-count-item">{{ config.count.favorites }}</span>
            </v-list-item>

            <v-list-item
              v-if="isMini && $config.feature('folders')"
              :to="{ name: 'folders' }"
              variant="text"
              class="nav-folders"
              :ripple="false"
              @click.stop=""
            >
              <v-icon class="ma-auto">mdi-folder</v-icon>
            </v-list-item>
            <v-list-item
              v-else-if="!isMini && $config.feature('folders')"
              :to="{ name: 'folders' }"
              variant="text"
              class="nav-folders"
              :ripple="false"
              @click.stop=""
            >
              <v-list-item-title class="nav-menu-item">
                <v-icon>mdi-folder</v-icon>
                <p class="nav-item-title">
                  {{ $gettext(`Folders`) }}
                </p>
              </v-list-item-title>
              <span v-show="config.count.folders > 0" class="nav-count-item">{{ config.count.folders }}</span>
            </v-list-item>

            <v-list-item
              v-if="isMini && $config.feature('archive')"
              :to="{ name: 'archive' }"
              variant="text"
              class="nav-archive"
              :ripple="false"
              @click.stop=""
            >
              <v-icon class="ma-auto">mdi-package-down</v-icon>
            </v-list-item>
            <v-list-item
              v-else-if="!isMini && $config.feature('archive')"
              :to="{ name: 'archive' }"
              variant="text"
              class="nav-archive"
              :ripple="false"
              @click.stop=""
            >
              <v-list-item-title class="nav-menu-item">
                <v-icon>mdi-package-down</v-icon>
                <p class="nav-item-title">
                  {{ $pgettext("Noun", "Archive") }}
                </p>
              </v-list-item-title>
              <span v-show="config.count.archived > 0" class="nav-count-item">{{ config.count.archived }}</span>
            </v-list-item>

            <v-list-item
              v-if="isMini && $config.feature('library')"
              :to="{ name: 'library_index' }"
              variant="text"
              class="nav-library"
              :ripple="false"
              @click.stop=""
            >
              <v-icon class="ma-auto">mdi-film</v-icon>
            </v-list-item>
            <v-list-item
              v-else-if="!isMini && $config.feature('library')"
              :to="{ name: 'library_index' }"
              variant="text"
              class="nav-library"
              :ripple="false"
              @click.stop=""
            >
              <v-list-item-title class="nav-menu-item">
                <v-icon>mdi-film</v-icon>
                <p class="nav-item-title">
                  {{ $gettext(`Library`) }}
                </p>
              </v-list-item-title>
            </v-list-item>

            <template v-if="!config.disable.settings">
              <v-list-item
                v-if="isMini"
                v-show="$config.feature('settings')"
                :to="{ name: 'settings' }"
                variant="text"
                class="nav-settings"
                :ripple="false"
                @click.stop=""
              >
                <v-icon class="ma-auto">mdi-cog</v-icon>
              </v-list-item>
              <div v-else-if="!isMini" v-show="$config.feature('settings')">
                <v-list-item :to="{ name: 'settings' }" variant="text" class="nav-settings activator" @click.stop="">
                  <v-list-item-title class="nav-menu-item">
                    <p class="nav-item-title">
                      {{ $gettext(`Settings`) }}
                    </p>
                  </v-list-item-title>
                </v-list-item>

                <v-list-group>
                  <template #activator="{ props }">
                    <v-list-item v-bind="props" variant="text" class="nav-settings" @click.stop="">
                      <v-icon>mdi-cog</v-icon>
                    </v-list-item>
                  </template>

                  <v-list-item :to="{ name: 'about' }" :exact="true" variant="text" class="nav-about" :ripple="false" @click.stop="">
                    <v-list-item-title :class="`menu-item`">
                      {{ $gettext(`About`) }}
                    </v-list-item-title>
                  </v-list-item>
                </v-list-group>
              </div>
            </template>

            <v-list-item v-show="!auth" :to="{ name: 'login' }" variant="text" class="nav-login" @click.stop="">
              <v-list-item-title>
                <v-icon>mdi-lock</v-icon>
                {{ $gettext(`Login`) }}
              </v-list-item-title>
            </v-list-item>
          </v-list>

          <div v-if="disconnected" class="nav-info connection-info clickable" @click.stop="showServerConnectionHelp">
            <div class="nav-info__underlay"></div>
            <div class="text-center">
              <v-icon icon="mdi-wifi-off" color="warning" size="21"></v-icon>
            </div>
            <div v-if="!isMini" class="text-start text-body-2">
              {{ $gettext(`Lost server connection`) }}
            </div>
          </div>
          <div v-else-if="!isMini && featUsage" class="nav-info usage-info clickable" @click.stop="showUsageInfo">
            <div class="nav-info__underlay"></div>
            <div class="nav-info__content">
              <v-progress-linear
                :model-value="config.usage.filesUsedPct"
                :color="config.usage.filesUsedPct > 95 ? 'error' : 'selected'"
                height="16"
                max="100"
                min="0"
                width="100%"
                rounded
              >
                <div class="text-caption opacity-85">
                  {{
                    $gettext(`%{n} GB of %{q} GB used`, {
                      n: $util.gigaBytes(config.usage.filesUsed),
                      q: $util.gigaBytes(config.usage.filesTotal),
                    })
                  }}
                </div>
              </v-progress-linear>
            </div>
          </div>

          <div v-show="auth && !isPublic && !disconnected" class="nav-info user-info">
            <div class="nav-info__underlay"></div>
            <p-auth-menu @account="onAccount" @logout="onLogout">
              <div class="nav-user-avatar text-center my-1 mx-2">
                <img :src="userAvatarURL" :alt="accountInfo" :title="accountInfo" class="rounded-circle" />
              </div>
              <template v-if="!isMini">
                <div class="nav-user-text text-start mt-1 flex-grow-1">
                  <p class="text-body-2"><bdi>{{ $util.balanceBidi(displayName) }}</bdi></p>
                  <p class="text-caption opacity-70">{{ accountInfo }}</p>
                </div>
                <div class="text-center">
                  <v-btn icon="mdi-dots-vertical" variant="text" :elevation="0"></v-btn>
                </div>
              </template>
            </p-auth-menu>
          </div>
        </div>
      </v-navigation-drawer>
      <div v-if="config.legalInfo" id="legal-info">
        <span v-if="config.legalUrl" class="clickable" @click.stop.prevent="showLegalInfo()">{{ config.legalInfo }}</span>
        <span v-else>{{ config.legalInfo }}</span>
      </div>
    </template>
    <div id="mobile-menu" :class="{ active: speedDial }" @click.stop="speedDial = false">
      <div class="menu-content grow-top-end">
        <div class="menu-icons">
          <a v-if="auth && !isPublic" href="#" :title="$gettext('Sign Out')" class="menu-action navigation-logout" @click.prevent="onLogout">
            <v-icon>mdi-power</v-icon>
          </a>
          <a href="#" :title="$gettext('Reload')" class="menu-action nav-reload" @click.prevent="reloadApp">
            <v-icon>mdi-refresh</v-icon>
          </a>
          <router-link
            v-if="auth && $config.feature('account')"
            :to="{ name: 'settings_account' }"
            :title="$gettext('Account')"
            class="menu-action nav-account"
          >
            <v-icon>mdi-shield-account-variant</v-icon>
          </router-link>
          <router-link
            v-if="auth && $config.feature('settings') && !routeName('settings')"
            :to="{ name: 'settings' }"
            :title="$gettext('Settings')"
            class="menu-action nav-settings"
          >
            <v-icon>mdi-cog</v-icon>
          </router-link>
          <a
            v-if="auth && !config.readonly && $config.feature('upload')"
            href="#"
            :title="$gettext('Upload')"
            class="menu-action nav-upload"
            @click.prevent="openUpload()"
          >
            <v-icon>mdi-cloud-upload</v-icon>
          </a>
          <router-link v-if="!auth && !isPublic" :to="{ name: 'login' }" :title="$gettext('Login')" class="menu-action nav-login">
            <v-icon>mdi-login</v-icon>
          </router-link>
        </div>
        <div class="menu-actions">
          <div v-if="auth && !routeName('browse') && $config.feature('search')" class="menu-action nav-search">
            <router-link to="/browse">
              <v-icon>mdi-magnify</v-icon>
              {{ $gettext(`Search`) }}
            </router-link>
          </div>
          <div v-if="auth && !routeName('albums') && $config.feature('albums')" class="menu-action nav-albums">
            <router-link to="/albums">
              <v-icon>mdi-bookmark</v-icon>
              {{ $gettext(`Albums`) }}
            </router-link>
          </div>
          <div v-if="auth && !routeName('videos') && $config.feature('videos')" class="menu-action nav-videos">
            <router-link :to="{ name: 'videos' }">
              <v-icon>mdi-play-circle</v-icon>
              {{ $gettext(`Videos`) }}
            </router-link>
          </div>
          <div v-if="auth && !routeName('favorites') && $config.feature('favorites')" class="menu-action nav-favorites">
            <router-link :to="{ name: 'favorites' }">
              <v-icon>mdi-star</v-icon>
              {{ $gettext(`Favorites`) }}
            </router-link>
          </div>
          <div v-if="auth && !routeName('folders') && $config.feature('folders')" class="menu-action nav-folders">
            <router-link :to="{ name: 'folders' }">
              <v-icon>mdi-folder</v-icon>
              {{ $gettext(`Folders`) }}
            </router-link>
          </div>
          <div v-if="auth && !routeName('library_index') && $config.feature('library')" class="menu-action nav-index">
            <router-link :to="{ name: 'library_index' }">
              <v-icon>mdi-film</v-icon>
              {{ $gettext(`Library`) }}
            </router-link>
          </div>
          <div v-if="config.legalUrl" class="menu-action nav-legal">
            <a :href="config.legalUrl" target="_blank" rel="noopener">
              <v-icon>mdi-information</v-icon>
              {{ $gettext(`Legal Information`) }}
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import links from "common/links";
import { getAppStorage } from "common/storage";
import PAuthMenu from "component/auth/menu.vue";

const appStorage = getAppStorage();

export default {
  name: "PNavigation",
  components: {
    PAuthMenu,
  },
  data() {
    const appName = this.$config.getName();

    let appNameSuffix = "";
    let appNameParts = appName.split(" ");

    if (appNameParts.length > 1) {
      appNameSuffix = appNameParts.slice(1, 9).join(" ");
    }

    const canManagePhotos = this.$config.allow("photos", "manage");
    const canSeeUsage = this.$session.isAdmin() || (canManagePhotos && this.$config.feature("files"));
    const isDemo = this.$config.get("demo");
    const isPro = this.$config.isPro();
    const isPublic = this.$config.get("public");
    const isReadOnly = this.$config.get("readonly");
    const isRestricted = this.$config.deny("photos", "access_library");
    const isSuperAdmin = this.$session.isSuperAdmin();
    const tier = this.$config.getTier();

    return {
      links,
      canSearchPlaces: this.$config.allow("places", "search"),
      canAccessPrivate: !isRestricted && this.$config.allow("photos", "access_private"),
      canManagePhotos: canManagePhotos,
      canManagePeople: this.$config.allow("people", "manage"),
      canManageUsers: (!isPublic || isDemo) && this.$config.allow("users", "access_all"),
      canManageServices: this.$config.feature("services") && this.$config.allow("services", "manage"),
      appNameSuffix: appNameSuffix,
      appName: this.$config.getName(),
      appAbout: this.$config.getAbout(),
      appIcon: this.$config.getIcon(),
      siteUrl: this.$config.values.siteUrl,
      disconnected: this.$config.disconnected,
      indexing: false,
      drawer: null,
      featUpgrade: false,
      featMembership: false,
      featFiles: this.$config.feature("files"),
      featUsage: canSeeUsage && this.$config.values?.usage?.filesTotal,
      isRestricted: isRestricted,
      isMini: appStorage.getItem("navigation.mode") !== "false" || isRestricted,
      isDemo: isDemo,
      isPro: isPro,
      isPublic: isPublic,
      isReadOnly: isReadOnly,
      isAdmin: this.$session.isAdmin(),
      isSuperAdmin: isSuperAdmin,
      isSponsor: this.$config.isSponsor(),
      isTest: this.$config.test,
      session: this.$session,
      config: this.$config.values,
      page: this.$config.page,
      user: this.$session.getUser(),
      speedDial: false,
      rtl: this.$isRtl,
      subscriptions: [],
    };
  },
  computed: {
    auth() {
      return this.session.auth || this.isPublic;
    },
    visible() {
      return !this.$route.meta.hideNav;
    },
    displayName() {
      const user = this.$session.getUser();
      if (user) {
        return user.getDisplayName();
      }

      return this.$gettext("Unregistered");
    },
    userAvatarURL() {
      return this.$session.getUser().getAvatarURL("tile_50", this.$config);
    },
    accountInfo() {
      const user = this.$session.getUser();
      if (user) {
        return user.getAccountInfo();
      }

      return this.$gettext("Account");
    },
  },
  created() {
    this.subscriptions.push(this.$event.subscribe("index", this.onIndex));
    this.subscriptions.push(this.$event.subscribe("import", this.onIndex));
  },
  beforeUnmount() {
    for (let i = 0; i < this.subscriptions.length; i++) {
      this.$event.unsubscribe(this.subscriptions[i]);
    }
  },
  methods: {
    routeName(name) {
      if (!name || !this.$route.name) {
        return false;
      }

      return this.$route.name.startsWith(name);
    },
    reloadApp() {
      this.$notify.info(this.$gettext("Reloading…"));
      this.$notify.blockUI();
      setTimeout(() => window.location.reload(), 100);
    },
    openUpload() {
      this.$event.publish("dialog.upload");
    },
    onHome(ev) {
      if (this.$vuetify.display.smAndDown) {
        this.toggleDrawer(ev);
        return;
      }

      if (this.$route.name !== "home") {
        this.$router.push({ name: "home" });
      }
    },
    showDrawer() {
      if (this.auth) {
        this.drawer = true;
        this.isMini = this.isRestricted;
      }
    },
    hideDrawer() {
      if (this.auth) {
        this.drawer = false;
        this.isMini = this.isRestricted;
      }
    },
    toggleDrawer(ev) {
      if (!ev || ev.target === undefined) {
        return;
      }

      if (this.$vuetify.display.smAndDown) {
        if (this.drawer) {
          this.hideDrawer();
        } else {
          this.showDrawer();
        }
      } else {
        this.toggleIsMini();
      }
    },
    toggleIsMini() {
      if (this.isRestricted) {
        return;
      }

      this.isMini = !this.isMini;
      appStorage.setItem("navigation.mode", `${this.isMini}`);
    },
    showUsageInfo() {
      this.$router.push({ name: "library_index" });
    },
    showServerConnectionHelp() {
      this.reloadApp();
    },
    showLegalInfo() {
      if (this.config.legalUrl) {
        this.$util.openExternalUrl(this.config.legalUrl);
      } else {
        this.$router.push({ name: "about" });
      }
    },
    // onAccount opens the account settings tab, falling back to general settings
    // when the account feature is unavailable.
    onAccount() {
      this.$router.push({ name: this.$config.feature("account") ? "settings_account" : "settings" });
    },
    onLogout() {
      this.$session.logoutEverywhere();
    },
    onIndex(ev) {
      if (!ev) {
        return;
      }

      const type = ev.split(".")[1];

      switch (type) {
        case "file":
        case "folder":
        case "indexing":
          this.indexing = true;
          break;
        case "completed":
          this.indexing = false;
          break;
      }
    },
  },
};
</script>
