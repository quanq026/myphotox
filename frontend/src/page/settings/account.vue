<template>
  <div class="p-tab p-settings-account">
    <div class="width-lg pa-3">
      <v-form ref="form" v-model="valid" class="p-form-account ma-0 pa-0" accept-charset="UTF-8" @submit.prevent="onChange">
        <input ref="upload" type="file" class="d-none input-upload" accept="image/png, image/jpeg" @change.stop="onUploadAvatar()" />
        <v-card flat tile class="bg-background ma-0 pa-0">
          <v-card-actions class="ma-0 pa-0">
            <v-row class="align-start" density="compact">
              <v-col cols="8" sm="9" md="10" class="pa-0 d-flex align-self-stretch">
                <v-row class="align-start" density="compact">
                  <v-col cols="12" md="5">
                    <v-text-field
                      v-model="user.DisplayName"
                      :disabled="busy"
                      maxlength="200"
                      autocomplete="off"
                      autocorrect="off"
                      autocapitalize="none"
                      :label="$gettext('Display Name')"
                      class="input-display-name"
                      :rules="rules.text(true, 1, 200, $gettext('Display Name'))"
                      @change="onChange"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="12" md="7">
                    <v-text-field
                      v-model="user.Email"
                      type="email"
                      maxlength="255"
                      :disabled="busy || !canEditEmail"
                      autocomplete="off"
                      autocorrect="off"
                      autocapitalize="none"
                      :label="$gettext('Email')"
                      class="input-email"
                      :rules="rules.email()"
                      @change="onChange"
                    ></v-text-field>
                  </v-col>
                </v-row>
              </v-col>
              <v-col class="text-center align-self-center" cols="4" sm="3" md="2">
                <v-avatar :size="$vuetify.display.md ? 100 : 112" :class="{ clickable: !busy }" @click.stop.prevent="onChangeAvatar()">
                  <v-img
                    :alt="accountInfo"
                    :title="$gettext('Change Avatar')"
                    :src="$vuetify.display.md ? user.getAvatarURL('tile_100') : user.getAvatarURL('tile_224')"
                  ></v-img>
                </v-avatar>
              </v-col>
              <v-col v-if="user.Details.Bio" cols="12">
                <v-textarea
                  v-model="user.Details.Bio"
                  auto-grow
                  rows="2"
                  class="input-bio"
                  autocorrect="off"
                  autocapitalize="none"
                  autocomplete="off"
                  :disabled="busy"
                  maxlength="2000"
                  :rules="rules.text(false, 1, 2000, $gettext('Bio'))"
                  :label="$gettext('Bio')"
                  @change="onChange"
                ></v-textarea>
              </v-col>
              <v-col cols="12">
                <v-textarea
                  v-model="user.Details.About"
                  auto-grow
                  rows="2"
                  class="input-about"
                  autocorrect="off"
                  autocapitalize="none"
                  autocomplete="off"
                  :disabled="busy"
                  maxlength="500"
                  :rules="rules.text(false, 1, 500, $gettext('About'))"
                  :label="$gettext('About')"
                  @change="onChange"
                ></v-textarea>
              </v-col>
            </v-row>
          </v-card-actions>
        </v-card>
        <v-card flat tile class="my-3 pa-0 bg-background">
          <v-card-title class="ma-0 pa-2 text-subtitle-2">
            {{ $gettext(`Security and Access`) }}
          </v-card-title>
          <v-card-actions class="ma-0 pa-0">
            <v-row class="align-start" density="compact">
              <v-col cols="12" sm="6">
                <v-btn
                  block
                  variant="flat"
                  color="button"
                  class="action-change-password"
                  :disabled="isPublic || isDemo || user.Name === '' || getProvider() !== 'local'"
                  @click.stop="showDialog('password')"
                >
                  {{ $gettext(`Change Password`) }}
                  <v-icon end>mdi-lock</v-icon>
                </v-btn>
              </v-col>
            </v-row>
          </v-card-actions>
        </v-card>
      </v-form>
    </div>
    <p-settings-password :visible="dialog.password" :model="user" @close="dialog.password = false"></p-settings-password>
  </div>
</template>

<script>
import PSettingsPassword from "component/settings/password.vue";
import User from "model/user";
import { rules } from "common/form";

export default {
  name: "PSettingsAccount",
  components: {
    PSettingsPassword,
  },
  data() {
    const isDemo = this.$config.isDemo();
    const isPublic = this.$config.isPublic();
    const user = this.$session.getUser();

    return {
      busy: isDemo || isPublic,
      rules,
      isDemo,
      isPublic,
      valid: true,
      rtl: this.$isRtl,
      user: user,
      session: this.$session,
      dialog: {
        password: false,
      },
    };
  },
  computed: {
    // showConnectedApps reports whether users can connect apps with OAuth sign-in, so the dialog is offered.
    showConnectedApps() {
      return this.$config.featOAuthApps();
    },
    accountInfo() {
      const user = this.$session.getUser();
      if (user) {
        return user.getAccountInfo();
      }

      return this.$gettext("Unregistered");
    },
    displayName() {
      const user = this.$session.getUser();
      if (user) {
        return user.getDisplayName();
      }

      return this.$gettext("Unregistered");
    },
    // canEditEmail checks if the account may change its email address, which requires permission to manage accounts.
    canEditEmail() {
      const user = this.$session.getUser();
      return !!user && (user.SuperAdmin || this.$config.allowAny("users", ["manage", "manage_own"]));
    },
  },
  created() {
    if (this.isPublic && !this.isDemo) {
      this.$router.push({ name: "settings" });
    }
  },
  mounted() {
    this.$refs.form.validate();
  },
  methods: {
    getProvider() {
      return this.$session.provider ? this.$session.provider : this.user.AuthProvider;
    },
    showDialog(name) {
      if (!name) {
        return;
      }
      this.dialog[name] = true;
    },
    updateUser() {
      this.$notify.blockUI("busy");
      this.$session
        .refresh()
        .then(() => {
          this.user = this.$session.getUser();
        })
        .finally(() => {
          this.$notify.unblockUI();
        });
    },
    onChangeAvatar() {
      if (this.busy) {
        return;
      }
      this.$refs.upload.click();
    },
    onChangeName() {
      this.user.Details.NameSrc = "manual";
      return this.onChange();
    },
    onChange() {
      if (this.busy || !this?.$refs?.form) {
        return;
      }

      this.busy = true;

      this.$refs.form
        .validate()
        .then((form) => {
          if (form.valid) {
            this.user
              .update()
              .then((u) => {
                this.user = new User(u);
                this.$session.setUser(u);
                this.$notify.success(this.$gettext("Settings saved"));
              })
              .finally(() => {
                this.busy = false;
              });
          } else {
            this.$notify.error(this.$gettext("Changes could not be saved"));
            this.busy = false;
          }
        })
        .catch(() => {
          this.$notify.error(this.$gettext("Changes could not be saved"));
          this.busy = false;
        });
    },
    onUploadAvatar() {
      if (this.busy) {
        return;
      }

      this.busy = true;

      this.$notify.info(this.$gettext("Uploading…"));

      this.user
        .uploadAvatar(this.$refs.upload.files)
        .then((u) => {
          this.user = new User(u);
          this.$session.setUser(u);
          this.$notify.success(this.$gettext("Settings saved"));
        })
        .finally(() => (this.busy = false));
    },
  },
};
</script>
