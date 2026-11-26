<template>
  <div class="p-pdf-viewer" />
</template>

<script>
export default {
  name: "PPdfViewer",
  props: {
    src: { type: [String, Object], default: "" },
    pages: { type: Number, default: 0 },
    hasPrev: { type: Boolean, default: false },
    hasNext: { type: Boolean, default: false },
    controlsVisible: { type: Boolean, default: true },
  },
  emits: ["media-prev", "media-next"],
};
</script>
