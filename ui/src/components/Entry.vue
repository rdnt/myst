<template>
	<div class="entry" :class="{empty: !entry}">
		<template v-if="entry">
			DATA: {{JSON.stringify(entry, null, 2)}}
		</template>
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
			required: false
		}
	},
	data: () => ({}),
	computed: {},
	methods: {}
})
</script>

<style lang="scss" scoped>
.entry {
	display: flex;
	align-items: center;
	justify-content: center;
	background-color: #101519;
	border-left: 2px solid #1a2025;
	height: 100%;
	padding: 20px;
	box-sizing: border-box;
	flex-basis: 40%;

	&.empty {
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
}
</style>
