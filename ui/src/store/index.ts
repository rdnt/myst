import {defineStore} from 'pinia'
import {Entry, Keystore} from "../api/generated/index";

interface State {
    keystores: Keystore[];
    keystore?: Keystore;
    entry?: Entry;
}

const state: State = {
    keystores: [],
    keystore: undefined,
    entry: undefined,
}

export const useMainStore = defineStore('main', {
    // a function that returns a fresh state
    state: () => ({
        ...state
    }),
    // optional getters
    getters: {},
    // optional actions
    actions: {
        setKeystores(keystores: Keystore[]) {
            this.keystores = keystores.sort((a, b) => {
                return a.id < b.id ? 1 : -1;
            });
        },
        setKeystore(keystore: Keystore | undefined) {
            this.keystore = keystore;
        },
        setEntry(entry: Entry | undefined) {
            this.entry = entry;
        },
    },
})
