import { z } from "zod";

export const zodToFormik =
  (schema: z.ZodTypeAny) =>
  (values: any): Record<string, string> => {
    const result = schema.safeParse(values);

    if (result.success) return {};

    const errors: Record<string, string> = {};
    result.error.issues.forEach((issue) => {
      const path = issue.path.join(".");
      errors[path] = issue.message;
    });
    return errors;
  };
