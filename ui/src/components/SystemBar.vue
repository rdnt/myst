<template>
  <div id="system-bar">
    <div class="title">
      <!-- <img src="images/bar/vaulttext.svg"> -->
    </div>
    <button class="option close" @click="close()">
      <img src="/assets/images/bar/close.svg" alt="" />
    </button>
    <button class="option maximize" @click="maximize()">
      <img src="/assets/images/bar/maximize.svg" alt="" />
    </button>
    <button class="option minimize" @click="minimize()">
      <img src="/assets/images/bar/minimize.svg" alt="" />
    </button>
  </div>
</template>

<script>
export default {
  data() {
    return {};
  },
  methods: {
    close: function() {
      window.close();
    },
    maximize: function() {
      if (this.client == "electron") {
        if (!window.isMaximized()) {
          window.maximize();
        } else {
          window.unmaximize();
        }
      }
    },
    minimize: function() {
      if (this.client == "electron") {
        window.minimize();
      }
    }
  }
};
</script>
