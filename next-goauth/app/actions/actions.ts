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

  const { accessToken } = await res.json();

  const verifyResult = await verifyToken(accessToken);

  cookies().set({
    name: "session",
    value: accessToken,
    httpOnly: true,
    path: "/",
    expires: verifyResult.expiresAt,
  });

  redirect("/dashboard");
}

export async function login(formData: FormData) {
  const rawData = Object.fromEntries(formData);

  const res = await fetch("http://localhost:5000/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(rawData),
  });

  if (!res.ok) {
    throw new Error("couldn't login user");
  }

  const { accessToken } = await res.json();

  const verifyResult = await verifyToken(accessToken);

  cookies().set({
    name: "session",
    value: accessToken,
    httpOnly: true,
    path: "/",
    expires: verifyResult.expiresAt,
  });

  redirect("/dashboard");
}

export type CurrentUserResponse = {
  user: User;
  accessToken: string | null;
};

export type User = {
  id: string;
  email: string;
};

export async function getCurrentUser() {
  const session = cookies().get("session");
  if (!session) {
    redirect("/login");
  }

  const res = await fetch("http://localhost:5000/current-user", {
    headers: {
      Authorization: `Bearer ${session.value}`,
    },
  });

  if (!res.ok) {
    const body = await res.text();
    console.log(body);

    redirect("/login");
  }

  const data: CurrentUserResponse = await res.json();

  const { accessToken } = data;

  if (accessToken) {
    const verifyResult = await verifyToken(accessToken);

    cookies().set({
      name: "session",
      value: accessToken,
      httpOnly: true,
      path: "/",
      expires: verifyResult.expiresAt,
    });
  }

  return data.user;
}

export async function logout() {
  const session = cookies().get("session");
  if (!session) {
    console.log("session already removed");
    return;
  }

  const res = await fetch("http://localhost:5000/logout", {
    headers: {
      Authorization: `Bearer ${session.value}`,
    },
  });

  if (!res.ok) {
    const body = await res.text();
    console.log(body);
    return;
  }

  cookies().set({
    name: "session",
    value: "",
    httpOnly: true,
    path: "/",
    expires: new Date(0),
  });

  redirect("/login");
}
