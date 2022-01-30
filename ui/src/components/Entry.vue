<template>
	<div class="entry">
		DATA: {{JSON.stringify(entry)}}
	</div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {Entry} from "../api/generated";

export default defineComponent({
	name: 'Entry',
	components: {},
	props: {
		entry: {
			type: Object as () => Entry,
			required: true
		}
	},
	data: () => ({}),
	computed: {},
	methods: {}
})
</script>

<style lang="scss" scoped>
.entry {
	background-color: #101519;
	border-left: 2px solid #1a2025;
	height: 100%;
	padding: 20px;
	box-sizing: border-box;
	flex-basis: 40%;
}
</style>
