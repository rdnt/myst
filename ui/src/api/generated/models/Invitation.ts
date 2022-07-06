/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type Invitation = {
    id: string;
    keystoreId: string;
    keystoreName: string;
    inviterId: string;
    inviteeId: string;
    inviterPublicKey: string;
    inviteePublicKey: string;
    encryptedKeystoreKey: string;
    status: 'pending' | 'accepted' | 'declined' | 'deleted' | 'finalized';
    createdAt: string;
    updatedAt: string;
    acceptedAt: string;
    declinedAt: string;
    deletedAt: string;
};