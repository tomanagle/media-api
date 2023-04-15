import sharp from "sharp";

const variations = [
  {
    name: "thumbnail",
    width: 200,
    height: 200,
  },
  {
    name: "small",
    width: 400,
    height: 400,
  },
  {
    name: "medium",
    width: 800,
    height: 800,
  },
  {
    name: "large",
    width: 1200,
    height: 1200,
  },
];

type GenerateVariationsProps = {
  id: string;
  extension: string;
  file: Buffer;
};

export async function generateVariations({
  id,
  extension,
  file,
}: GenerateVariationsProps) {
  const original = sharp(file);

  const metadata = await original.metadata();

  const { width, height } = metadata;

  if (!width || !height) {
    throw new Error("no width or height");
  }

  const variationsToGenerate = variations.filter((variation) => {
    const { width: variationWidth, height: variationHeight } = variation;

    return width > variationWidth || height > variationHeight;
  });

  const promises = variationsToGenerate.map(async (variation) => {
    const { name, width, height } = variation;

    const resized = original.resize(width, height);

    const buffer = await resized.toBuffer();

    return {
      name,
      buffer,
      width,
      height,
      fileName: `${id}-${name}.${extension}`,
    };
  });

  return Promise.allSettled(promises);
}
