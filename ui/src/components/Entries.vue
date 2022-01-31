<template>
	<transition :duration="500" name="show">
		<div class="entries-list">
			<div class="entry header">
					<span class="name">
						Domain
						<button><img alt="" src="/assets/sort-arrow.svg"/></button>
					</span>
				<span class="user">
						Username
						<button><img alt="" src="/assets/sort-arrow.svg"/></button>
					</span>
				<span class="pass">
						Password
						<button><img alt="" src="/assets/sort-arrow.svg"/></button>
					</span>
			</div>
			<div class="entries" v-if="keystore">
				<router-link
					v-for="e in keystore.entries"
					:key="e.id"
					:to="{
            name: 'entry',
            params: { keystoreId: this.keystore.id, entryId: e.id }
          }"
					class="entry"
					:class="{active: entry && entry.id === e.id}"
				>
          <span class="icon">
            <img :src="`http://${e.label}/favicon.ico`" alt=""/>
          </span>
					<span class="name">
            {{ e.label }}
          </span>
					<span class="user">
            {{ e.username }}
          </span>
					<span class="pass">
            {{ e.password.replace(/./g, "∗") }}
						<button tabindex="-1"	><img alt="" src="/assets/eye.svg"/></button>
          </span>
				</router-link>
			</div>
		</div>
	</transition>
	<KeystoreEntry :entry="entry"/>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import KeystoreEntry from "./Entry.vue";
import {Entry, Keystore} from "../api/generated";
import {useMainStore} from "../store";

export default defineComponent({
	name: 'Entries',
	components: {KeystoreEntry},
	setup() {
		const main = useMainStore()

		return {
			main,
		}
	},
	props: {
	},
	data() {
		return {}
	},
	computed: {
		keystores(): Keystore[] {
			return this.main.keystores
		},
		keystore(): Keystore | undefined {
			return this.main.keystore
		},
		entry(): Entry | undefined {
			return this.main.entry
		},
	},
	watch: {
		// $route: {
			// handler: function (route) {
			// 	console.log('handler')
			// 	if (route.params.keystoreId) {
			// 		console.log("keystoreId", route.params.keystoreId);
			// 		api.keystore(route.params.keystoreId).then((keystore) => {
			// 			this.keystore = keystore
			// 			console.log('keystore set')
			// 		})
			// 	} else {
			// 		this.keystore = undefined
			// 	}
			// 	// if (!route.params.entryId) {
			// 	// 	this.entry = undefined;
			// 	// } else {
			// 	// 	this.entry = this.entries.find(entry => entry.id === route.params.entryId);
			// 	// }
			// }
		// }
	},
	methods: {}
})
</script>

<style lang="scss" scoped>



//transition: transform .5s cubic-bezier(.68,.09,.13,.89);

.entries-list {
	position: relative;
	background-color: #101519;
	height: 100%;
	flex-grow: 1;

	.entries {
		overflow-y: auto;
		height: calc(100% - 60px);
		padding: 0 20px;
	}

	.entry {
		position: relative;
		display: flex;
		flex-direction: row;
		flex-wrap: nowrap;
		justify-content: space-between;
		align-items: center;
		padding: 10px 14px;
		box-sizing: border-box;
		border-radius: 5px;
		min-height: 20px;
		text-decoration: none;
		margin-bottom: 2px;

		&:last-child {
			margin-bottom: 20px;
		}

		&.header {
			color: rgb(138, 143, 152);
			height: 60px;
			padding: 0 34px;

			.name {
				flex-basis: calc(30% + 34px);
			}

			span {
				color: rgb(138, 143, 152);
				font-size: 1rem;
			}

			button {
				background-color: transparent;
				border: none;
				padding: 0;
				margin: 0;

				img {
					width: 10px;
					height: 10px;
					margin-right: 10px;
					opacity: .5;
				}
			}
		}

		.icon {
			flex-basis: 24px;
			padding-right: 10px;

			img {
				display: block;
				height: 20px;
				vertical-align: baseline;
			}
		}

		.name {
			flex-basis: calc(30%);
		}

		.name, .user, .pass {
			//flex-basis: 0;
			flex-grow: 1;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
			padding-right: 20px;
			min-height: 20px;
			box-sizing: border-box;
			vertical-align: middle;
			display: inline-block;
			line-height: 1.3;
			font-size: 1.1rem;
		}

		.user, .pass {
			flex-basis: 35%;
		}

		&:not(.header) {

			&:hover {
				background-color: #1e2328;
				color: #f7f8f8;
			}

			&:active {
				background-color: rgba(#2d2f36, .75);
				color: rgb(215, 216, 219);
				opacity: 1;
			}

			&.active {
				background-color: lighten(#1e2328, 3%);
			}

			.pass {
				position: relative;
				button {
					position: absolute;
					right: 0;
					background-color: transparent;
					border: none;
					padding: 0;
					outline: none;
					display: inline-block;
					opacity: 0;

					img {
						height: 20px;
						display: block;
					}
				}
			}
		}
	}
}

#entry {
	flex-basis: 100%;
}


</style>
