<script lang="ts">
  import AuthModal from "@/components/AuthModal.svelte";
  import CreateKeystoreModal from "@/components/CreateKeystoreModal.svelte";
  import Sidebar from "@/components/Sidebar.svelte";
  import Invitations from "@/pages/Invitations.svelte";
  import Keystores from "@/pages/Keystores.svelte";
  import {getKeystores, keystores} from "@/stores/keystores.ts";
  import {Route, Router} from "svelte-navigator";
  import {currentUser} from "@/stores/user";

  let showCreateKeystoreModal: boolean = false;
  let showAuthModal: boolean = false;
</script>
<Router>
  <Sidebar keystores={$keystores} bind:showCreateKeystoreModal={showCreateKeystoreModal} bind:showAuthModal={showAuthModal}/>

  <main>
    <Route>
      <Keystores keystores={$keystores}/>
    </Route>

    <Route path="/keystore/:keystoreId">
      <Keystores keystores={$keystores}/>
    </Route>

    <Route path="/keystore/:keystoreId/entry/:entryId">
      <Keystores keystores={$keystores}/>
    </Route>

    {#if $currentUser}
      <Route path="/invitations">
        <Invitations/>
      </Route>
    {/if}
  </main>
</Router>

<CreateKeystoreModal bind:show={showCreateKeystoreModal} on:created={() => {getKeystores()}}/>

<AuthModal bind:show={showAuthModal}/>
