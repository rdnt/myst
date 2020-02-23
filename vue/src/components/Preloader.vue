<template>
  <div
    id="preloader"
    :class="{
      pulse: pulse == true,
      hide: hide == true,
      loaded: loaded == true
    }"
  >
    <img class="black" src="/assets/images/vault.svg" />
    <img class="logo" src="/assets/images/vault.svg" />
    <div class="glow">
      <img src="/assets/images/vault.svg" />
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      pulse: false,
      hide: false,
      loaded: false
    };
  },
  mounted() {
    this.pulse = true;
    setTimeout(() => {
      this.hide = true;
    }, 1000);
    setTimeout(() => {
      this.loaded = true;
    }, 1500);
  }
};
</script>

$("#preloader #black").addClass("pulse"); setTimeout(function() { $("#preloader
#logo").addClass("hide"); $("#preloader #glow").addClass("hide");
setTimeout(function() { $("#preloader #logo").addClass("loaded");
$("#preloader").addClass("loaded"); setTimeout(function() {
$(".section").addClass("loaded"); }, 600); setTimeout(function() { $("#preloader
#logo, #preloader #glow, #preloader #black").addClass("invis"); $("html,
body").removeClass("unscroll"); }, 1000); }, 500); }, 1000);
