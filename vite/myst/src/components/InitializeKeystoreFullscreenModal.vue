<template>
  <div
    class="modal"
    v-if="show"
  >
    <div
      class="modal-overlay"
      v-on:click="hide"
    />
    <div class="modal-content">
      <form
        class="setup"
        @submit.prevent="submit"
      >
        <h4>First time setup</h4>
        <br />
        <div class="separator" />
        <br /><br /><br /><br /><br />
        <h6 v-if="step === 1">Let's choose a name for your first keystore. Keystores help organize your secrets
          into
          groups.
          <br /><br>
        </h6>
        <h6 v-else>Selected keystore name:<br></h6>
        <div
          :class="{invalid: !nameValid && warnings}"
          class="input"
        >
          <label>Keystore name</label>
          <input
            v-model="name"
            :class="{selected: step === 2}"
            :readonly="step !== 1"
            :tabindex="step === 1 ? 0 : -1"
            autofocus
            class="input keystore-name"
            placeholder="Keystore name"
            @focusout="warnings = true"
            @input="warnings = true"
          />
          <span>can only contain lowercase alphanumeric characters and hyphens</span>
        </div>

        <transition
          :duration="500"
          name="fade"
        >
          <div
            v-if="step === 2"
            class="step-2"
          >
            <h6>Your keystores will be encrypted using a master password. The security of your secrets will
              depend on
              its
              complexity. Choose one wisely and make sure you remember it.</h6>
            <br />
            <div
              :class="{invalid: !passwordValid && warnings}"
              class="input"
            >
              <label>Master password</label>
              <expanding-textarea
                ref="password"
                v-model="password"
                class="password-input"
                placeholder="Master password"
                @keydown.enter.prevent="submit"
              />
              <span>too weak</span>
            </div>
            <br />
          </div>
        </transition>

        <div class="footer">
          <span>Step {{ step }} of 2</span>
          <button
            v-if="step === 1"
            :class="{disabled: !valid}"
            class="step-button"
            type="submit"
          >Next
          </button>
          <button
            v-else
            :class="{disabled: !valid}"
            class="step-button green"
            type="submit"
          >Create</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import ExpandingTextarea from "./ExpandingTextarea.vue";

export default defineComponent({
  name: "InitializeKeystoreFullscreenModal",
  components: { ExpandingTextarea },
  props: ["show"],
  data: () => ({
    name: "",
    password: "",
    step: 1,
    nameRegex: /^[a-z0-9-]+$/,
    warnings: false,
  }),
  computed: {
    valid() {
      if (!this.nameValid) {
        return false;
      }

      if (this.step === 2 && !this.passwordValid) {
        return false;
      }

      return true;
    },
    nameValid() {
      if (this.name.length == 0) {
        return false;
      } else if (!this.nameRegex.test(this.name)) {
        return false;
      }

      return true;
    },
    passwordValid() {
      if (this.password.length == 0) {
        return false;
      } else if (this.password.length < 8) {
        return false;
      }

      return true;
    },
  },
  methods: {
    hide() {
      this.$emit("hide");
    },
    submit() {
      if (this.step === 1 && this.nameValid) {
        this.step = 2;
        this.warnings = true;
        this.$nextTick(() => {
          const textarea = this.$refs.password as ExpandingTextarea;
          textarea.$el.focus();
        });
      } else if (this.nameValid && this.passwordValid) {
        this.$store
          .dispatch("keystore/createKeystore", {
            name: this.name,
            password: this.password,
          })
          .then(() => {
            // TODO: do not use vuex or switch back
            // to vue 2 where everything works as expected
            // this.$emit('hide')

            this.$store.commit("setReady", true);
          });
      }
    },
  },
});
</script>

<style lang="scss" scoped>
$bg: #0a0e11;
$accent: #00edb1;
//$text-color: #f4f8fb;
$text-color: #fff;

@import url("https://rsms.me/inter/inter.css");

* {
  font-family: "Inter", sans-serif;
  font-weight: 300;
  font-size: 100%;
  color: $text-color;
  line-height: 1.325;
}

.step-2 {
  //opacity: 1;
  //transform: translateY(0);
  //transition: 1s ease;
}

//.show-enter-active,
//.show-leave-active {
//  opacity: 1;
//
//  .modal-content {
//    transform: scale(1);
//  }
//}
//
//.show-enter,
//.show-leave-to {
//  opacity: 0;
//  .modal-content {
//    transform: scale(0.95);
//  }
//}

