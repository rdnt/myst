/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { AuthenticateRequest } from '../models/AuthenticateRequest';
import type { CreateEntryRequest } from '../models/CreateEntryRequest';
import type { CreateInvitationRequest } from '../models/CreateInvitationRequest';
import type { CreateKeystoreRequest } from '../models/CreateKeystoreRequest';
import type { Entry } from '../models/Entry';
import type { Error } from '../models/Error';
import type { Invitation } from '../models/Invitation';
import type { Invitations } from '../models/Invitations';
import type { Keystore } from '../models/Keystore';
import type { Keystores } from '../models/Keystores';
import type { UpdateEntryRequest } from '../models/UpdateEntryRequest';
import type { User } from '../models/User';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class DefaultService {

    /**
     * Triggers a health check
     * @returns any OK
     * @throws ApiError
     */
    public static healthCheck(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/health',
        });
    }

    /**
     * Attempts to authenticate the user
     * @returns any OK
     * @returns Error Error
     * @throws ApiError
     */
    public static authenticate({
requestBody,
}: {
requestBody: AuthenticateRequest,
}): CancelablePromise<any | Error> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/authenticate',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Returns the currently signed in user
     * @returns User OK
     * @returns Error Error
     * @throws ApiError
     */
    public static currentUser(): CancelablePromise<User | Error> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/user',
            errors: {
                404: `Not Found`,
            },
        });
    }

    /**
     * Creates a new encrypted keystore with the given password
     * @returns Error Error
     * @returns Keystore Keystore
     * @throws ApiError
     */
    public static createKeystore({
requestBody,
}: {
requestBody: CreateKeystoreRequest,
}): CancelablePromise<Error | Keystore> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/keystores',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Returns all keystores
     * @returns Keystores Keystores
     * @returns Error Error
     * @throws ApiError
     */
    public static keystores(): CancelablePromise<Keystores | Error> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/keystores',
        });
    }

    /**
     * Get a keystore if it exists and return it
     * @returns Keystore Keystore
     * @returns Error Error
     * @throws ApiError
     */
    public static keystore({
keystoreId,
}: {
/** unique identifier for a keystore **/
keystoreId: string,
}): CancelablePromise<Keystore | Error> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/keystore/{keystoreId}',
            path: {
                'keystoreId': keystoreId,
            },
        });
    }

    /**
     * Creates a new entry and adds it to the keystore
     * @returns Entry Entry
     * @returns Error Error
     * @throws ApiError
     */
    public static createEntry({
keystoreId,
requestBody,
}: {
/** unique identifier for a keystore **/
keystoreId: string,
requestBody: CreateEntryRequest,
}): CancelablePromise<Entry | Error> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/keystore/{keystoreId}/entries',
            path: {
                'keystoreId': keystoreId,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Get the entries of a keystore
     * @returns Entry Entries
     * @returns Error Error
     * @throws ApiError
     */
    public static keystoreEntries({
keystoreId,
}: {
/** unique identifier for a keystore **/
keystoreId: string,
}): CancelablePromise<Array<Entry> | Error> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/keystore/{keystoreId}/entries',
            path: {
                'keystoreId': keystoreId,
            },
        });
    }

    /**
     * Update a keystore entry
     * @returns Entry Entry
     * @returns Error Error
     * @throws ApiError
     */
    public static updateEntry({
keystoreId,
entryId,
requestBody,
}: {
keystoreId: string,
entryId: string,
requestBody: UpdateEntryRequest,
}): CancelablePromise<Entry | Error> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/keystore/{keystoreId}/entry/{entryId}',
            path: {
                'keystoreId': keystoreId,
                'entryId': entryId,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Delete a keystore entry
     * @returns any OK
     * @returns Error Error
     * @throws ApiError
     */
    public static deleteEntry({
keystoreId,
entryId,
}: {
keystoreId: string,
entryId: string,
}): CancelablePromise<any | Error> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/keystore/{keystoreId}/entry/{entryId}',
            path: {
                'keystoreId': keystoreId,
                'entryId': entryId,
            },
        });
    }

    /**
     * Create a keystore invitation
     * @returns Invitation Invitation
     * @returns Error Error
     * @throws ApiError
     */
    public static createInvitation({
requestBody,
keystoreId,
}: {
requestBody: CreateInvitationRequest,
keystoreId?: string,
}): CancelablePromise<Invitation | Error> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/keystore/{keystoreId}/invitations',
            path: {
                'keystoreId': keystoreId,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Accepts an invitation
     * @returns Invitation Invitation
     * @returns Error Error
     * @throws ApiError
     */
    public static acceptInvitation({
invitationId,
}: {
invitationId?: string,
}): CancelablePromise<Invitation | Error> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/invitation/{invitationId}',
            path: {
                'invitationId': invitationId,
            },
        });
    }

    /**
     * Returns an invitation
     * @returns Invitation Invitation
     * @returns Error Error
     * @throws ApiError
     */
    public static getInvitation({
invitationId,
}: {
invitationId?: string,
}): CancelablePromise<Invitation | Error> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/invitation/{invitationId}',
            path: {
                'invitationId': invitationId,
            },
        });
    }

    /**
     * Get all keystore invitations
     * @returns Invitations Invitations
     * @returns Error Error
     * @throws ApiError
     */
    public static getInvitations(): CancelablePromise<Invitations | Error> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/invitations',
        });
    }

}