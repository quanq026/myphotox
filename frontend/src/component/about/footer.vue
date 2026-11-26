<template>
  <div id="p-about-footer">
    <div v-if="legalInfo" class="row legal text-center text-caption py-2 text-secondary-dark">
      {{ legalInfo }}
    </div>
    <v-footer class="bg-background pa-2 text-caption text-secondary-dark">
      <div class="footer-about text-truncate">
        <router-link :to="{ name: 'about' }" class="body-link text-truncate">{{ appName }} · {{ buildVersion }}</router-link>
      </div>
      <div class="footer-copyright text-truncate">
        <span>{{ appName }}</span>
      </div>
    </v-footer>
  </div>
</template>
<script>
export default {
  name: "PAboutFooter",
  data() {
    const version = this.$config.values.version || "";
    const shortVersion = version.split("-")[0];
    return {
      appName: this.$config.getName(),
      buildVersion: shortVersion ? `v${shortVersion}` : version,
      legalInfo: this.$config.get("legalInfo"),
    };
  },
};
</script>
