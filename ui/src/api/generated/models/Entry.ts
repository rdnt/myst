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
 * @interface Entry
 */
export interface Entry {
    /**
     * 
     * @type {string}
     * @memberof Entry
     */
    id: string;
    /**
     * 
     * @type {string}
     * @memberof Entry
     */
    label: string;
    /**
     * 
     * @type {string}
     * @memberof Entry
     */
    username: string;
    /**
     * 
     * @type {string}
     * @memberof Entry
     */
    password: string;
    /**
     * 
     * @type {string}
     * @memberof Entry
     */
    notes: string;
}

export function EntryFromJSON(json: any): Entry {
    return EntryFromJSONTyped(json, false);
}

export function EntryFromJSONTyped(json: any, ignoreDiscriminator: boolean): Entry {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'label': json['label'],
        'username': json['username'],
        'password': json['password'],
        'notes': json['notes'],
    };
}

export function EntryToJSON(value?: Entry | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'label': value.label,
        'username': value.username,
        'password': value.password,
        'notes': value.notes,
    };
}

