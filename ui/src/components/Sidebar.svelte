<script lang="ts">
  import Link from "@/components/Link.svelte";
  import {myInvitations} from "@/stores/invitations";
  import {onMount} from "svelte";

  export let keystores;

  onMount(() => {


  })
</script>

<div class="sidebar">
  <h4>Myst</h4>
  <div class="list">
    <h5>Keystores</h5>
    {#each keystores as keystore}
      <Link path="/keystore/{keystore.id}">{keystore.name}</Link>
    {/each}
  </div>

  <div class="list bottom">
    <h5>More</h5>
    <div class="rel">
      <Link path="/sharing">
        Sharing
        {#if $myInvitations.length > 0}
          <div class="badge">{$myInvitations.length}</div>
        {/if}
      </Link>

<!--      <Link active={showInvitations} on:click={() => showInvitations = !showInvitations}>-->
<!--        Invitations-->
<!--        {#if $invitations.length > 0 || true}-->
<!--          <div class="badge">{$invitations.length + 2}</div>-->
<!--        {/if}-->
<!--      </Link>-->
    </div>
    <h5 class="version">v0.0.0-0123456</h5>
  </div>

</div>

<style lang="scss">
  .sidebar {
    position: relative;
    background-color: #0a0e11;
    height: 100%;
    padding: 20px;
    box-sizing: border-box;
    flex-basis: 300px;
    display: flex;
    flex-direction: column;

    h4 {
      font-weight: 700;
      font-size: 2rem;
      padding: 0 20px;
      margin: 0;
      margin-bottom: 40px;
    }

    h5 {
      height: 20px;
      padding: 0 20px;
      margin: 0 0 10px;
      color: #8a8f9f;
      text-transform: uppercase;
      font-size: .85rem;
      font-weight: 600;
      letter-spacing: .5px;
    }

    .rel {
      position: relative;
      width: 100%;
    }

    .version {
      margin-top: 40px;
      opacity: .75;
      font-weight: 500;
      text-transform: none;
    }

    .list {
      display: flex;
      flex-direction: column;
      flex-wrap: wrap;
      margin: 0;
      margin-top: 20px;

      &.bottom {
        bottom: 0;
        margin-top: auto;
      }

      :global(.link a), :global(.link button) {
        margin: 0;
        line-height: 1.4;
        background-color: transparent;
        border: none;
        color: inherit;
        text-decoration: none;
        display: flex;
        width: 100%;
        align-items: center;
        border-radius: 5px;
        position: relative;
        cursor: pointer;
        padding: 10px 20px 10px 20px;
        font-size: 1.1rem;
        white-space: nowrap;
        text-overflow: ellipsis;
        font-weight: 500;
        box-sizing: border-box;

        &:hover {
          background-color: rgba(#1e2328, .75);
        }
      }

      :global(.link.active a), :global(.link.active button) {
        font-weight: 500;
        color: #00edb1;
        background-color: #0c1d19;
        background-color: rgba(#1e2328, .75);

        &:after {
          content: '';
          position: absolute;
          top: 0;
          width: 10px;
          height: 100%;
          left: -20px;
          border-top-right-radius: 5px;
          border-bottom-right-radius: 5px;
        }
      }
    }
  }

  .badge {
    margin-left: auto;
    font-weight: bold;
  }
</style>
