<template>
  <transition name="show">
    <div class="modal" v-if="show">
      <div class="modal-overlay" v-on:click="hide" />
      <div class="modal-content">
        <div class="modal-header">
          <slot name="header"></slot>
        </div>
        <div class="modal-body">
          <slot></slot>
        </div>
        <div class="modal-footer">
          <button class="button" v-on:click="hide">
            Cancel
          </button>
          <slot name="footer"></slot>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
export default {
  name: "modal",
  props: ["show"],
  methods: {
    hide() {
      this.$emit("hide");
    }
  }
};
</script>

<style scoped lang="scss">
$bg: #1f2023;

.show-enter-active,
.show-leave-active {
  opacity: 1;

  .modal-content {
    transform: scale(1);
  }
}

.show-enter,
.show-leave-to {
  opacity: 0;
  .modal-content {
    transform: scale(0.95);
  }
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  pointer-events: none;
  transition: opacity 0.1s ease;

  .modal-content {
    background: $bg;
    box-shadow: 0 0 40px rgba(darken($bg, 10%), 0.25);
    position: absolute;
    display: flex;
    flex-direction: column;
    max-width: calc(100% - 40px);
    width: min(600px, 100%);
    max-height: calc(100% - 140px);
    top: 100px;
    pointer-events: all;
    transition: transform 0.1s ease;
    border-radius: 5px;
    overflow: hidden;

    .modal-header {
      flex-shrink: 0;
      border-bottom: 1px solid rgba(48, 50, 54, 0.5);
      padding: 24px;
      //padding: 16px;
      box-sizing: border-box;
      font-size: 1.2rem;
      color: rgba(215, 216, 219, 1);
    }

    .modal-body {
      overflow-y: auto;
      flex-basis: 100%;
      min-height: 60px;
      padding: 16px;
      box-sizing: border-box;
      color: rgb(247, 248, 248);
      color: #63666d;
    }

    .modal-footer {
      display: flex;
      flex-direction: row;
      justify-content: flex-end;
      flex-shrink: 0;
      border-top: 1px solid rgba(48, 50, 54, 0.5);
      //padding: 0 20px 20px 20px;
      padding: 18px 16px;
      box-sizing: border-box;
      color: rgb(215, 216, 219);
    }
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: darken($bg, 10%);
    opacity: 0.66;
    pointer-events: all;
  }
}
</style>
