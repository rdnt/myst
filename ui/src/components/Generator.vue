<template>
  <div id="generator">
    <div class="container card">
      <div class="result">
        !#^TBVtbn89!#TB1#TB71tN!T^3
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
            <div>
              <div class="field checkbox">
                <input name="test" type="checkbox" checked />
                <label>Numbers (123 etc.)</label>
              </div>
            </div>
            <div>
              <div class="field checkbox">
                <input name="test" type="checkbox" checked />
                <label>Uppercase characters (ABC etc.)</label>
              </div>
            </div>
            <div>
              <div class="field checkbox">
                <input name="test" type="checkbox" checked />
                <label>Lowercase characters (abc etc.)</label>
              </div>
            </div>
            <div>
              <div class="field checkbox">
                <input name="test" type="checkbox" checked />
                <label>Symbols (!@# etc.)</label>
              </div>
            </div>
          </div>
          <div class="column">
            <div>
              <div class="field checkbox">
                <input name="test" type="checkbox" />
                <label>Avoid similar characters (e.g. I and l)</label>
              </div>
            </div>
            <div>
              <div class="field checkbox">
                <input name="test" type="checkbox" />
                <label>Show phonetic words (e.g. Sierra Hotel Tango)</label>
              </div>
            </div>
            <div>
              <div class="field checkbox">
                <input name="test" type="checkbox" />
                <label>Generate multiple passwords</label>
              </div>
            </div>
            <div>
              <div class="field checkbox">
                <input name="test" type="checkbox" />
                <label>Save settings for later</label>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      settingsVisible: true
    };
  },
  mounted() {},
  methods: {
    toggleSettings: function() {
      this.settingsVisible = !this.settingsVisible;
    }
  }
};
</script>
