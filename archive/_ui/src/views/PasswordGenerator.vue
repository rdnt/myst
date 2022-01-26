<template>
  <div id="generator" class="container">
    <div class="field input result">
      <input
        type="text"
        autocomplete="off"
        autocorrect="off"
        autocapitalize="off"
        spellcheck="false"
        v-model="password"
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
      <label>{{ labels[strength] }}</label>
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
        <div class="label">
          Options
        </div>
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
  name: "PasswordGenerator",
  data() {
    return {
      password: "",
      settingsVisible: true,
      strength: 1,
      // strengthLabel: "Very Weak",
      labels: {
        1: "Very weak",
        2: "Weak",
        3: "Good",
        4: "Strong"
      }
    };
  },
  mounted() {},
  methods: {
    toggleSettings: function() {
      this.settingsVisible = !this.settingsVisible;
    },
    setStrength: function() {
      if (this.password.length < 8) {
        this.strength = 1;
      } else if (this.password.length < 12) {
        this.strength = 2;
      } else if (this.password.length < 16) {
        this.strength = 3;
      } else if (this.password.length < 32) {
        this.strength = 4;
      } else {
        this.strength = 5;
      }
    }
  }
};
</script>
