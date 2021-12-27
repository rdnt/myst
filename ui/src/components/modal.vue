<template>
  <div class="modal">
    <router-link tag="div" class="modal-overlay" :to="link" />
    <div class="modal-content">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: "modal",
  props: ["link"]
};
</script>

<style scoped lang="scss">
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
    background: #292e39;
    box-shadow: 0px 0px 47px rgba(#000, 0.5);
    position: absolute;
    display: flex;
    flex-direction: column;
    max-width: calc(100% - 40px);
    min-width: min(600px, 100%);
    max-height: calc(100% - 140px);
    top: 100px;
    backface-visibility: hidden;
    pointer-events: all;
    transition: transform 0.1s ease;
    border-radius: 5px;
    overflow: hidden;
  }
  .modal-header {
    flex-basis: 100px;
  }
  .modal-body {
    flex-grow: 1;
    overflow-y: auto;
  }
  .modal-footer {
    flex-basis: 100px;
  }
  > .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: #000;
    opacity: 0.66;
    pointer-events: all;
  }
}
</style>
