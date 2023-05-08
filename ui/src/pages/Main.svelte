<script lang="ts">
  import SignInModal from "@/components/SignInModal.svelte";
  import RegisterModal from "@/components/RegisterModal.svelte";
  import CreateKeystoreModal from "@/components/CreateKeystoreModal.svelte";
  import Sidebar from "@/components/Sidebar.svelte";
  import Invitations from "@/pages/Invitations.svelte";
  import Keystores from "@/pages/Keystores.svelte";
  import {getKeystores} from "@/stores/keystores.ts";
  import {Route, Router} from "svelte-navigator";
  import {onMount} from "svelte";

  let showCreateKeystoreModal: boolean = false;
  let showSignInModal: boolean = false;
  let showRegisterModal: boolean = false;
</script>
<Router>
  <Sidebar
           bind:showCreateKeystoreModal={showCreateKeystoreModal}
           bind:showSignInModal={showSignInModal}
           bind:showRegisterModal={showRegisterModal}
  />

  <main>
    <Route>
      <Keystores keystores={$keystores}/>
    </Route>

<!--    <Route path="/keystore/:keystoreId">-->
<!--      <Keystores keystores={$keystores}/>-->
<!--    </Route>-->

<!--    <Route path="/keystore/:keystoreId/entry/:entryId">-->
<!--      <Keystores keystores={$keystores}/>-->
<!--    </Route>-->

    <Route path="/invitations">
      <Invitations/>
    </Route>
  </main>
</Router>

<CreateKeystoreModal bind:show={showCreateKeystoreModal} on:created={() => {getKeystores()}}/>

<SignInModal bind:show={showSignInModal}/>
<RegisterModal bind:show={showRegisterModal}/>

<style lang="scss">
  $bg: #0a0e11;

  main {
    display: flex;
    align-items: stretch;
    background-color: $bg;
    width: 100%;
  }
</style>
