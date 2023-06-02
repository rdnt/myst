<script lang="ts">
  import {Link, useLocation} from "svelte-navigator";

  export let path: string = '';
  const location = useLocation();
  export let active;
  $: active = $location.pathname.startsWith(path);
</script>

<div class="link" class:active style="max-width: 100%">
  <Link to={path}>
    <slot/>
  </Link>
</div>
