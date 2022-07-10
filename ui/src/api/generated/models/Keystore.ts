/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Entry } from './Entry';

export type Keystore = {
    id: string;
    remoteId: string;
    name: string;
    entries: Array<Entry>;
};