<script lang="ts">

</script>

<div class="entry-placeholder"></div>

<style lang="scss">
  .entry-placeholder {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    background-color: #101519;
    border-left: 2px solid #1a2025;
    height: 100%;
    padding: 20px;
    box-sizing: border-box;
    flex-basis: 50%;
    overflow-y: auto;

    &:after {
      content: "";
      background-color: #1b2025;
      border-radius: 50%;
      font-size: 1.5em;
      font-weight: bold;
      text-align: center;
      display: block;
      width: 100px;
      height: 100px;
      line-height: 100%;
    }

  }
</style>
