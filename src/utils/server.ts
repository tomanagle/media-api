import axios from "axios";
import { randomUUID } from "crypto";
import Fastify from "fastify";
import mime from "mime-types";
import {
  S3Client,
  CreateMultipartUploadCommand,
  UploadPartCommand,
  CompleteMultipartUploadCommand,
} from "@aws-sdk/client-s3";

import FastifyMultipart from "@fastify/multipart";
import { env } from "../env";
import { generateVariations } from "../media/media.utils";
import tryToCatch from "./tryToCatch";
import { uploadFile } from "../aws/s3";

const s3 = new S3Client({
  region: env.AWS_REGION,
  credentials: {
    accessKeyId: env.AWS_ACCESS_KEY_ID,
    secretAccessKey: env.AWS_SECRET_ACCESS_KEY,
  },
});

export async function buildServer() {
  const app = Fastify({});

  app.register(FastifyMultipart);

  app.post("/upload", async (req, res) => {
    const file = await req.file({
      limits: {
        fileSize: 1e7, // 10 mb
      },
    });

    if (!file) {
      return res.code(400).send({
        error: "no file",
      });
    }

    const id = randomUUID();

    const originalFileName = file.filename;

    const mimeType = file.mimetype;

    const extension = mime.extension(mimeType);

    if (!extension) {
      return res.code(400).send({
        error: "invalid file type",
      });
    }

    const buffer = await file.toBuffer();

    const v = await generateVariations({
      id,
      extension,
      file: buffer,
    });

    const variations = v
      .filter((i) => i.status === "fulfilled")
      .map((i) => i.value);

    const prom = variations.map((variation) => {
      const { fileName, buffer } = variation;

      return uploadFile({
        fileName,
        mimeType,
        file: buffer,
        s3,
      });
    });

    const result = await Promise.allSettled(prom);

    return result;
  });

  return app;
}
