<template>
  <div id="login" :class="{ show: !loggedIn }">
    <div id="loader"></div>
    <div class="form">
      <img
        class="logo"
        :class="{ submitting: loggingIn }"
        src="/assets/images/vault.svg"
        alt=""
      />
      <div class="label">Myst</div>
      <div class="master-password">
        <input
          v-model="passphrase"
          :placeholder="'keystore passphrase'"
          autocapitalize="off"
          autocomplete="off"
          autofocus
          class="field"
          spellcheck="false"
          type="password"
          @focusout="setEnter(false)"
          @keyup="setEnter(false)"
          @keydown.enter="setEnter(true)"
        />
        <div class="prompt" :class="{ 'enter-pressed': enter }">
          <span>↵ enter</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  data() {
    return {
      passphrase: "",
      enter: false,
      loggingIn: false
    };
  },
  computed: mapState({
    loggedIn: state => state.keystore.keystore !== null
  }),
  mounted() {},
  watch: {
    passphrase: function() {
      this.placeholder = "•".repeat(this.passphrase.length);
    }
  },
  methods: {
    setEnter(pressed) {
      this.enter = pressed;
      if (pressed && !this.loggingIn) {
        // calculate hash and login
        this.loggingIn = true;

        this.$store
          .dispatch("keystore/authenticate", {
            keystoreId: "0000000000000000000000",
            passphrase: this.passphrase
          })
          .then(() => {
            this.$router.push(
              "/keystore/" + this.$store.state.keystore.keystore.id
            );
          })
          .finally(() => {
            this.loggingIn = false;
          });
      }
    }
  }
};
</script>
