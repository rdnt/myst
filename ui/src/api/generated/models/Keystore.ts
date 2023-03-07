/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Entry } from './Entry';

export type Keystore = {
    id: string;
    remoteId: string;
    readOnly: boolean;
    name: string;
    entries: Array<Entry>;
};
