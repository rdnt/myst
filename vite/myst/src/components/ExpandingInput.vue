<template>
  <input
    ref="input"
    v-model="value"
    :placeholder="placeholder"
    :style="{
      width: width + 'px',
      minWidth: minWidth + 'px',
    }"
    @input="() => {
      $emit('change', value);
      update()
    }"
  />
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "ExpandingInput",
  props: {
    placeholder: {
      type: String,
    },
  },
  data: (): {
    value: string;
    width: number;
    minWidth: number;
  } => ({
    value: "default",
    width: 0,
    minWidth: 0,
  }),
  mounted() {
    this.$nextTick(() => {
      const input = this.$refs.input as HTMLInputElement;
      this.value = this.$props.placeholder || "";

      this.update().then(() => {
        this.minWidth = input.scrollWidth;
        this.value = "";
      });

      window.addEventListener("resize", this.update);
    });
  },
  methods: {
    update() {
      const input = this.$refs.input as HTMLInputElement;
      this.width = 0;

      return this.$nextTick(() => {
        this.width = input.scrollWidth;
      });
    },
  },
});
</script>
