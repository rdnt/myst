/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type Invitation = {
    id: string;
    keystoreId: string;
    inviterId: string;
    inviteeId: string;
    inviterKey: string;
    inviteeKey: string;
    keystoreKey: string;
    status: 'pending' | 'accepted' | 'rejected' | 'finalized';
    createdAt: number;
    updatedAt: number;
};