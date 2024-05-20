"use server";

import { redirect } from "next/navigation";
import { cookies } from "next/headers";
import { verifyToken } from "../lib/jwt";

export async function register(formData: FormData) {
  const rawData = Object.fromEntries(formData);

  const res = await fetch("http://localhost:5000/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(rawData),
  });

  if (!res.ok) {
    throw new Error("couldn't register user");
  }

  const rawToken = await res.json();

  const { session, expiresAt } = await verifyToken(rawToken);

  cookies().set({
    name: "session",
    value: session,
    httpOnly: true,
    path: "/",
    expires: expiresAt,
  });

  redirect("/dashboard");
}
