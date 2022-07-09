<script lang="ts">
  import type {Invitation} from "@/api";

  export let invitation: Invitation;
</script>

<div class="invitation" class:empty={!invitation}>
  {#if invitation}

  {/if}
</div>

<style lang="scss">
  $accent: #00edb1;

  h5 {
    font-weight: 600;
    font-size: 1.8rem;
    margin: 0;
  }

  .invitation {
    display: flex;
    flex-direction: column;
    background-color: #101519;
    border-left: 2px solid #1a2025;
    height: 100%;
    padding: 20px;
    box-sizing: border-box;
    flex-basis: 50%;
    overflow-y: auto;
    padding: 32px;
    //flex-grow: 1;
    //padding-top: 100px;

    &.empty {
      &:after {
        content: "";
        background-color: #1b2025;
        border-radius: 50%;
        font-size: 1.5em;
        font-weight: bold;
        text-align: center;
        display: block;
        width: 100px;
        height: 100px;
        line-height: 100%;
      }
    }
  }

  .button {
    outline: none;
    border: none;
    height: 40px;
    font-size: 1.1rem;
    font-weight: 500;
    padding: 0 16px;
    border-radius: 5px;
    background-color: rgba(#202228, 1);
    color: #fff;
    margin-left: 10px;

    &.left {
      //margin-right: auto;
    }

    &.disabled {
      //background-color: #161819;
      opacity: .5;
    }

    &.green {
      background-color: rgba(#002e23, .9);
      color: $accent;

      &.disabled {
        background-color: #0c1d19;
      }
    }

    &.transparent {
      background-color: transparent;
      padding: 0 12px;

      &.disabled {

      }
    }

    &.red {
      background-color: #2e2020;
      color: #ff9999;

      &.disabled {
        background-color: rgba(29, 29, 12, 0.99);
      }
    }
  }

  .modal-header {
    //margin-bottom: 22px;
  }

  .modal-footer {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    margin-top: 22px;
  }

  .delete-title {
    padding: 4px;
    box-sizing: border-box;
    font-size: 1.1rem;
  }

  .fields {
    margin-bottom: 22px;
  }
</style>
