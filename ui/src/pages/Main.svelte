<script lang="ts">
  import RegisterModal from "@/components/RegisterModal.svelte";
  import CreateKeystoreModal from "@/components/CreateKeystoreModal.svelte";
  import Sidebar from "@/components/Sidebar.svelte";
  import Invitations from "@/pages/Invitations.svelte";
  import Keystores from "@/pages/Keystores.svelte";
  import {getKeystores, keystores} from "@/stores/keystores";
  import {Route, Router} from "svelte-navigator";
  import {onDestroy, onMount} from "svelte";
  import {getCurrentUser} from "@/stores/user";
  import {getInvitations} from "@/stores/invitations";
  import {useNavigate} from "svelte-navigator";
  import type {Keystore} from "@/api";
  import api from "@/api";

  const navigate = useNavigate();

  let showCreateKeystoreModal: boolean = false;
  let showSignInModal: boolean = false;
  let showRegisterModal: boolean = false;

  let ready: boolean = false;
  $: ready;

  onMount(async () => {
    await getKeystores()
    const u = await getCurrentUser()
    if (u) {
      await getInvitations()
    }
    ready = true
  })

  const onKeystoreCreated = async (keystore: Keystore) => {
    await getKeystores()
    navigate(`/keystore/${keystore.id}`)
  }
</script>

{#if ready}
  <Router>
    <Sidebar
             bind:showCreateKeystoreModal={showCreateKeystoreModal}
             bind:showSignInModal={showSignInModal}
             bind:showRegisterModal={showRegisterModal}
             keystores={$keystores}
    />

    <main>
      <Route>
        <Keystores keystores={$keystores} bind:showCreateKeystoreModal={showCreateKeystoreModal}/>
      </Route>

      <Route path="/keystore/:keystoreId">
        <Keystores keystores={$keystores} bind:showCreateKeystoreModal={showCreateKeystoreModal}/>
      </Route>

      <Route path="/keystore/:keystoreId/entry/:entryId">
        <Keystores keystores={$keystores} bind:showCreateKeystoreModal={showCreateKeystoreModal}/>
      </Route>

      <Route path="/invitations">
        <Invitations/>
      </Route>
    </main>
  </Router>

  <CreateKeystoreModal bind:show={showCreateKeystoreModal} on:created={(e) => {onKeystoreCreated(e.detail)}}/>

  <!--<SignInModal bind:show={showSignInModal}/>-->
  <RegisterModal bind:show={showRegisterModal} />

{/if}

<style lang="scss">
  $bg: #0a0e11;

  main {
    position: relative;
    display: flex;
    align-items: stretch;
    background-color: $bg;
    flex-grow: 0;
    flex-basis: 100%;
    min-width: 0;
  }
</style>
