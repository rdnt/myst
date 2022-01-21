<template>
  <transition name="show">
    <div class="modal" v-if="show">
      <div class="modal-overlay" v-on:click="hide" />
      <div class="modal-content">
        <div class="setup">
          <h4>First time setup</h4>
          <br/>
          <div class="separator"/>
          <br/><br/><br/><br/>
          <h6 v-if="step === 1">Let's choose a name for your first keystore. Keystores help organize your secrets into groups.</h6>
          <h6 v-if="step === 2">Your selected keystore name:</h6>
          <br v-if="step === 1" />
          <div class="input">
            <label>Keystore name</label>
            <input class="password-input" :class="step === 2 ? 'selected': ''" placeholder="Keystore name" />
          </div>

          <template v-if="step === 1">
            <div class="footer">
              <button class="step-button" @click="step = 2">Next</button>
<!--              <button class="step-button green">Create</button>-->
            </div>
          </template>

          <template v-else>
            <br/><br/><br/><br/><br/>
            <h6>Your keystores will be encrypted using a master password. The security of your secrets will depend on its
              complexity. Choose wisely.</h6>
            <br/>
            <div class="input">
              <label>Master password</label>
              <expanding-textarea class="password-input" placeholder="Master password" @change="masterPasswordChanged" />
            </div>
            <br/>
            <div class="footer">
              <!--<button class="step-button">Cancel</button>-->
              <button class="step-button green">Create</button>
            </div>
          </template>

        </div>
      </div>
    </div>
  </transition>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import ExpandingTextarea from '@/components/ExpandingTextarea.vue'

export default defineComponent({
  name: 'InitializeKeystoreFullscreenModal',
  components: { ExpandingTextarea },
  props: ['show'],
  data: () => ({
    keystoreName: '',
    masterPassword: '',
    step: 1
  }),
  methods: {
    hide () {
      this.$emit('hide')
    },
    keystoreNameChanged (name: string) {
      this.keystoreName = name
    },
    masterPasswordChanged (password: string) {
      this.masterPassword = password
    }
  }
})
</script>

<style scoped lang="scss">
$bg: #0a0e11;
$accent: #00edb1;
//$text-color: #f4f8fb;
$text-color: #fff;

@import url('https://rsms.me/inter/inter.css');
* {
  font-family: 'Inter', sans-serif;
  font-weight: 300;
  font-size: 100%;
  color: $text-color;
  line-height: 1;
}

textarea {
  margin: 0;
  padding: 0;
  border: none;
  outline: none;
  width: 100%;
  resize: none;
  line-height: 1;
  padding: 15px 16px !important;
}

.input {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;

  label {
    flex-basis: 100%;
    margin-bottom: 10px;
    opacity: .8;
    padding: 4px 8px;
    display: none;
  }
}

.separator {
  width: 100px;
  height: 2px;
  background-color: #31363d;
  margin: 0 auto;
}

input {
  margin: 0;
  border: 0;
  padding: 0;
  background-color: transparent;
  color: #fff;
  outline: none;
  width: 100%;
}

h4 {
  font-weight: 700;
  font-size: 1.8rem;
  margin: 0;
}

h6 {
  font-size: 1rem;
  margin: 0;
  line-height: 1.4;
}

button {
  margin: 0;
  padding: 0;
  outline: none;
  border: none;
  background-color: transparent;
}

.step-button {
  font-size: 1.1rem;
  //text-transform: uppercase;
  font-weight: 500;
  padding: 14px 30px;
  border-radius: 5px;
  margin: 0 5px;
  background-color: rgba(#fff, .05);
  color: #fff;
  //text-decoration: underline;
  //text-decoration-color: rgba(#fff, .25);
  //text-underline-offset: 2px;

  &.green {
    background-color: #0d2322;
    color: #56e9b6;
  }

  &:hover {
    //text-decoration: underline;
  }
}

.setup {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #e1e4e8;
  width: 480px;

  .header {
    position: relative;
    margin-bottom: 60px;
    width: 500px;
    box-sizing: border-box;
  }

  .footer {
    margin-top: 60px;
    display: flex;
    justify-content: center;

  }

  h4 {
    font-size: 1.8rem;
    margin-bottom: 12px;
    text-align: center;
  }

  h6 {
    font-size: 1.1rem;
    margin-bottom: 6px;
    font-weight: 300;
  }
}

.password-input {
  display: block;
  //border-bottom: 2px solid rgba($accent, .1);
  font-size: 1.1rem;
  font-weight: 400;
  //width: 100%;
  caret-color: $accent;
  box-sizing: border-box;
  background-color: rgba(#abc, .05);
  border-radius: 5px;
  padding: 15px 16px;
  color: #fff;
  line-height: 1.5;
  overflow: hidden;

  &.selected {
    background: transparent;
    font-weight: 600;
    color: $accent;
    padding: 0;
    height: auto;
    font-size: 1.4rem;
  }

  &::placeholder {
    color: lighten(#68737e, 5%);
  }

  &:focus {
    //border-bottom: 2px solid lighten(#3b4048, 10%);

    &::placeholder {
      color: lighten(#68737e, 15%);
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
