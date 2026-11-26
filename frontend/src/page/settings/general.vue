<template>
  <div class="p-tab p-settings-general py-2">
    <v-form ref="form" validate-on="invalid-input" class="p-form-settings" accept-charset="UTF-8" @submit.prevent="onChange">
      <v-card flat tile class="mt-0 px-1 bg-background">
        <v-card-title class="pb-2 text-subtitle-2">
          {{ $gettext(`User Interface`) }}
        </v-card-title>

        <v-card-actions>
          <v-row class="align-start" density="compact">
            <v-col cols="12" sm="4">
              <v-select
                v-model="settings.ui.theme"
                :disabled="busy"
                :items="themes"
                item-title="text"
                item-value="value"
                :label="$gettext('Theme')"
                :menu-props="{ maxHeight: 346 }"
                hide-details
                class="input-theme"
                @update:model-value="onChangeTheme"
              ></v-select>
            </v-col>

            <v-col cols="12" sm="4">
              <v-select
                v-model="settings.ui.language"
                :disabled="busy"
                :items="languages"
                item-title="text"
                item-value="value"
                :label="$gettext('Language')"
                :menu-props="{ maxHeight: 346 }"
                hide-details
                class="input-language"
                @update:model-value="onChange"
              ></v-select>
            </v-col>

            <v-col cols="12" sm="4">
              <v-select
                v-model="settings.ui.timeZone"
                :disabled="busy"
                item-value="ID"
                item-title="Name"
                :items="options.TimeZones($gettext('Local'))"
                :label="$gettext('Time Zone')"
                :menu-props="{ maxHeight: 346 }"
                hide-details
                class="input-timezone"
                @update:model-value="onChangeTheme"
              ></v-select>
            </v-col>
          </v-row>
        </v-card-actions>
      </v-card>

      <v-card flat tile class="mt-2 px-1 bg-background">
        <v-card-title class="pb-2 text-subtitle-2">
          {{ $gettext(`Display`) }}
        </v-card-title>

        <v-card-actions>
          <v-row class="align-start" density="compact">
            <v-col cols="12" sm="4" class="px-2 pb-2 pt-2">
              <v-checkbox
                v-model="settings.search.listView"
                :disabled="busy"
                class="ma-0 pa-0 input-search-listview"
                density="compact"
                :label="$gettext('List View')"
                :hint="$gettext('View search results as a list.')"
                prepend-icon="mdi-view-list"
                persistent-hint
                @update:model-value="onChange"
              ></v-checkbox>
            </v-col>

            <v-col cols="12" sm="4" class="px-2 pb-2 pt-2">
              <v-checkbox
                v-model="settings.search.showTitles"
                :disabled="busy"
                class="ma-0 pa-0 input-search-titles"
                density="compact"
                :label="$gettext('Show Titles')"
                :hint="$gettext('Display picture titles in search results.')"
                prepend-icon="mdi-format-text"
                persistent-hint
                @update:model-value="onChange"
              ></v-checkbox>
            </v-col>

            <v-col cols="12" sm="4" class="px-2 pb-2 pt-2">
              <v-checkbox
                v-model="settings.search.showCaptions"
                :disabled="busy"
                class="ma-0 pa-0 input-search-captions"
                density="compact"
                :label="$gettext('Show Captions')"
                :hint="$gettext('Display picture captions in search results.')"
                prepend-icon="mdi-text"
                persistent-hint
                @update:model-value="onChange"
              ></v-checkbox>
            </v-col>
          </v-row>
        </v-card-actions>
      </v-card>

      <v-card v-if="canManageLibrary" flat tile class="mt-2 px-1 bg-background">
        <v-card-title class="pb-2 text-subtitle-2">
          {{ $gettext(`Library`) }}
        </v-card-title>

        <v-card-actions class="px-2 pb-3 pt-1">
          <div class="d-flex flex-wrap align-center ga-3">
            <v-btn
              variant="flat"
              color="highlight"
              class="action-rescan"
              :loading="indexing"
              :disabled="indexing"
              prepend-icon="mdi-refresh"
              @click.stop="startRescan"
            >
              {{ $gettext(`Rescan Library`) }}
            </v-btn>
            <v-btn
              v-if="indexing"
              variant="flat"
              color="button"
              class="action-cancel-rescan"
              @click.stop="cancelRescan"
            >
              {{ $gettext(`Cancel`) }}
            </v-btn>
            <span v-if="indexStatus" class="text-caption text-secondary-dark text-break">{{ indexStatus }}</span>
          </div>
        </v-card-actions>
      </v-card>
    </v-form>
    <p-about-footer></p-about-footer>
  </div>
