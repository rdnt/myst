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
          class="field"
          type="password"
          autocomplete="off"
          autocapitalize="off"
          spellcheck="false"
          autofocus
          v-model="password"
          :placeholder="'master password'"
          @keydown.enter="setEnter(true)"
          @keyup="setEnter(false)"
          @focusout="setEnter(false)"
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
      password: "",
      enter: false,
      loggingIn: false
    };
  },
  computed: mapState({
    loggedIn: state => state.loggedIn
  }),
  mounted() {},
  watch: {
    password: function() {
      this.placeholder = "•".repeat(this.password.length);
    }
  },
  methods: {
    setPassword(event) {
      this.password = event.target.value;
    },
    setEnter(pressed) {
      this.enter = pressed;
      if (pressed) {
        // calculate hash and login
        // this.loggingIn = true;
        // setTimeout(() => {
        //   this.loggingIn = false;
        // }, 1000);

        setTimeout(() => {
          this.$store.commit("login");
        }, 1000);
      }
    }
  }
};
</script>
