<template>
  <div class="p-map" />
</template>

<script>
export default {
  name: "PMap",
  props: {
    latlng: { type: Array, default: () => [0.0, 0.0] },
    zoom: { type: Number, default: 9 },
    style: { type: String, default: "embedded" },
    interactive: { type: Boolean, default: false },
    draggable: { type: Boolean, default: false },
    showControls: { type: Boolean, default: false },
    clickable: { type: Boolean, default: false },
    markerClickable: { type: Boolean, default: false },
    animateDuration: { type: Number, default: -1 },
  },
  emits: ["update:latlng", "marker-moved", "marker-clicked", "map-clicked"],
};
</script>
