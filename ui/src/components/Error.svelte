<script lang="ts">
  import {slide} from 'svelte/transition';
  import {error} from "src/stores/messages"
</script>

{#if $error !== undefined}
  <div class="error" transition:slide>
    <div class="title">{$error.title}</div>
    <div class="message">{$error.message}</div>
  </div>
{/if}

<style lang="scss">
  .error {
    display: block;
    position: fixed;
    bottom: 0;
    right: 0;
    border-radius: 4px;
    margin: 20px;
    padding: 14px 20px;
    background-color: #2e2020;
    color: #ff9999;

    .title {
      color: #ff9999;
      font-weight: 600;
      font-size: 1.1rem;
      margin-bottom: 4px;
    }

    .message {
      color: #ff9999;
      font-weight: 400;
      font-size: 1rem;
    }
  }
</style>
