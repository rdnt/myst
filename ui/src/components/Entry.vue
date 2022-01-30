<template>
	<div class="entry">
		entry content
	</div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {Entry} from "../api/generated";

export default defineComponent({
	name: 'Entry',
	components: {},
	props: {
		entry: Object as () => Entry,
	},
	data: () => ({}),
	computed: {},
	methods: {}
})
</script>

<style lang="scss" scoped>
.entry {
	background-color: #14191d;
	height: 100%;
	padding: 20px;
	box-sizing: border-box;
	flex-basis: 40%;
}
</style>
