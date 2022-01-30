import {createRouter, createWebHistory} from "vue-router";
import Keystore from "../views/Keystore.vue";
import Entries from "../components/Entries.vue";
import Entry from "../components/Entry.vue";

export default createRouter({
    history: createWebHistory(),
    routes: [
        {
            name: 'keystore',
            path: '/',
            component: Keystore,
            // props: true,
            children: [
                {
                    path: '/keystore/:keystoreId',
                    name: 'entries',
                    component: Entries,
                    // props: true,
                    children: [
                        {
                            path: '/keystore/:keystoreId/entry/:entryId',
                            name: 'entry',
                            component: Entry,
                            // props: true,
                        },
                    ]
                },
            ]
        }
    ],
})
