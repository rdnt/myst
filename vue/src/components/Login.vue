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
      <div class="label">Master Password</div>
      <div class="master-password">
        <input
          class="field placeholder"
          type="text"
          autocomplete="off"
          autocorrect="off"
          autocapitalize="off"
          spellcheck="false"
          :value="placeholder"
          tabindex="-1"
        />
        <input
          class="field"
          type="text"
          autocomplete="off"
          autocorrect="off"
          autocapitalize="off"
          spellcheck="false"
          autofocus
          v-model="password"
          @keydown.enter="setEnter(true)"
          @keyup="setEnter(false)"
          @focusout="setEnter(false)"
        />
      </div>
      <div class="prompt" :class="{ 'enter-pressed': enter }">
        Press <span class="code">enter</span> to submit
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
      placeholder: "",
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

        this.$store.commit("login");
      }
    }
  }
};
</script>
