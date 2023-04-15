type Awaited<T> = T extends PromiseLike<infer U> ? Awaited<U> : T;

async function tryToCatch<T extends (...args: any) => any>(
  fn: T,
  ...args: Parameters<T>
): Promise<[null, Awaited<ReturnType<T>>] | [any]> {
  try {
    if (!args) {
      return await fn();
    }

    const result: Awaited<ReturnType<T>> = Array.isArray(args)
      ? await fn(...(args as any[]))
      : await fn(args);

    return [null, result];
  } catch (e: any) {
    return [e];
  }
}
export default tryToCatch;
