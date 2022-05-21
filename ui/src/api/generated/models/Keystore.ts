/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Entry } from './Entry';

export type Keystore = {
    id: string;
    name: string;
    entries: Array<Entry>;
};