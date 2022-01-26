import { Configuration, DefaultApi } from "./generated";

const configuration = new Configuration({
    basePath: window.location.origin.replace(":8082", ":8081/api"),
});

export default new DefaultApi(configuration);
