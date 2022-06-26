<script lang="ts">
  import Keystore from "@/components/Keystore.svelte";
  import {onMount} from "svelte";
  import {useNavigate, useParams} from "svelte-navigator";

  const navigate = useNavigate();
  const params = useParams();

  export let keystores;
  let keystore;

  onMount(() => {
    console.log('keystore mount')
    if (!$params.keystoreId) {
      // TODO: always select first
      // TODO: maybe select default keystore once that functionality is implemented
      // keystore = keystores[0];
      // keystore = keystores.find((keystore) => keystore.name === "Passwords");
      console.log(keystores)
      if (keystores.length > 0 ) {
        navigate("/keystore/" + keystores[0].id);
      }

    }
  });

  $: keystore = (keystores || []).find(
    (keystore) => keystore.id === $params.keystoreId
  );
</script>

{#if keystore}
  <Keystore {keystore}/>
{/if}
