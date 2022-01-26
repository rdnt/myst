import { defineStore } from 'pinia'
import api from "../api"

export const mainStore = defineStore('main', {
    // a function that returns a fresh state
    state: () => ({
        ready: false,
        onboarding: false
    }),
    // optional getters
    getters: {
        // // getters receive the state as first parameter
        // doubleCount: (state) => state.counter * 2,
        // // use getters in other getters
        // doubleCountPlusOne(): number {
        //     return this.doubleCount * 2 + 1
        // },
    },
    // optional actions
    actions: {
        // setReady(ready: boolean) {
        //     this.ready = ready;
        // },
        // setOnboarding(onboarding: boolean) {
        //     this.onboarding = onboarding;
        // },
        load() {
            this.ready = true
            console.log(api)
            api.keystoreIds().then((ids) => {
                this.onboarding = ids.length == 0;
            }).finally(() => {
                this.ready = true
            })
        }
    },
})
