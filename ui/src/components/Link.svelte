<script lang="ts">
  import {Link, useLocation} from "svelte-navigator";

  export let path: string = '';
  const location = useLocation();
  let active;
  $: active = $location.pathname.startsWith(path);
</script>

<div class="link" class:active>
  <Link to={path}>
    <slot/>
  </Link>
</div>
