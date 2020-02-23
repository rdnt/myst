<template>
  <div id="generator" class="container">
    <div class="field input result">
      <input
        type="text"
        value="!#^TBVtbn89!#TB1#TB71tN!T^3"
        @change="setStrength"
        @input="setStrength"
      />
    </div>
    <div
      class="strength-indicator "
      :class="{
        s1: strength == 1,
        s2: strength == 2,
        s3: strength == 3,
        s4: strength == 4
      }"
    >
      <div class="bars">
        <div class="bar"></div>
        <div class="bar"></div>
        <div class="bar"></div>
        <div class="bar"></div>
      </div>
      <label>{{ strengthLabel }}</label>
    </div>
    <div class="settings-container" :class="{ expanded: settingsVisible }">
      <div class="settings-toggle">
        <div>
          Settings
        </div>
        <!-- <div class="field toggle">
            <input name="toggle" type="checkbox" @click="toggleSettings" />
            <template v-if="settingsVisible">
              <label>Hide settings</label>
            </template>
            <template v-else>
              <label>Show settings</label>
            </template>
          </div> -->
      </div>
      <div class="settings" :class="{ show: settingsVisible }">
        <div class="column">
          <div class="field checkbox">
            <input name="test" type="checkbox" checked />
            <label>Numbers (123 etc.)</label>
          </div>
          <div class="field checkbox">
            <input name="test" type="checkbox" checked />
            <label>Uppercase characters (ABC etc.)</label>
          </div>
          <div class="field checkbox">
            <input name="test" type="checkbox" checked />
            <label>Lowercase characters (abc etc.)</label>
          </div>
          <div class="field checkbox">
            <input name="test" type="checkbox" checked />
            <label>Symbols (!@# etc.)</label>
          </div>
        </div>
        <div class="column">
          <div class="field checkbox">
            <input name="test" type="checkbox" />
            <label>Avoid similar characters (e.g. I and l)</label>
          </div>
          <div class="field checkbox">
            <input name="test" type="checkbox" />
            <label>Show phonetic words (e.g. Sierra Hotel Tango)</label>
          </div>
          <div class="field checkbox">
            <input name="test" type="checkbox" />
            <label>Generate multiple passwords</label>
          </div>
          <div class="field checkbox">
            <input name="test" type="checkbox" />
            <label>Save settings for later</label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
// import Tabs from "@/components/Tabs.vue";

export default {
  data() {
    return {
      settingsVisible: true,
      strength: 0,
      strengthLabel: ""
    };
  },
  mounted() {},
  methods: {
    toggleSettings: function() {
      this.settingsVisible = !this.settingsVisible;
    },
    setStrength: function(event) {
      var password = event.target.value;
      if (password.length < 8) {
        this.strength = 1;
        this.strengthLabel = "Weak";
      } else if (password.length < 12) {
        this.strength = 2;
        this.strengthLabel = "Good";
      } else if (password.length < 16) {
        this.strength = 3;
        this.strengthLabel = "Strong";
      } else {
        this.strength = 4;
        this.strengthLabel = "Insane";
      }
    }
  }
};
</script>
