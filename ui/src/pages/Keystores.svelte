<script lang="ts">
  import {useNavigate, useParams} from "svelte-navigator";
  import Keystore from "../components/Keystore.svelte";
  import {onMount} from "svelte";

  const navigate = useNavigate();
  const params = useParams();

  export let keystores;
  let keystore;

  onMount(() => {
    if (!$params.keystoreId) {
      // TODO: always select first
      // TODO: maybe select default keystore once that functionality is implemented
      // keystore = keystores[0];
      keystore = keystores.find((keystore) => keystore.name === "Passwords");
      navigate("/keystore/"+ keystore.id)
    }
  })

  $: keystore = (keystores || []).find((keystore) => keystore.id === $params.keystoreId);
</script>

{#if keystore}
  <Keystore {keystore} />
{/if}
