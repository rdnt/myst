<template>
  <textarea
    ref="textarea"
    v-model="value"
    :placeholder="placeholder"
    :style="{
      height: height + 'px',
    }"
    @input="
      $emit('change', value);
      update()
    "
  />
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "ExpandingTextarea",
  props: {
    placeholder: {
      type: String,
    },
  },
  data: (): {
    value: string;
    height: number;
  } => ({
    value: "",
    height: 0,
  }),
  mounted() {
    this.$nextTick(() => {
      this.value = this.$props.placeholder || "";

      this.update().then(() => {
        this.value = "";
      });

      window.addEventListener("resize", this.update);
    });
  },
  unmounted() {
    window.removeEventListener("resize", this.update);
  },
  methods: {
    update() {
      this.height = 0;

      return this.$nextTick(() => {
        const textarea = this.$refs.textarea as HTMLTextAreaElement;
        this.height = textarea.scrollHeight;
      });
    },
  },
});
</script>
