/* tslint:disable */
/* eslint-disable */
/**
 * Myst REST API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: apiteam@swagger.io
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import {
    AuthenticateRequest,
    AuthenticateRequestFromJSON,
    AuthenticateRequestToJSON,
    CreateEntryRequest,
    CreateEntryRequestFromJSON,
    CreateEntryRequestToJSON,
    CreateKeystoreRequest,
    CreateKeystoreRequestFromJSON,
    CreateKeystoreRequestToJSON,
    Entry,
    EntryFromJSON,
    EntryToJSON,
    Keystore,
    KeystoreFromJSON,
    KeystoreToJSON,
    UnlockKeystoreRequest,
    UnlockKeystoreRequestFromJSON,
    UnlockKeystoreRequestToJSON,
} from '../models';

export interface AuthenticateOperationRequest {
    authenticateRequest: AuthenticateRequest;
}

export interface CreateEntryOperationRequest {
    keystoreId: string;
    createEntryRequest: CreateEntryRequest;
}

export interface CreateKeystoreOperationRequest {
    createKeystoreRequest: CreateKeystoreRequest;
}

export interface GetKeystoreRequest {
    keystoreId: string;
}

export interface KeystoreRequest {
    keystoreId: string;
}

export interface UnlockKeystoreOperationRequest {
    keystoreId: string;
    unlockKeystoreRequest: UnlockKeystoreRequest;
}

/**
 * 
 */
export class DefaultApi extends runtime.BaseAPI {

