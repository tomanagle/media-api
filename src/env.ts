import envSchema from "env-schema";
import { FromSchema } from "json-schema-to-ts";

const schema = {
  type: "object",
  required: [
    "PORT",
    "AWS_REGION",
    "AWS_ACCESS_KEY_ID",
    "AWS_SECRET_ACCESS_KEY",
    "AWS_BUCKET",
  ],
  properties: {
    PORT: {
      type: "number",
      default: 3000,
    },
    AWS_REGION: {
      type: "string",
      default: "us-east-1",
    },
    AWS_ACCESS_KEY_ID: {
      type: "string",
      default: "AKIAIOSFODNN7EXAMPLE",
    },
    AWS_SECRET_ACCESS_KEY: {
      type: "string",
      default: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
    },
    AWS_BUCKET: {
      type: "string",
      default: "my-bucket",
    },
  },
} as const;

type Env = FromSchema<typeof schema>;

export const env = envSchema<Env>({
  schema: schema,
  dotenv: true,
});
