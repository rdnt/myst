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
                    name: 'entries',
                    path: '/keystore/:keystoreId',
                    component: Entries,
                    // props: true,
                    children: [
                        {
                            name: 'entry',
                            path: '/keystore/:keystoreId/entry/:entryId',
                            component: Entry,
                            // props: true,
                        },
                    ]
                },
            ]
        }
    ],
})