</template>

<script>
import Settings from "model/settings";
import * as options from "options/options";
import * as themes from "options/themes";
import $api from "common/api";
import PAboutFooter from "component/about/footer.vue";

export default {
  name: "PSettingsGeneral",
  components: {
    PAboutFooter,
  },
  data() {
    return {
      isPublic: this.$config.isPublic(),
      canManageLibrary: !this.$config.get("readonly") && !this.$config.isPublic() && this.$config.allow("files", "manage"),
      settings: new Settings(this.$config.getSettings()),
      options: options,
      busy: this.$config.loading(),
      indexing: false,
      indexStatus: "",
      subscriptions: [],
      themes: [],
      currentTheme: this.$config.themeName,
      languages: options.Languages(),
    };
  },
  created() {
    this.load();
    this.subscriptions.push(this.$event.subscribe("config.updated", (ev, data) => this.settings.setValues(data.config.settings)));
    this.subscriptions.push(this.$event.subscribe("index", this.onIndexEvent));
  },
  beforeUnmount() {
    for (let i = 0; i < this.subscriptions.length; i++) {
      this.$event.unsubscribe(this.subscriptions[i]);
    }
  },
  methods: {
    load() {
      this.busy = true;
      this.$notify.blockUI("busy");

      this.$config
        .load()
        .then(() => {
          this.settings.setValues(this.$config.getSettings());
          this.themes = themes.Translated().filter(Boolean);
          this.currentTheme = this.settings.ui.theme;
        })
        .finally(() => {
          this.busy = false;
          this.$notify.unblockUI();
        });
    },
    onChangeTheme() {
      this.onChange().then(() => {
        if (this.currentTheme !== this.settings.ui.theme || this.settings.changed("ui", "timeZone")) {
          this.$notify.info(this.$gettext("Reloading…"));
          this.$notify.blockUI();
          setTimeout(() => window.location.reload(), 100);
        }
      });
    },
    onChange() {
      if (this.busy) {
        return Promise.resolve();
      }

      const reload = this.settings.changed("ui", "language");

      this.busy = true;
      this.$notify.blockUI("busy");

      return this.settings
        .save()
        .then(() => {
          this.$config.setSettings(this.settings);

          if (reload) {
            this.$notify.info(this.$gettext("Reloading…"));
            setTimeout(() => window.location.reload(), 100);
          } else {
            this.$notify.success(this.$gettext("Changes successfully saved"));
          }
        })
        .finally(() => {
          this.busy = false;
          this.$notify.unblockUI();
        });
    },
    startRescan() {
      if (this.indexing || !this.canManageLibrary) {
        return;
      }
      this.indexing = true;
      this.indexStatus = this.$gettext("Indexing…");

      $api
        .post("index", { path: "", rescan: false, cleanup: true })
        .then(() => {
          this.indexing = false;
          this.indexStatus = "";
          this.$notify.success(this.$gettext("Done."));
        })
        .catch((e) => {
          if (e?.response?.status !== 409) {
            this.$notify.error(this.$gettext("Indexing failed"));
          }
          this.indexing = false;
          this.indexStatus = "";
        });
    },
    cancelRescan() {
      $api.delete("index");
    },
    onIndexEvent(ev, data) {
      const type = ev ? ev.split(".")[1] : "";
      switch (type) {
        case "folder":
        case "indexing":
          this.indexing = true;
          this.indexStatus = data?.fileName || data?.filePath || this.$gettext("Indexing…");
          break;
        case "converting":
        case "thumbnails":
        case "updating":
          this.indexing = true;
          this.indexStatus = data?.fileName || this.$gettext("Updating index");
          break;
        case "completed":
          this.indexing = false;
          this.indexStatus = "";
          break;
      }
    },
  },
};
</script>
