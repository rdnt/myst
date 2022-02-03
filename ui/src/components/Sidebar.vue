<template>
	<div class="sidebar">
		<h4>Myst</h4>
		<div class="keystores-list">
			<h5>Keystores</h5>
			<router-link
				v-for="k in keystores"
				:key="k.id"
				:class="{active: keystore && keystore.id === k.id}"
				:to="{
            name: 'entries',
            params: { keystoreId: k.id }
          }"
				class="keystore">
				{{ k.name }}
			</router-link>
		</div>
	</div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {useMainStore} from "../store";
import {Keystore} from "../api/generated/index";

export default defineComponent({
	name: 'Sidebar',
	components: {},
	setup() {
		const main = useMainStore()

		return {
			main,
		}
	},
	computed: {
		keystores(): Keystore[] {
			return this.main.keystores
		},
		keystore(): Keystore | undefined {
			return this.main.keystore
		},
	},
	data: () => ({
	}),
	methods: {}
})
</script>

<style lang="scss" scoped>
$accent: #00edb1;

h4 {
	font-weight: 700;
	font-size: 2rem;
	margin: 10px 0 60px;
	padding: 0 20px;
}

h5 {
	height: 20px;
	padding: 0 20px;
	margin: 0 0 10px;
	color: #8a8f9f;
	text-transform: uppercase;
	font-size: .85rem;
	font-weight: 600;
	letter-spacing: .5px;
}

.sidebar {
	background-color: #0a0e11;
	height: 100%;
	padding: 20px;
	box-sizing: border-box;
	flex-basis: 250px;
}

.keystores-list {
	display: flex;
	flex-direction: column;

	.keystore {
		display: flex;
		align-items: center;
		border-radius: 5px;
		position: relative;
		color: #fff;
		cursor: pointer;
		height: 22px;
		padding: 10px 20px 10px 20px;
		font-size: 1.1rem;
		white-space: nowrap;
		text-overflow: ellipsis;
		margin-bottom: 2px;
		text-decoration: none;

		&:hover {
			background-color: rgba(#1e2328, .75);
		}

		&.active {
			font-weight: 500;
			color: #00edb1;
			background-color: #0c1d19;

			&:after {
				content: '';
				position: absolute;
				top: 0;
				width: 10px;
				height: 100%;
				left: -20px;
				border-top-right-radius: 5px;
				border-bottom-right-radius: 5px;
			}
		}
	}
}
</style>
