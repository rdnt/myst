<script lang="ts">
  import {fade} from 'svelte/transition';

  export let show: boolean;

  const keydown = (event: KeyboardEvent) => {
    if (event.key === 'Escape') {
      show = false;
    }
  }
</script>

{#if show}
  <div class="modal">
    <div class="overlay"></div>
    <div transition:fade={{delay: 0, duration: 180, ease: 'in'}} class="body">
      <div class="header">
        <slot name="header"></slot>
      </div>
      <div class="content">
        <slot></slot>
      </div>
      <div class="footer">
        <slot name="footer"></slot>
      </div>
    </div>
  </div>
{/if}

<svelte:window on:keydown={keydown}/>

<style lang="scss">
  .modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    height: -webkit-fill-available;
    padding: 50px 0;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    z-index: 1;

    .body {
      position: relative;
      max-height: 100%;
      max-width: 500px;
      flex-basis: 100%;
      background: #111519;
      display: flex;
      flex-direction: column;
      overflow: hidden;
      margin: 0 auto;
      padding: 20px;
      border-radius: 5px;
      box-sizing: border-box;
      box-shadow: 0 8px 24px 4px rgba(0, 0, 0, .5);
      will-change: transform;

      .content {
        overflow-y: auto;
      }
    }

    .overlay {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      height: -webkit-fill-available;
      background: #000;
      transition: opacity .18s ease;
      opacity: .85;
    }
  }
</style>
