<template>
  <div class="auth-footer">
    <div class="text-center footer-links text-caption">
      <router-link :to="{ name: 'about' }" class="text-white">{{ $gettext("About") }}</router-link>
    </div>
    <div class="text-center footer-links text-caption opacity-80">
      <span class="text-white">{{ appName }}</span>
    </div>
  </div>
</template>
<script>
export default {
  name: "PAuthFooter",
  data() {
    return {
      appName: this.$config.getName(),
    };
  },
};
</script>
