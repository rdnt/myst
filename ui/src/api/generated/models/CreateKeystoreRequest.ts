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
 * @interface CreateKeystoreRequest
 */
export interface CreateKeystoreRequest {
    /**
     * 
     * @type {string}
     * @memberof CreateKeystoreRequest
     */
    name: string;
    /**
     * 
     * @type {string}
     * @memberof CreateKeystoreRequest
     */
    password: string;
}

export function CreateKeystoreRequestFromJSON(json: any): CreateKeystoreRequest {
    return CreateKeystoreRequestFromJSONTyped(json, false);
}

export function CreateKeystoreRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateKeystoreRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'name': json['name'],
        'password': json['password'],
    };
}

export function CreateKeystoreRequestToJSON(value?: CreateKeystoreRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'name': value.name,
        'password': value.password,
    };
}

