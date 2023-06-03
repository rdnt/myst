<script lang="ts">
  import Keystore from "@/components/Keystore.svelte";
  import {onMount} from "svelte";
  import {useLocation, useNavigate, useParams} from "svelte-navigator";
  import EntryPlaceholder from "@/components/EntryPlaceholder.svelte";
  import {getKeystores} from "@/stores/keystores";
  import CreateKeystoreModal from "@/components/CreateKeystoreModal.svelte";

  const navigate = useNavigate();
  const params = useParams();
  const location = useLocation();

  export let keystores;
  let keystore;

  export let showCreateKeystoreModal: boolean;

  $: {
    keystores;

    // redirect to root if no keystores
    if (keystores && keystores.length === 0 && $location.pathname !== '/') {
      navigate("/", {replace: true})
    }

    // if no keystore selected try to select the first one
    if (!$params.keystoreId && keystores?.length > 0) {
        navigate("/keystore/" + keystores?.[0].id, {replace: true});
    } else {
      // navigate to root if keystore not found, otherwise render
      const keystore = keystores.find((keystore) => keystore.id === $params.keystoreId);
      if (!keystore && $location.pathname !== '/') {
        console.log('redirect loop')
        navigate("/", {replace: true})
      }
    }
  }

  // onMount(async () => {
  //   if (!$params.keystoreId) {
  //     // TODO: always select first (right now it's buggy on login)
  //     // TODO: maybe select default keystore once that functionality is implemented
  //     // keystore = keystores[0];
  //     // keystore = keystores.find((keystore) => keystore.name === "Passwords");
  //     if (keystores.length > 0 ) {
  //       navigate("/keystore/" + keystores[0].id, {replace: true});
  //     }
  //
  //   } else {
  //       const keystore = keystores.find((keystore) => keystore.id === $params.keystoreId);
  //       if (!keystore) {
  //
  //         console.log("NO KEYSTORE FOUND??")
  //           navigate("/", {replace: true})
  //       }
  //   }
  // });

  $: keystore = (keystores || []).find(
    (keystore) => keystore.id === $params.keystoreId
  );
</script>

{#if keystore}
  <Keystore {keystore}/>
{:else}
  <div class="empty">
    <div class="title">Create a new keystore to get started</div>
    <button class="button green" on:click={() => {showCreateKeystoreModal = true}}>Create</button>
  </div>
{/if}

<style lang="scss">
  $accent: #00edb1;

  .empty {
    display: flex;
    flex-direction: column;
    position: relative;
    background-color: #101519;
    height: 100%;
    flex-basis: 100%;
    justify-content: center;
    align-items: center;
  }

  .title {
    margin-bottom: 20px;
  }

  .button {
    outline: none;
    border: none;
    height: 48px;
    font-size: 1.1rem;
    font-weight: 500;
    padding: 0 40px;
    border-radius: 5px;
    margin: 0 5px;
    background-color: rgba(#202228, 1);
    color: #fff;

    &.disabled {
      background-color: #161819;
    }

    &.green {
      background-color: #002e23;
      color: $accent;

      &.disabled {
        background-color: #0c1d19;
      }
    }
  }
</style>
