// import { Configuration, DefaultApi } from "./generated";
//
// const configuration = new Configuration({
//   basePath: window.location.origin.replace(":8082", ":8081") + "/api",
// });
//
// export default new DefaultApi(configuration);
// export * from './generated/models';

OpenAPI.BASE = window.location.origin.replace(":8082", ":8081").replace(":9092", ":9091") + '/api'

import {DefaultService, OpenAPI} from "@/api/generated";

export default DefaultService;
export * from '@/api/generated';
