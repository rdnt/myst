<script lang="ts">
  import {messages, MessageType} from "@/stores/messages";
  import {slide} from 'svelte/transition';
</script>

{#if $messages.length > 0}
  <div class="messages">
    {#each $messages as message}
      <div class="message"
           class:info={message.type === MessageType.Info}
           class:error={message.type === MessageType.Error}
           transition:slide>
        <div class="title">{message.title}</div>
        {#if message.message}
          <div class="msg">{message.message}</div>
        {/if}
      </div>
    {/each}
  </div>
{/if}


<style lang="scss">
  .messages {
    position: fixed;
    bottom: 0;
    right: 0;
    margin: 20px;
  }

  .message {
    margin-top: 10px;
    border-radius: 4px;
    padding: 14px 20px;

    &.info {
      background-color: #002e23;
      color: #00edb1;
    }

    &.error {
      background-color: #2e2020;
      color: #ff9999;
    }

    .title {
      font-weight: 600;
      font-size: 1.1rem;
    }

    .msg {
      margin-top: 4px;
      font-weight: 400;
      font-size: 1rem;
    }
  }
</style>
