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
     * @param requestBody 
     * @returns any OK
     * @returns Error Error
     * @throws ApiError
     */
    public static authenticate(
requestBody: AuthenticateRequest,
): CancelablePromise<any | Error> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/authenticate',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Creates a new encrypted keystore with the given password
     * @param requestBody 
     * @returns Error Error
     * @returns Keystore Keystore
     * @throws ApiError
     */
    public static createKeystore(
requestBody: CreateKeystoreRequest,
): CancelablePromise<Error | Keystore> {
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
     * @param keystoreId unique identifier for a keystore
     * @returns Keystore Keystore
     * @returns Error Error
     * @throws ApiError
     */
    public static keystore(
keystoreId: string,
): CancelablePromise<Keystore | Error> {
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
     * @param keystoreId unique identifier for a keystore
     * @param requestBody 
     * @returns Entry Entry
     * @returns Error Error
     * @throws ApiError
     */
    public static createEntry(
keystoreId: string,
requestBody: CreateEntryRequest,
): CancelablePromise<Entry | Error> {
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
     * @param keystoreId unique identifier for a keystore
     * @returns Entry Entries
     * @returns Error Error
     * @throws ApiError
     */
    public static keystoreEntries(
keystoreId: string,
): CancelablePromise<Array<Entry> | Error> {
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
     * @param keystoreId 
     * @param entryId 
     * @param requestBody 
     * @returns Entry Entry
     * @returns Error Error
     * @throws ApiError
     */
    public static updateEntry(
keystoreId: string,
entryId: string,
requestBody: UpdateEntryRequest,
): CancelablePromise<Entry | Error> {
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
     * @param keystoreId 
     * @param entryId 
     * @returns any OK
     * @returns Error Error
     * @throws ApiError
     */
    public static deleteEntry(
keystoreId: string,
entryId: string,
): CancelablePromise<any | Error> {
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
     * @param requestBody 
     * @param keystoreId 
     * @returns Invitation Invitation
     * @returns Error Error
     * @throws ApiError
     */
    public static createInvitation(
requestBody: CreateInvitationRequest,
keystoreId?: string,
): CancelablePromise<Invitation | Error> {
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