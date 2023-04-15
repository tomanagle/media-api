import axios from "axios";
import {
  S3Client,
  CreateMultipartUploadCommand,
  UploadPartCommand,
  CompleteMultipartUploadCommand,
  GetObjectCommand,
} from "@aws-sdk/client-s3";
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";
import S3 from "aws-sdk/clients/s3";
import { env } from "../env";
import tryToCatch from "../utils/tryToCatch";

type GetPreSignedUrl = {
  fileName: string;
  mimeType: string;
  partNumber: number;
  s3: S3Client;
};

async function getPreSignedUrl({
  fileName,
  mimeType,
  partNumber,
  s3,
}: GetPreSignedUrl) {
  const createMultipartUploadCommand = new CreateMultipartUploadCommand({
    Bucket: env.AWS_BUCKET,
    Key: fileName,
    ContentType: mimeType,
  });

  const result = await s3.send(createMultipartUploadCommand);

  const uploadId = result.UploadId;
  const key = result.Key;

  const uploadPartCommand = new UploadPartCommand({
    Bucket: process.env.AWS_BUCKET,
    Key: key,
    UploadId: uploadId,
    PartNumber: partNumber,
  });

  const signedUrl = await getSignedUrl(s3, uploadPartCommand, {
    expiresIn: 60,
  });

  return { signedUrl, uploadId, key };
}

type CompleteMultipartUpload = {
  key: string;
  uploadId: string;
  parts: S3.CompletedPart[];
  s3: S3Client;
};

async function completeMultipartUpload({
  key,
  uploadId,
  parts,
  s3,
}: CompleteMultipartUpload) {
  const command = new CompleteMultipartUploadCommand({
    Bucket: process.env.AWS_BUCKET,
    UploadId: uploadId,
    Key: key,
    MultipartUpload: {
      Parts: parts,
    },
  });

  const cmd = await s3.send(command);

  return { ok: true };
}

type UploadFileProps = {
  fileName: string;
  mimeType: string;
  file: Buffer;
  s3: S3Client;
};

export async function uploadFile({
  fileName,
  mimeType,
  file,
  s3,
}: UploadFileProps) {
  const [signError, { signedUrl, key, uploadId }] = await tryToCatch(
    getPreSignedUrl,
    {
      fileName,
      mimeType,
      partNumber: 1,
      s3,
    }
  );

  const [uploadError, data] = await tryToCatch(axios.put, signedUrl, file);

  if (!data) {
    return { ok: false, error: uploadError };
  }

  if (uploadError) {
    return { ok: false, error: uploadError };
  }

  if (!data.headers.etag) {
    return { ok: false, error: "no etag" };
  }

  const etag = data.headers.etag;

  const [completeError] = await tryToCatch(completeMultipartUpload, {
    s3,
    uploadId,
    key: fileName,
    parts: [
      {
        ETag: etag,
        PartNumber: 1,
      },
    ],
  });

  if (completeError) {
    return { ok: false, error: completeError };
  }

  const getCmd = new GetObjectCommand({
    Bucket: env.AWS_BUCKET,
    Key: fileName,
  });

  const url = await getSignedUrl(s3, getCmd);

  return {
    fileName,
    url,
  };
}
