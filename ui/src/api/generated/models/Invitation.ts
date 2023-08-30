/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { InvitationKeystore } from './InvitationKeystore';
import type { User } from './User';

export type Invitation = {
    id: string;
    keystore: InvitationKeystore;
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