.fade-enter-active,
.fade-leave-active {
  transition: 0.5s;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.fade-enter-to {
  opacity: 1;
  transform: translateY(0);
}

.fade-leave-from {
  opacity: 1;
  transform: translateY(0);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

textarea {
  margin: 0;
  padding: 0;
  border: none;
  outline: none;
  width: 100%;
  resize: none;
  //line-height: 1;
  padding: 15px 16px !important;
}

.separator {
  width: 100px;
  height: 2px;
  background-color: #31363d;
  margin: 0 auto;
}

//TODO focus

//box-shadow: 0 0 0 0 rgba($accent, .5);
//
//&:focus {
//	transition: .1s ease;
//	box-shadow: 0 0 0 2px rgba($accent, .5);
//}

input {
  margin: 0;
  border: 0;
  padding: 0;
  background-color: transparent;
  color: #fff;
  outline: none;
  width: 100%;
  height: 56px !important;
}

h4 {
  font-weight: 700;
  font-size: 1.8rem;
  margin: 0;
}

h6 {
  font-size: 1rem;
  margin: 0;
  //line-height: 1.4;
}

button {
  margin: 0;
  padding: 0;
  outline: none;
  border: none;
  background-color: transparent;
  height: 48px;

  //&:disabled {
  //	opacity: .5;
  //	background-color: ;
  //}
}

.step-button {
  //text-transform: uppercase;
  font-size: 1.1rem;
  font-weight: 500;
  padding: 0 40px;
  border-radius: 5px;
  margin: 0 5px;
  background-color: rgba(#202228, 1);
  color: #fff;

  &.disabled {
    background-color: #161819;
  }

  &.green {
    background-color: #002e23;
    color: $accent;

    &.disabled {
      background-color: #0c1d19;
    }
  }

  &:hover {
    //text-decoration: underline;
  }
}

.setup {
  position: fixed;
  color: #e1e4e8;
  width: 500px;

  .footer {
    margin-top: 100px;
    display: flex;
    flex-direction: column;
    align-items: center;

    span {
      font-size: 0.9rem;
      opacity: 0.95;
      margin-bottom: 12px;
    }
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

.input {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;

  label {
    flex-basis: 100%;
    margin-bottom: 10px;
    opacity: 0.8;
    padding: 4px 8px;
    display: none;
  }

  input {
    display: block;
    font-size: 1.1rem;
    font-weight: 400;
    //caret-color: $accent;
    box-sizing: border-box;
    background-color: rgba(#abc, 0.05);
    border-radius: 5px;
    //padding: 15px 16px;
    padding: 0 16px;
    color: #fff;
    //line-height: 1.5;
    overflow: hidden;

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

  &.invalid {
    span {
      opacity: 1;
    }
  }

  span {
    color: #ff9999;
    font-weight: 500;
    font-size: 0.9rem;
    padding: 10px 16px;
    opacity: 0;
  }
}

.keystore-name {
  &.selected {
    position: relative;
    background: transparent;
    font-weight: 600;
    color: $accent;
    padding: 0;
    font-size: 1.4rem;
    pointer-events: none;
    top: -16px;
    margin-bottom: 4px;
  }
}

.password-input {
  display: block;
  //border-bottom: 2px solid rgba($accent, .1);
  font-size: 1.1rem;
  font-weight: 400;
  //width: 100%;
  box-sizing: border-box;
  background-color: rgba(#abc, 0.05);
  border-radius: 5px;
  padding: 15px 16px;
  color: #fff;
  //line-height: 1.5;
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

//.show-enter-active, .show-leave-active {
//  transition: opacity 1s ease, transform 1s ease !important;
//}
//
//.show-enter {
//  opacity: 0;
//  //transform: translateY(-20px);
//}

//.show-leave-to {
//  opacity: 0;
//  //transform: translateY(-20px);
//}

//.show-enter-active, .show-leave-active {
//  transition: opacity .5s ease;
//}
//.show-enter, .show-leave-to /* .fade-leave-active below version 2.1.8 */ {
//  opacity: 0;
//}

//.show-enter-active,
//.show-leave-active {
//  transition: .5s ease;
//  opacity: 1;
//
//  transform: translateY(0);
//}

//.show-enter,
//.show-leave-to {
//  opacity: 0;
//  transition: .5s ease;
//  .modal-content {
//    transform: scale(0.95);
//  }
//}

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
    flex-direction: column;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
    pointer-events: all;
    transition: transform 0.1s ease;
    overflow: hidden;
    display: flex;
    justify-content: center;
    align-items: center;

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
