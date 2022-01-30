<template>
	<router-view/>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {Keystore} from "../api/generated/index";

export default defineComponent({
	name: 'Keystore',
	components: {},
	props: {
		keystores: {
			type: Array as () => Keystore[],
			required: true,
		},
		keystore: {
			type: Object as () => Keystore,
			required: true,
		},
	},
	data: () => ({}),
	methods: {}
})
</script>

<style lang="scss" scoped>
$accent: #00edb1;

h4 {
	font-weight: 700;
	font-size: 1.6rem;
	margin: 0 0 66px;
	margin-top: 10px;
	height: 20px;
	padding: 0 20px;
}

h5 {
	font-size: 1rem;
	margin: 0;
	height: 20px;
	padding: 0 20px;
	margin-bottom: 10px;
	color: #8a8f9f;
	text-transform: uppercase;
	font-size: .7rem;
	font-weight: 600;
	letter-spacing: .5px;
}

.sidebar {
	background-color: #0a0e11;
	height: 100%;
	padding: 20px;
	box-sizing: border-box;
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
		padding: 10px 100px 10px 20px;
		font-size: 1.1rem;
		white-space: nowrap;
		text-overflow: ellipsis;
		margin-bottom: 4px;

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