    /**
     * Attempts to authenticate the user
     */
    async authenticateRaw(requestParameters: AuthenticateOperationRequest, initOverrides?: RequestInit): Promise<runtime.ApiResponse<void>> {
        if (requestParameters.authenticateRequest === null || requestParameters.authenticateRequest === undefined) {
            throw new runtime.RequiredError('authenticateRequest','Required parameter requestParameters.authenticateRequest was null or undefined when calling authenticate.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/authenticate`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: AuthenticateRequestToJSON(requestParameters.authenticateRequest),
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * Attempts to authenticate the user
     */
    async authenticate(requestParameters: AuthenticateOperationRequest, initOverrides?: RequestInit): Promise<void> {
        await this.authenticateRaw(requestParameters, initOverrides);
    }

    /**
     * Creates a new entry and adds it to the keystore
     */
    async createEntryRaw(requestParameters: CreateEntryOperationRequest, initOverrides?: RequestInit): Promise<runtime.ApiResponse<Entry>> {
        if (requestParameters.keystoreId === null || requestParameters.keystoreId === undefined) {
            throw new runtime.RequiredError('keystoreId','Required parameter requestParameters.keystoreId was null or undefined when calling createEntry.');
        }

        if (requestParameters.createEntryRequest === null || requestParameters.createEntryRequest === undefined) {
            throw new runtime.RequiredError('createEntryRequest','Required parameter requestParameters.createEntryRequest was null or undefined when calling createEntry.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/keystore/{keystoreId}/entries`.replace(`{${"keystoreId"}}`, encodeURIComponent(String(requestParameters.keystoreId))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: CreateEntryRequestToJSON(requestParameters.createEntryRequest),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => EntryFromJSON(jsonValue));
    }

    /**
     * Creates a new entry and adds it to the keystore
     */
    async createEntry(requestParameters: CreateEntryOperationRequest, initOverrides?: RequestInit): Promise<Entry> {
        const response = await this.createEntryRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Creates a new encrypted keystore with the given password
     */
    async createKeystoreRaw(requestParameters: CreateKeystoreOperationRequest, initOverrides?: RequestInit): Promise<runtime.ApiResponse<Keystore>> {
        if (requestParameters.createKeystoreRequest === null || requestParameters.createKeystoreRequest === undefined) {
            throw new runtime.RequiredError('createKeystoreRequest','Required parameter requestParameters.createKeystoreRequest was null or undefined when calling createKeystore.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/keystores`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: CreateKeystoreRequestToJSON(requestParameters.createKeystoreRequest),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => KeystoreFromJSON(jsonValue));
    }

    /**
     * Creates a new encrypted keystore with the given password
     */
    async createKeystore(requestParameters: CreateKeystoreOperationRequest, initOverrides?: RequestInit): Promise<Keystore> {
        const response = await this.createKeystoreRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get a keystore if it exists and returns it
     */
    async getKeystoreRaw(requestParameters: GetKeystoreRequest, initOverrides?: RequestInit): Promise<runtime.ApiResponse<Array<Keystore>>> {
        if (requestParameters.keystoreId === null || requestParameters.keystoreId === undefined) {
            throw new runtime.RequiredError('keystoreId','Required parameter requestParameters.keystoreId was null or undefined when calling getKeystore.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/keystore/{keystoreId}/entries`.replace(`{${"keystoreId"}}`, encodeURIComponent(String(requestParameters.keystoreId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(KeystoreFromJSON));
    }

    /**
     * Get a keystore if it exists and returns it
     */
    async getKeystore(requestParameters: GetKeystoreRequest, initOverrides?: RequestInit): Promise<Array<Keystore>> {
        const response = await this.getKeystoreRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Triggers a health check
     */
    async healthCheckRaw(initOverrides?: RequestInit): Promise<runtime.ApiResponse<void>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/health`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * Triggers a health check
     */
    async healthCheck(initOverrides?: RequestInit): Promise<void> {
        await this.healthCheckRaw(initOverrides);
    }

    /**
     * Get a keystore if it exists and return it
     */
    async keystoreRaw(requestParameters: KeystoreRequest, initOverrides?: RequestInit): Promise<runtime.ApiResponse<Array<Keystore>>> {
        if (requestParameters.keystoreId === null || requestParameters.keystoreId === undefined) {
            throw new runtime.RequiredError('keystoreId','Required parameter requestParameters.keystoreId was null or undefined when calling keystore.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/keystore/{keystoreId}`.replace(`{${"keystoreId"}}`, encodeURIComponent(String(requestParameters.keystoreId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(KeystoreFromJSON));
    }

    /**
     * Get a keystore if it exists and return it
     */
    async keystore(requestParameters: KeystoreRequest, initOverrides?: RequestInit): Promise<Array<Keystore>> {
        const response = await this.keystoreRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Returns all keystores
     */
    async keystoresRaw(initOverrides?: RequestInit): Promise<runtime.ApiResponse<Array<Keystore>>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/keystores`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(KeystoreFromJSON));
    }

    /**
     * Returns all keystores
     */
    async keystores(initOverrides?: RequestInit): Promise<Array<Keystore>> {
        const response = await this.keystoresRaw(initOverrides);
        return await response.value();
    }

    /**
     * Unlocks a keystore if it exists and returns it
     */
    async unlockKeystoreRaw(requestParameters: UnlockKeystoreOperationRequest, initOverrides?: RequestInit): Promise<runtime.ApiResponse<Array<Keystore>>> {
        if (requestParameters.keystoreId === null || requestParameters.keystoreId === undefined) {
            throw new runtime.RequiredError('keystoreId','Required parameter requestParameters.keystoreId was null or undefined when calling unlockKeystore.');
        }

        if (requestParameters.unlockKeystoreRequest === null || requestParameters.unlockKeystoreRequest === undefined) {
            throw new runtime.RequiredError('unlockKeystoreRequest','Required parameter requestParameters.unlockKeystoreRequest was null or undefined when calling unlockKeystore.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/keystore/{keystoreId}`.replace(`{${"keystoreId"}}`, encodeURIComponent(String(requestParameters.keystoreId))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: UnlockKeystoreRequestToJSON(requestParameters.unlockKeystoreRequest),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(KeystoreFromJSON));
    }

    /**
     * Unlocks a keystore if it exists and returns it
     */
    async unlockKeystore(requestParameters: UnlockKeystoreOperationRequest, initOverrides?: RequestInit): Promise<Array<Keystore>> {
        const response = await this.unlockKeystoreRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
