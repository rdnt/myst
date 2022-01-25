import {MystRestAPI} from "../generated/src";
import * as coreHttp from '@azure/core-http';

class TestTokenCredential implements coreHttp.TokenCredential {
    public token: string;
    public expiresOn: number;

    constructor(token: string, expiresOn?: Date) {
        this.token = token;
        this.expiresOn = expiresOn
            ? expiresOn.getTime()
            : Date.now() + 60 * 60 * 1000;
    }

    async getToken(
        _scopes: string | string[],
        _options?: coreHttp.GetTokenOptions
    ): Promise<coreHttp.AccessToken | null> {
        return {
            token: this.token,
            expiresOnTimestamp: this.expiresOn
        };
    }
}

export default new MystRestAPI({$host: "localhost:8081", allowInsecureConnection: true});
