<script lang="ts">
  export let ref = null;
  export let label: string;
  export let value: string;

  import {createEventDispatcher} from 'svelte';

  const dispatch = createEventDispatcher();
</script>

<div class={$$props.class} class:field={true}>
  <label>{label}</label>
  <div class="value">{value || '—'}</div>
</div>

<style lang="scss">
  .field {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    margin-bottom: 32px;
    align-items: center;

    label {
      flex-basis: 100%;
      color: darken(#8a8f9f, 5%);
      text-transform: uppercase;
      font-size: 0.9rem;
      font-weight: 600;
      letter-spacing: 0.5px;
      pointer-events: none;
    }

    .value {
      width: 100%;
      display: block;
      font-size: 1.1rem;
      font-weight: 400;
      box-sizing: border-box;
      padding-top: 6px;
      overflow: hidden;
    }

    //span {
    //  color: #ff9999;
    //  font-weight: 500;
    //  font-size: .9rem;
    //  padding: 0 16px;
    //  opacity: 0;
    //  pointer-events: none;
    //  margin-top: 12px;
    //  min-height: 20px;
    //
    //  &.show {
    //    opacity: 1;
    //    pointer-events: unset;
    //  }
    //}
  }
</style>
