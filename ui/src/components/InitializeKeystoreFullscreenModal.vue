<template>
  <transition name="show">
    <div class="modal" v-if="show">
      <div class="modal-overlay" v-on:click="hide" />
      <div class="modal-content">
        <div class="setup">
          <div class="header">
            <h4>First time setup</h4>
            <h6>It seems that this is the first time you are opening Myst on this device.</h6>
            <h6>Let's choose a master password. It will be used to securely store your secrets.</h6>
          </div>
          <label>
            <input spellcheck="false" class="password-input" placeholder="master password"/>
          </label>
          <div class="footer">
            <button class="step-button">Next</button>
          </div>

        </div>

        <!--        <div class="modal-header">-->
<!--          <slot name="header"></slot>-->
<!--        </div>-->
<!--        <div class="modal-body">-->
<!--          <slot></slot>-->
<!--        </div>-->
<!--        <div class="modal-footer">-->
<!--          <button class="button" v-on:click="hide">-->
<!--            Cancel-->
<!--          </button>-->
<!--          <slot name="footer"></slot>-->
<!--        </div>-->
      </div>
    </div>
  </transition>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'InitializeKeystoreFullscreenModal',
  props: ['show'],
  methods: {
    hide () {
      this.$emit('hide')
    }
  }
})
</script>

<style scoped lang="scss">
$bg: #24292e;

@import url('https://rsms.me/inter/inter.css');
* {
  font-family: 'Inter', sans-serif;
  font-weight: 300;
  font-size: 100%;
}

input {
  margin: 0;
  border: 0;
  padding: 0;
  background-color: transparent;
  color: #fff;
  outline: none;
}

h4 {
  font-size: 1.8rem;
  margin: 0;
}

h6 {
  font-size: 1rem;
  margin: 0;
}

button {
  margin: 0;
  padding: 0;
  outline: none;
  border: none;
  background-color: transparent;
}

.step-button {
  color: #fff;
  font-size: 1rem;
  font-weight: 400;
  padding: 8px 0;

  &:hover {
    text-decoration: underline;
  }
}

.setup {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #e1e4e8;

  .header {
    margin-bottom: 60px;
  }

  .footer {
    margin-top: 60px;
    display: flex;
    justify-content: flex-end;
  }

  h4 {
    font-size: 1.8rem;
    margin-bottom: 30px;
  }

  h6 {
    font-size: 1rem;
    margin-bottom: 6px;
  }
}

.password-input {
  border-bottom: 2px solid #3b4048;
  font-size: 2rem;
  padding: 10px 0;
  width: 100%;
  caret-color: #9ecbff;

  &::placeholder {
    color: #68737e;
  }

  &:focus {
    border-bottom: 2px solid lighten(#3b4048, 10%);

    &::placeholder {
      color: lighten(#68737e, 10%);
    }
  }
}

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
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
    pointer-events: all;
    transition: transform 0.1s ease;
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
