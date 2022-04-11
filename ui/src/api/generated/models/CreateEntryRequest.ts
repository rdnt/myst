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

import { exists, mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface CreateEntryRequest
 */
export interface CreateEntryRequest {
    /**
     * 
     * @type {string}
     * @memberof CreateEntryRequest
     */
    label: string;
    /**
     * 
     * @type {string}
     * @memberof CreateEntryRequest
     */
    username: string;
    /**
     * 
     * @type {string}
     * @memberof CreateEntryRequest
     */
    password: string;
    /**
     * 
     * @type {string}
     * @memberof CreateEntryRequest
     */
    notes: string;
}

export function CreateEntryRequestFromJSON(json: any): CreateEntryRequest {
    return CreateEntryRequestFromJSONTyped(json, false);
}

export function CreateEntryRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateEntryRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'label': json['label'],
        'username': json['username'],
        'password': json['password'],
        'notes': json['notes'],
    };
}

export function CreateEntryRequestToJSON(value?: CreateEntryRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'label': value.label,
        'username': value.username,
        'password': value.password,
        'notes': value.notes,
    };
}

