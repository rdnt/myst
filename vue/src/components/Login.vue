<template>
  <div id="login" :class="{ show: !loggedIn }">
    <form>
      <img class="logo" src="/assets/images/vault.svg" alt="" />
      <input
        class="field"
        name="PASSWORD"
        placeholder="PASSWORD"
        type="password"
        autocomplete="off"
        autofocus
      />
    </form>
  </div>
</template>

<script>
import { mapState } from "vuex";
export default {
  data() {
    return {};
  },
  computed: mapState({
    loggedIn: state => state.loggedIn
  }),
  mounted() {
    setTimeout(() => {
      this.login = true;
    }, 2500);
  }
};
</script>
