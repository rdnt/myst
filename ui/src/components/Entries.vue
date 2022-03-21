<template>
	<transition :duration="500" name="show">
		<div class="entries-list">
<!--			<div class="entry header">-->
<!--				<span class="name">-->
<!--					Domain-->
<!--					<button><img alt="" src="/assets/sort-arrow.svg"/></button>-->
<!--				</span>-->
<!--				<span class="user">-->
<!--					Username-->
<!--					<button><img alt="" src="/assets/sort-arrow.svg"/></button>-->
<!--				</span>-->
<!--				<span class="pass">-->
<!--					Password-->
<!--					<button><img alt="" src="/assets/sort-arrow.svg"/></button>-->
<!--				</span>-->
<!--			</div>-->
			<div v-if="keystore" class="entries">
				<router-link
					v-for="(e, i) in keystore.entries"
					:key="e.id"
					:class="{active: entry && entry.id === e.id}"
					:to="{
            name: 'entry',
            params: { keystoreId: keystore.id, entryId: e.id }
          }"
					class="entry"
				>
          <span class="icon">
            <img :src="`http://${e.website}/favicon.ico`" alt=""/>


          </span>

					<div class="info">
						<span class="name">
							{{ e.website }}
						</span>
						<span class="user">
							{{ i === 2 ? 'multiple accounts': e.username }}
						</span>
					</div>

<!--					<span class="pass">-->
<!--            {{ e.password.replace(/./g, "∗") }}-->
<!--						<button tabindex="-1"><img alt="" src="/assets/eye.svg"/></button>-->
<!--          </span>-->
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
	props: {},
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
		height: calc(100% - 0px);
		//padding: 0 20px;
		//padding-top: 23px;
		padding: 20px;
		box-sizing: border-box;
	}

	.entry {
		position: relative;
		display: flex;
		flex-direction: row;
		flex-wrap: nowrap;
		align-items: center;
		padding: 10px 14px;
		box-sizing: border-box;
		border-radius: 5px;
		min-height: 20px;
		text-decoration: none;
		margin-bottom: 2px;

		.info {
			display: flex;
			flex-direction: column;
		}

		//&:last-child {
		//	margin-bottom: 20px;
		//}

		&.header {
			color: rgb(138, 143, 152);
			height: 60px;
			padding: 0 34px;
			margin-top: 20px;

			.name {
				flex-basis: calc(30% + 34px);
			}

			span {
				color: rgba(#8a8f9f, .75);
				text-transform: uppercase;
				font-size: 0.85rem;
				font-weight: 600;
				letter-spacing: 0.5px;
				pointer-events: none;
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
			padding-right: 16px;

			img {
				display: block;
				width: 32px;
				height: 32px;
				vertical-align: baseline;
			}
		}

		.name {
			flex-basis: calc(30%);
			font-weight: 500;
			margin: 2px 0;
			font-size: 1.1rem;
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
			//font-size: 1.1rem;
		}

		.user {
			color: darken(#8a8f9f, 10%);
			font-weight: 500;
			font-size: 1rem;
		}

		.user, .pass {
			flex-basis: 35%;
		}

		&:not(.header) {

			&:hover {
				background-color: #191e23;
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
