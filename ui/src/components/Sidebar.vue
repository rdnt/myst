<template>
	<div class="sidebar">
		<h4>Myst</h4>
		<div class="keystores-list">
			<div class="keystore" :class="{active: keystore.id === k.id}" v-for="k in keystores" :key="k.id">
				{{ k.name }}
			</div>
		</div>
	</div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {Keystore} from "../api/generated";

export default defineComponent({
	name: 'Sidebar',
	components: {},
	props: {
		keystores: Array as () => Keystore[],
		keystore: Object as () => Keystore,
	},
	data: () => ({}),
	computed: {},
	methods: {}
})
</script>

<style lang="scss" scoped>
h4 {
	font-weight: 700;
	font-size: 1.2rem;
	padding: 0 10px;
	margin: 0 0 20px;
	height: 20px;
}

.sidebar {
	background-color: #14191d;
	height: 100%;
	padding: 20px;
	box-sizing: border-box;
}

.keystores-list {
	display: flex;
	flex-direction: column;

	.keystore {
		border-radius: 5px;
		color: #fff;
		cursor: pointer;
		height: 22px;
		padding: 10px 100px 10px 10px;
		font-size: 1.1rem;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;


		&:hover {
			background-color: #1e2328;
		}

		&.active {
			background-color: #272d32;
		}
	}
}
</style>
