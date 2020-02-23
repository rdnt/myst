<template>
  <div id="two-factor">
    <div class="header">
      <div class="label">
        <span>Two Factor Authentication</span>
      </div>
    </div>
    <div class="items-container">
      <div class="items">
        <div
          v-for="item in items"
          class="item card"
          :tabindex="offset + item.id + 1"
          :key="item.id"
        >
          <span class="code">980 501</span>
          <span class="name">Google</span>
          <div class="timer">
            <div class="anim inner">
              <div class="mask"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  computed: mapState({
    offset: state => state.sites.all.length
  }),
  data() {
    return {
      items: [
        {
          id: 0
        },
        {
          id: 1
        },
        {
          id: 2
        },
        {
          id: 3
        },
        {
          id: 4
        },
        {
          id: 5
        },
        {
          id: 6
        },
        {
          id: 7
        },
        {
          id: 8
        },
        {
          id: 9
        },
        {
          id: 10
        },
        {
          id: 11
        },
        {
          id: 12
        },
        {
          id: 13
        },
        {
          id: 14
        }
      ]
    };
  }
};
</script>
