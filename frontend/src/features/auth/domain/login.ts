import { z } from "zod";
import { emailField } from "./email";

export const loginSchema = z.object({ email: emailField });

export type Login = Readonly<z.infer<typeof loginSchema>>;

export const emptyLogin: Login = { email: "" };
