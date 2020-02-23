<template>
  <div id="login" :class="{ hide: login == true }">
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
export default {
  data() {
    return {
      login: true
    };
  },
  mounted() {
    setTimeout(() => {
      this.login = true;
    }, 2500);
  }
};
</script>
