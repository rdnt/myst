<script lang="ts">
  import {useFocus, useParams} from "svelte-navigator";
  import Link from "../components/Link.svelte";
  import Entry from "../components/Entry.svelte";
  import EntryPlaceholder from "../components/EntryPlaceholder.svelte";
  import CreateInvitationModal from "../components/CreateInvitationModal.svelte";

  export let keystore;

  let showCreateInvitationModal: boolean = false;

  const params = useParams();
  const registerFocus = useFocus();

  let entry;

  $: entry = keystore?.entries.find(e => e.id === $params.entryId);

  function createInvitation(detail: any) {
    console.log("createInvitation", detail);
  }
</script>

<div class="entries-list">
  <div class="entries" use:registerFocus>
    <div class="entries-list-header">
      <button on:click={() => showCreateInvitationModal = true} class="button"><span class="icon"></span>Share</button>
    </div>
    {#each keystore.entries as entry}
      <Link path={`/keystore/${keystore.id}/entry/${entry.id}`}>
        <div class="entry">
          <span class="icon">
            <img alt="" src={`https://${entry.website}/favicon.ico`}/>
          </span>
          <div class="info">
            <span class="name">
              {entry.website}
            </span>
            <span class="user">
              {entry.username}
            </span>
          </div>
        </div>
      </Link>
    {/each}
  </div>
</div>
{#if entry}
  <Entry {entry} {keystore}/>
{:else}
  <EntryPlaceholder />
{/if}

<CreateInvitationModal bind:show={showCreateInvitationModal} {keystore} on:submit={(e) => {createInvitation(e.detail)}}/>

<style lang="scss">
  .entries-list {
    position: relative;
    background-color: #101519;
    height: 100%;
    flex-basis: 50%;

    .entries {
      overflow-y: auto;
      height: calc(100% - 0px);
      padding: 20px;
      box-sizing: border-box;

      :global(.link a) {
        color: #f00;
        text-decoration: none;
      }

      .entry {
        position: relative;
        display: flex;
        flex-direction: row;
        flex-wrap: nowrap;
        align-items: center;
        padding: 10px 14px;
        box-sizing: border-box;
        border-radius: 5px;
        min-height: 20px;
        text-decoration: none;
        margin-bottom: 2px;

        .info {
          display: flex;
          flex-direction: column;
        }

        &.header {
          color: rgb(138, 143, 152);
          height: 60px;
          padding: 0 34px;
          margin-top: 20px;

          .name {
            flex-basis: calc(30% + 34px);
          }

          span {
            color: rgba(#8a8f9f, .75);
            text-transform: uppercase;
            font-size: 0.85rem;
            font-weight: 600;
            letter-spacing: 0.5px;
            pointer-events: none;
          }

          button {
            background-color: transparent;
            border: none;
            padding: 0;
            margin: 0;

            img {
              width: 10px;
              height: 10px;
              margin-right: 10px;
              opacity: .5;
            }
          }
        }

        .icon {
          padding-right: 16px;

          img {
            display: block;
            width: 32px;
            height: 32px;
            vertical-align: baseline;
          }
        }

        .name {
          flex-basis: calc(30%);
          font-weight: 500;
          margin: 2px 0;
          font-size: 1.1rem;
        }

        .name, .user, .pass {
          flex-grow: 1;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          padding-right: 20px;
          min-height: 20px;
          box-sizing: border-box;
          vertical-align: middle;
          display: inline-block;
          line-height: 1.3;
        }

        .user {
          color: darken(#8a8f9f, 5%);
          font-weight: 500;
          font-size: 1rem;
        }

        .user, .pass {
          flex-basis: 35%;
        }

        &:not(.header) {

          &:hover {
            background-color: #191e23;
            color: #f7f8f8;
          }

          &:active {
            background-color: rgba(#2d2f36, .75);
            color: rgb(215, 216, 219);
            opacity: 1;
          }

          &.active {
            background-color: lighten(#1e2328, 3%);
          }

          .pass {
            position: relative;

            button {
              position: absolute;
              right: 0;
              background-color: transparent;
              border: none;
              padding: 0;
              outline: none;
              display: inline-block;
              opacity: 0;

              img {
                height: 20px;
                display: block;
              }
            }
          }
        }
      }
    }
  }

  .entries-list-header {
    height: 84px;
    padding: 12px 0;
    box-sizing: border-box;
  }

  //#entry {
  //  flex-basis: 100%;
  //}

  $accent: #00edb1;
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
</style>
