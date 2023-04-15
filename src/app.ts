import { env } from "./env";
import { buildServer } from "./utils/server";

async function main() {
  const app = await buildServer();

  await app.listen({
    port: env.PORT,
  });

  console.log("ready to work");
}

main();
