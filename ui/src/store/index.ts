import { defineStore } from 'pinia'

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
        setReady(ready: boolean) {
            this.ready = ready;
        },
        setOnboarding(onboarding: boolean) {
            this.onboarding = onboarding;
        }
    },
})
