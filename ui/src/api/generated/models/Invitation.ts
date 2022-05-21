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
    status: Invitation.status;
    createdAt: number;
    updatedAt: number;
};

export namespace Invitation {

    export enum status {
        PENDING = 'pending',
        ACCEPTED = 'accepted',
        REJECTED = 'rejected',
        FINALIZED = 'finalized',
    }


}