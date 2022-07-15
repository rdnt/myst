/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Keystore } from './Keystore';
import type { User } from './User';

export type Invitation = {
    id: string;
    keystore: Keystore;
    inviter: User;
    invitee: User;
    encryptedKeystoreKey: string;
    status: 'pending' | 'accepted' | 'declined' | 'deleted' | 'finalized';
    createdAt: string;
    updatedAt: string;
    acceptedAt: string;
    declinedAt: string;
    deletedAt: string;
};