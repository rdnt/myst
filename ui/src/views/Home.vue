<template>
  <div class="home">
    <InitializeKeystoreFullscreenModal :show="true"/>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import InitializeKeystoreFullscreenModal from '@/components/InitializeKeystoreFullscreenModal.vue'

export default defineComponent({
  name: 'Home',
  components: {
    InitializeKeystoreFullscreenModal
  }
})
</script>
