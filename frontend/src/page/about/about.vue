<template>
  <div class="p-page p-page-about" tabindex="-1">
    <v-toolbar flat :density="$vuetify.display.smAndDown ? 'compact' : 'default'" color="secondary" class="page-toolbar p-page__navigation">
      <v-toolbar-title>{{ $config.getName() }}</v-toolbar-title>
    </v-toolbar>

    <div class="pa-6 text-selectable">
      <div class="pb-6">
        <div class="d-flex flex-no-wrap justify-space-between align-center">
          <div>
            <p class="text-h5 font-weight-medium mb-2">{{ $config.getName() }}</p>
            <p class="text-body-2 text-secondary-dark mb-0">
              {{ $config.values.siteCaption || "Self-Hosted Photos App" }} · {{ versionLabel }}
            </p>
          </div>
          <img :src="$config.getIcon()" :alt="$config.getName()" class="ma-2" width="64" height="64" />
        </div>
      </div>

      <v-divider class="my-4" />

      <div class="text-body-1">
        <p class="mb-3">
          MyPhotoX is a lightweight, self-hosted personal photo and video library designed for fast local browsing, album management, and public link sharing.
        </p>
        <p class="mb-2">Core features included in this edition:</p>
        <ul class="ml-6 mb-4">
          <li>Photo &amp; video library indexing, upload, and import</li>
          <li>Timeline browsing, search filters, folders, favorites, and archive</li>
          <li>Albums, batch editing, metadata editing, and ZIP downloads</li>
          <li>Anonymous public share links</li>
          <li>English &amp; Vietnamese interface</li>
        </ul>
      </div>

      <p class="text-caption mt-6 mb-0 text-secondary-dark">
        Licensed under GNU Affero General Public License v3.0 (AGPL-3.0).
      </p>
    </div>
    <p-about-footer></p-about-footer>
  </div>
</template>

<script>
import PAboutFooter from "component/about/footer.vue";

export default {
  name: "PPageAbout",
  components: {
    PAboutFooter,
  },
  computed: {
    versionLabel() {
      const v = this.$config.values.version || "";
      return v ? `v${v}` : "v1.0";
    },
  },
  mounted() {
    this.$view.enter(this);
  },
  unmounted() {
    this.$view.leave(this);
  },
};
</script>
